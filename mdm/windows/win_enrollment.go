package windows

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/mdm/windows/wap"
	"github.com/mattrax/xml"
	"github.com/rs/zerolog/log"
)

// Discovery handles the discovery phase of enrollment.
func Discovery(srv *mattrax.Server) http.HandlerFunc {
	enrollmentPolicyServiceURL := getNamedRouteURLOrFatal(srv.GlobalRouter, "winmdm-policy")
	enrollmentServiceURL := getNamedRouteURLOrFatal(srv.GlobalRouter, "winmdm-enrollment")
	federationServiceURL := getNamedRouteURLOrFatal(srv.GlobalRouter, "login")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			return
		}

		var cmd soap.DiscoverRequest
		if errored := soap.Read(&cmd, r, w); errored {
			return
		}

		var res = soap.ResponseEnvelope{
			Header: soap.ResponseHeader{
				RelatesTo: cmd.Header.MessageID,
			},
			Body: soap.ResponseEnvelopeBody{
				Body: soap.DiscoverResponse{
					AuthPolicy:                 "Federated",
					EnrollmentVersion:          cmd.Body.RequestVersion,
					EnrollmentPolicyServiceURL: enrollmentPolicyServiceURL,
					EnrollmentServiceURL:       enrollmentServiceURL,
					AuthenticationServiceURL:   federationServiceURL,
				},
			},
		}
		res.Populate("http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse")
		soap.Respond(res, w)
	}
}

// Policy instructs the client how the generate the identity certificate.
// This endpoint is part of the spec MS-XCEP.
func Policy(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cmd soap.PolicyRequest
		if errored := soap.Read(&cmd, r, w); errored {
			return
		}

		if url, err := url.ParseRequestURI(cmd.Header.To); cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy/IPolicy/GetPolicies" || err != nil || url.Host != srv.Args.Domain || cmd.Header.WSSESecurity.BinarySecurityToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		authenticationToken, err := base64.StdEncoding.DecodeString(cmd.Header.WSSESecurity.BinarySecurityToken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := srv.Auth.Token(string(authenticationToken)); err != nil {
			log.Error().Err(err).Msg("error verifying authentication token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var res = soap.ResponseEnvelope{
			Header: soap.ResponseHeader{
				RelatesTo: cmd.Header.MessageID,
			},
			Body: soap.ResponseEnvelopeBody{
				Body: soap.PolicyResponse{
					Response: soap.PolicyXCEPResponse{
						PolicyID:           "mattrax-identity",
						PolicyFriendlyName: "Mattrax Identity Certificate Policy",
						NextUpdateHours:    soap.NillableField,
						PoliciesNotChanged: soap.NillableField,
						Policies: []soap.XCEPPolicy{
							{
								OIDReferenceID: 0, // References to OID defined in OIDs section
								CAs:            soap.NillableField,
								Attributes: soap.XCEPAttributes{
									PolicySchema: 3,
									PrivateKeyAttributes: soap.XCEPPrivateKeyAttributes{
										MinimalKeyLength:      4096,
										KeySpec:               soap.NillableField,
										KeyUsageProperty:      soap.NillableField,
										Permissions:           soap.NillableField,
										AlgorithmOIDReference: soap.NillableField,
										CryptoProviders:       soap.NillableField,
									},
									SupersededPolicies:        soap.NillableField,
									PrivateKeyFlags:           soap.NillableField,
									SubjectNameFlags:          soap.NillableField,
									EnrollmentFlags:           soap.NillableField,
									GeneralFlags:              soap.NillableField,
									HashAlgorithmOIDReference: 0,
									RARequirements:            soap.NillableField,
									KeyArchivalAttributes:     soap.NillableField,
									Extensions:                soap.NillableField,
								},
							},
						},
					},
					OIDs: []soap.XCEPoID{
						{
							OIDReferenceID: 0,
							DefaultName:    "szOID_OIWSEC_SHA256",
							Group:          2, // 2 = Encryption algorithm identifier
							Value:          "2.16.840.1.101.3.4.2.1",
						},
					},
				},
			},
		}
		res.Populate("http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy/IPolicy/GetPoliciesResponse")
		soap.Respond(res, w)
	}
}

// Enrollment provisions the device's management client and issues it a certificate which is used for authentication
func Enrollment(srv *mattrax.Server) http.HandlerFunc {
	managementServiceURL := getNamedRouteURLOrFatal(srv.GlobalRouter, "winmdm-manage")

	if identityCert, identityKey, err := srv.Cert.Get(context.Background(), "identity"); err == sql.ErrNoRows {
		settings, err := srv.DB.Settings(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("error retrieving settings")
			os.Exit(1)
		}

		if settings.TenantName == "" {
			settings.TenantName = "Mattrax"
		}

		identityCert, identityKey, err = srv.Cert.Create(context.Background(), "identity", false, pkix.Name{
			CommonName: settings.TenantName + " Identity",
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating new identity root certificate")
			os.Exit(1)
		}
	} else if err != nil || identityCert == nil || identityKey == nil {
		log.Fatal().Err(err).Msg("Error loading identity certificates")
		os.Exit(1)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var cmd soap.EnrollmentRequest
		if errored := soap.Read(&cmd, r, w); errored {
			return
		}

		if url, err := url.ParseRequestURI(cmd.Header.To); cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RST/wstep" || err != nil || url.Host != srv.Args.Domain || cmd.Header.WSSESecurity.BinarySecurityToken == "" {
			var res = soap.NewFault("s:Receiver", "s:MessageFormat", "", "The request was not destined for this server", "")
			soap.Respond(res, w)
			return
		}

		authenticationToken, err := base64.StdEncoding.DecodeString(cmd.Header.WSSESecurity.BinarySecurityToken)
		if err != nil {
			var res = soap.NewFault("s:Receiver", "a:InvalidSecurity", "", "The security header could not be parsed", "")
			soap.Respond(res, w)
			return
		}

		authClaims, err := srv.Auth.Token(string(authenticationToken))
		if err != nil {
			log.Error().Err(err).Msg("error verifying authentication token")
			var res = soap.NewFault("s:Receiver", "s:Authentication", "", "The user's authenticity could not be verified", "")
			soap.Respond(res, w)
			return
		}

		settings, err := srv.DB.Settings(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("error retrieving settings")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		}

		if settings.DisableEnrollment {
			var res = soap.NewFault("s:Receiver", "s:Authorization", "NotSupported", "Mattrax device enrollments have been disabled", "")
			soap.Respond(res, w)
			return
		}

		existingDevice, err := srv.DB.GetDeviceByUDID(r.Context(), cmd.GetAdditionalContextItem("DeviceID"))
		if err != nil && err != sql.ErrNoRows {
			log.Error().Err(err).Msg("error checking if managed device already exists")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		}

		if err != sql.ErrNoRows && existingDevice.State != db.DeviceStateUserUnenrolled {
			log.Debug().Msg("Device already enrolled in Mattrax.")
			var res = soap.NewFault("s:Receiver", "s:Authorization", "DeviceCapReached", "This device is already enrolled into Mattrax. Please remove before enrolling again", "")
			soap.Respond(res, w)
			return
		}

		if authClaims.MicrosoftSpecificAuthClaims.TenantID != "" {
			if err := srv.DB.NewAzureADUser(r.Context(), db.NewAzureADUserParams{
				Upn:        authClaims.Subject,
				Fullname:   authClaims.Name,
				AzureadOid: sql.NullString{authClaims.MicrosoftSpecificAuthClaims.ObjectID, true},
			}); err != nil {
				if pgerr, ok := err.(*pq.Error); ok {
					if pgerr.Code == "23505" {
						log.Warn().Err(err).Msg("Ignoring duplicate key warning.") // TODO: Fix SQL query to prevent this being needed
					} else {
						log.Error().Err(err).Msg("error creating new AzureAD user")
						var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
						soap.Respond(res, w)
						return
					}
				}
			}
		}

		device := db.NewDeviceParams{
			Udid:            cmd.GetAdditionalContextItem("DeviceID"),
			State:           db.DeviceStateManaged,
			Name:            cmd.GetAdditionalContextItem("DeviceName"),
			SerialNumber:    cmd.GetAdditionalContextItem("HWDevID"),
			OperatingSystem: cmd.GetAdditionalContextItem("OSVersion"),
			AzureDid:        authClaims.MicrosoftSpecificAuthClaims.DeviceID,
			EnrolledBy:      authClaims.Subject,
		}
		if cmd.GetAdditionalContextItem("EnrollmentType") == "Device" {
			device.EnrollmentType = db.EnrollmentTypeDevice
		} else {
			device.EnrollmentType = db.EnrollmentTypeUser
		}

		if existingDevice.ID == 0 {
			if err := srv.DB.NewDevice(r.Context(), device); err != nil {
				log.Error().Err(err).Msg("error creating new device")
				var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
				soap.Respond(res, w)
				return
			}
		} else {
			if err := srv.DB.NewDeviceReplacingExisting(r.Context(), db.NewDeviceReplacingExistingParams(device)); err != nil {
				log.Error().Err(err).Msg("error updating existing device as new device")
				var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
				soap.Respond(res, w)
				return
			}
			if err := srv.DB.NewDeviceReplacingExistingReset(r.Context(), existingDevice.ID); err != nil {
				log.Error().Err(err).Msg("error updating existing device as new device")
				var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
				soap.Respond(res, w)
				return
			}
		}

		certificateSigningRequestDer, err := base64.StdEncoding.DecodeString(cmd.Body.BinarySecurityToken.Value)
		if err != nil {
			var res = soap.NewFault("s:Receiver", "a:InvalidSecurity", "", "The binary security token could not be decoded", "")
			soap.Respond(res, w)
			return
		}

		csr, err := x509.ParseCertificateRequest(certificateSigningRequestDer)
		if err != nil {
			log.Error().Err(err).Msg("error parsing binary security token certificate signing request")
			var res = soap.NewFault("s:Receiver", "a:InvalidSecurity", "", "The binary security token could not be parsed", "")
			soap.Respond(res, w)
			return
		} else if err = csr.CheckSignature(); err != nil {
			log.Error().Err(err).Msg("error checking binary security token signature")
			var res = soap.NewFault("s:Receiver", "a:InvalidSecurity", "", "The binary security token could not be verified", "")
			soap.Respond(res, w)
			return
		}

		identityCertificate, identityCertificateKey, err := srv.Cert.Get(r.Context(), "identity")
		if err != nil || identityCertificate == nil || identityCertificateKey == nil {
			log.Error().Bool("cert-null", identityCertificate == nil).Bool("key-null", identityCertificateKey == nil).Err(err).Msg("Error retrieving the identity certificate")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		}

		if settings.TenantName == "" {
			settings.TenantName = "Mattrax"
		}

		var certStore = "User"
		if cmd.GetAdditionalContextItem("EnrollmentType") == "Device" {
			certStore = "System"
		}

		var notBefore = time.Now().Add(time.Duration(mathrand.Int31n(120)) * -time.Minute)
		clientCertificate := &x509.Certificate{
			Version:            csr.Version,
			Signature:          csr.Signature,
			SignatureAlgorithm: x509.SHA256WithRSA,
			PublicKey:          csr.PublicKey,
			PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
			Subject: pkix.Name{
				CommonName:         cmd.GetAdditionalContextItem("DeviceID"),
				OrganizationalUnit: []string{"WinMDM"},
			},
			Extensions:      csr.Extensions,
			ExtraExtensions: csr.ExtraExtensions,
			DNSNames:        csr.DNSNames,
			EmailAddresses:  csr.EmailAddresses,
			IPAddresses:     csr.IPAddresses,
			URIs:            csr.URIs,

			SerialNumber:          big.NewInt(2), // TODO: Increasing (Should be unqiue for CA)
			Issuer:                identityCertificate.Issuer,
			NotBefore:             notBefore,
			NotAfter:              notBefore.Add(365 * 24 * time.Hour),
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			BasicConstraintsValid: true,
			IsCA:                  false,
		}

		clientCertificateRaw, err := x509.CreateCertificate(rand.Reader, clientCertificate, identityCertificate, csr.PublicKey, identityCertificateKey)
		if err != nil {
			log.Error().Err(err).Msg("error creating client certificate")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		}

		var DMCLientProviderParameters []wap.WapParameter

		if authClaims.MicrosoftSpecificAuthClaims.DeviceID != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, wap.WapParameter{
				Name:     "AADResourceID",
				Value:    authClaims.MicrosoftSpecificAuthClaims.DeviceID,
				DataType: "string",
			})
		}

		if settings.TenantEmail != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, wap.WapParameter{
				Name:     "HelpEmailAddress",
				Value:    settings.TenantEmail,
				DataType: "string",
			})
		}

		if settings.TenantWebsite != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, wap.WapParameter{
				Name:     "HelpWebsite",
				Value:    settings.TenantWebsite,
				DataType: "string",
			})
		}

		if settings.TenantPhone != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, wap.WapParameter{
				Name:     "HelpPhoneNumber",
				Value:    settings.TenantPhone,
				DataType: "string",
			})
		}

		var wapProvisioningDoc = wap.WapProvisioningDoc{
			Version: "1.1",
			Characteristic: []wap.WapCharacteristic{
				{
					Type: "CertificateStore",
					Characteristics: []wap.WapCharacteristic{
						{
							Type: "Root",
							Characteristics: []wap.WapCharacteristic{
								{
									Type: "System",
									Characteristics: []wap.WapCharacteristic{
										{
											Type: strings.ToUpper(fmt.Sprintf("%x", sha1.Sum(identityCertificate.Raw))),
											Params: []wap.WapParameter{
												{
													Name:  "EncodedCertificate",
													Value: base64.StdEncoding.EncodeToString(identityCertificate.Raw),
												},
											},
										},
									},
								},
							},
						},
						{
							Type: "My",
							Characteristics: []wap.WapCharacteristic{
								{
									Type: certStore,
									Characteristics: []wap.WapCharacteristic{
										{
											Type: strings.ToUpper(fmt.Sprintf("%x", sha1.Sum(clientCertificateRaw))),
											Params: []wap.WapParameter{
												{
													Name:  "EncodedCertificate",
													Value: base64.StdEncoding.EncodeToString(clientCertificateRaw),
												},
											},
										},
										{
											Type: "PrivateKeyContainer",
											Params: []wap.WapParameter{
												{
													Name:  "KeySpec",
													Value: "2",
												},
												{
													Name:  "ContainerName",
													Value: "ConfigMgrEnrollment",
												},
												{
													Name:  "ProviderType",
													Value: "1",
												},
											},
										},
									},
								},
							},
						},
						{
							Type: "My",
							Characteristics: []wap.WapCharacteristic{
								{
									Type: "WSTEP",
									Characteristics: []wap.WapCharacteristic{
										{
											Type: "Renew",
											Params: []wap.WapParameter{
												{
													Name:     "ROBOSupport",
													Value:    "true",
													DataType: "boolean",
												},
												{
													Name:     "RenewPeriod",
													Value:    "41",
													DataType: "integer",
												},
												{
													Name:     "RetryInterval",
													Value:    "7",
													DataType: "integer",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Type: "APPLICATION",
					Params: []wap.WapParameter{
						{
							Name:  "APPID",
							Value: "w7",
						},
						{
							Name:  "PROVIDER-ID",
							Value: "MattraxMDM",
						},
						{
							Name:  "ADDR",
							Value: managementServiceURL,
						},
						{
							Name:  "NAME",
							Value: settings.TenantName,
						},
						{
							Name: "BACKCOMPATRETRYDISABLED",
						},
						{
							Name:  "CONNRETRYFREQ",
							Value: "6",
						},
						{
							Name:  "DEFAULTENCODING",
							Value: "application/vnd.syncml.dm+xml",
						},
						{
							Name:  "INITIALBACKOFFTIME",
							Value: timeInMiliseconds(30 * time.Second),
						},
						{
							Name:  "MAXBACKOFFTIME",
							Value: timeInMiliseconds(120 * time.Second),
						},
						{
							Name:  "SSLCLIENTCERTSEARCHCRITERIA",
							Value: "Subject=CN%3d" + strings.ReplaceAll(url.PathEscape(clientCertificate.Subject.String()), "=", "%3D") + "&Stores=MY%5C" + certStore,
						},
					},
					Characteristics: []wap.WapCharacteristic{
						{
							Type: "APPAUTH",
							Params: []wap.WapParameter{
								{
									Name:  "AAUTHLEVEL",
									Value: "CLIENT",
								},
								{
									Name:  "AAUTHTYPE",
									Value: "DIGEST",
								},
								{
									Name:  "AAUTHSECRET",
									Value: "dummy",
								},
								{
									Name:  "AAUTHDATA",
									Value: "nonce",
								},
							},
						},
						{
							Type: "APPAUTH",
							Params: []wap.WapParameter{
								{
									Name:  "AAUTHLEVEL",
									Value: "APPSRV",
								},
								{
									Name:  "AAUTHTYPE",
									Value: "DIGEST",
								},
								{
									Name:  "AAUTHNAME",
									Value: "dummy",
								},
								{
									Name:  "AAUTHSECRET",
									Value: "dummy",
								},
								{
									Name:  "AAUTHDATA",
									Value: "nonce",
								},
							},
						},
					},
				},
				{
					Type: "DMClient",
					Characteristics: []wap.WapCharacteristic{
						{
							Type: "Provider",
							Characteristics: []wap.WapCharacteristic{
								{
									Type: "MattraxMDM",
									Params: append([]wap.WapParameter{
										{
											Name:     "EntDeviceName",
											Value:    cmd.GetAdditionalContextItem("DeviceName"),
											DataType: "string",
										},
										{
											Name:     "EntDMID",
											Value:    cmd.GetAdditionalContextItem("DeviceID"),
											DataType: "string",
										},
										{
											Name:     "UPN",
											Value:    authClaims.Subject,
											DataType: "string",
										},
									}, DMCLientProviderParameters...),
									Characteristics: []wap.WapCharacteristic{
										{
											Type: "Poll",
											Params: []wap.WapParameter{
												{
													Name:     "IntervalForFirstSetOfRetries",
													Value:    "3",
													DataType: "integer",
												},
												{
													Name:     "NumberOfFirstRetries",
													Value:    "5",
													DataType: "integer",
												},
												{
													Name:     "IntervalForSecondSetOfRetries",
													Value:    "15",
													DataType: "integer",
												},
												{
													Name:     "NumberOfSecondRetries",
													Value:    "8",
													DataType: "integer",
												},
												{
													Name:     "IntervalForRemainingScheduledRetries",
													Value:    "480",
													DataType: "integer",
												},
												{
													Name:     "NumberOfRemainingScheduledRetries",
													Value:    "0",
													DataType: "integer",
												},
												{
													Name:     "PollOnLogin",
													Value:    "true",
													DataType: "boolean",
												},
												{
													Name:     "AllUsersPollOnFirstLogin",
													Value:    "true",
													DataType: "boolean",
												},
											},
										},
										{
											Type: "CustomEnrollmentCompletePage",
											Params: []wap.WapParameter{
												{
													Name:     "Title",
													Value:    "Mattrax Enrollment Complete",
													DataType: "string",
												},
												{
													Name:     "BodyText",
													Value:    "Your device is now being managed by '" + settings.TenantName + "'. Please contact your IT administrators for support if you have any problems.",
													DataType: "string",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		provisioningProfileXML, err := xml.Marshal(wapProvisioningDoc)
		if err != nil {
			log.Error().Err(err).Msg("error marshalling wap provisioning profile")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		}

		var res = soap.ResponseEnvelope{
			Header: soap.ResponseHeader{
				RelatesTo: cmd.Header.MessageID,
			},
			Body: soap.ResponseEnvelopeBody{
				Body: soap.EnrollmentResponse{
					TokenType:          "http://schemas.microsoft.com/5.0.0.0/ConfigurationManager/Enrollment/DeviceEnrollmentToken",
					DispositionMessage: "", // TODO: Wrong type + What does it do?
					BinarySecurityToken: soap.BinarySecurityToken{
						ValueType:    "http://schemas.microsoft.com/5.0.0.0/ConfigurationManager/Enrollment/DeviceEnrollmentProvisionDoc",
						EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd#base64binary",
						Value:        base64.StdEncoding.EncodeToString(provisioningProfileXML),
					},
					RequestID: 0,
				},
			},
		}
		res.Populate("http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RSTRC/wstep")
		soap.Respond(res, w)
	}
}

func timeInMiliseconds(d time.Duration) string {
	return strconv.Itoa(int(d / time.Millisecond))
}

// RsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type RsaPublicKey struct { // TODO: Try and remove this
	N *big.Int
	E int
}

// TODO: Remove Fatal bit and return error. Do fatalling at caller
func getNamedRouteURLOrFatal(r *mux.Router, name string, pairs ...string) string {
	route := r.GetRoute(name)
	if route == nil {
		log.Fatal().Str("name", name).Msg("Error acquiring named route")
	}

	url, err := route.URL(pairs...)
	if err != nil {
		log.Fatal().Str("name", name).Err(err).Msg("Error acquiring url of named route")
	}

	return url.String()
}
