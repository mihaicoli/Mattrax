package windows

import (
	"crypto/x509/pkix"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/pkg"
	"github.com/mattrax/Mattrax/pkg/soap"
	wap "github.com/mattrax/Mattrax/pkg/wap_provisioning_doc"
	"github.com/mattrax/xml"
	"github.com/rs/zerolog/log"
)

// Discovery handles the discovery phase of enrollment.
func Discovery(srv *mattrax.Server) http.HandlerFunc {
	enrollmentPolicyServiceURL, err := pkg.GetNamedRouteURL(srv.GlobalRouter, "winmdm-policy")
	enrollmentServiceURL, err2 := pkg.GetNamedRouteURL(srv.GlobalRouter, "winmdm-enrollment")
	federationServiceURL, err3 := pkg.GetNamedRouteURL(srv.GlobalRouter, "login")
	if err != nil || err2 != nil || err3 != nil {
		log.Fatal().Err(err).Err(err2).Err(err3).Msg("Error determining route URL") // TODO: Move error handling to main package
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			return
		}

		var cmd soap.DiscoverRequest
		if errored := soap.Read(&cmd, r, w); errored {
			return
		}

		var res = soap.NewDiscoverResponse(cmd.Header.MessageID)
		res.Body.Body = soap.DiscoverResponse{
			AuthPolicy:                 "Federated",
			EnrollmentVersion:          cmd.Body.RequestVersion,
			EnrollmentPolicyServiceURL: enrollmentPolicyServiceURL,
			EnrollmentServiceURL:       enrollmentServiceURL,
			AuthenticationServiceURL:   federationServiceURL,
		}
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

		if url, err := url.ParseRequestURI(cmd.Header.To); cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy/IPolicy/GetPolicies" || err != nil || url.Host != srv.Args.Domain {
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

		soap.Respond(soap.NewPolicyResponse(cmd.Header.MessageID, "mattrax-identity", "Mattrax Identity Certificate Policy"), w)
	}
}

// Enrollment provisions the device's management client and issues it a certificate which is used for authentication
func Enrollment(srv *mattrax.Server) http.HandlerFunc {
	managementServiceURL, err := pkg.GetNamedRouteURL(srv.GlobalRouter, "winmdm-manage") // TODO: Move error handling to main package
	if err != nil {
		log.Fatal().Err(err).Msg("Error determining route URL")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var cmd soap.EnrollmentRequest
		if errored := soap.Read(&cmd, r, w); errored {
			return
		}

		if url, err := url.ParseRequestURI(cmd.Header.To); cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RST/wstep" || err != nil || url.Host != srv.Args.Domain {
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

		settings := srv.Settings.Get()
		if settings.DisableEnrollment {
			var res = soap.NewFault("s:Receiver", "s:Authorization", "NotSupported", "Mattrax device enrollments have been disabled", "")
			soap.Respond(res, w)
			return
		}

		var user db.GetUserRow
		if authClaims.MicrosoftSpecificAuthClaims.TenantID != "" {
			aadUser, err := srv.DB.NewAzureADUser(r.Context(), db.NewAzureADUserParams{
				Upn:        authClaims.Subject,
				Fullname:   authClaims.Name,
				AzureadOid: sql.NullString{authClaims.MicrosoftSpecificAuthClaims.ObjectID, true},
			})
			if err != nil {
				log.Error().Str("upn", authClaims.Subject).Str("oid", authClaims.MicrosoftSpecificAuthClaims.ObjectID).Err(err).Msg("error importing AzureAD user")
				var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
				soap.Respond(res, w)
				return
			}
			user = db.GetUserRow(aadUser) // TODO: Error handling??
		} else if user, err = srv.DB.GetUser(r.Context(), authClaims.Subject); err != nil {
			log.Error().Str("upn", authClaims.Subject).Err(err).Msg("error retrieving user")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		}

		existingDevice, err := srv.DB.GetDeviceByUDID(r.Context(), cmd.GetAdditionalContextItem("DeviceID"))
		if err != nil && err != sql.ErrNoRows {
			log.Error().Err(err).Msg("error checking if managed device already exists")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		} else if err == nil && existingDevice.State != db.DeviceStateUserUnenrolled && existingDevice.State != db.DeviceStateDeploying {
			log.Debug().Int32("id", existingDevice.ID).Msg("Device already enrolled in Mattrax.")
			var res = soap.NewFault("s:Receiver", "s:Authorization", "DeviceCapReached", "This device is already enrolled into Mattrax. Please remove before enrolling again", "")
			soap.Respond(res, w)
			return
		}

		csr, err := cmd.Body.BinarySecurityToken.ParseVerifyCSR(srv.Cert.IsIssuerIdentity)
		if err != nil {
			if aerr, ok := err.(pkg.AdvancedError); ok {
				if err != nil && aerr.InternalDescription != "" {
					log.Error().Err(err).Msg(aerr.InternalDescription)
				}

				var res = soap.NewFault(aerr.FaultCauser, aerr.FaultType, "", aerr.FaultReason, "")
				soap.Respond(res, w)
			}
			return
		}

		device := db.NewDeviceParams{
			Udid:            cmd.GetAdditionalContextItem("DeviceID"),
			State:           db.DeviceStateDeploying,
			Name:            cmd.GetAdditionalContextItem("DeviceName"),
			HwDevID:         cmd.GetAdditionalContextItem("HWDevID"),
			OperatingSystem: cmd.GetAdditionalContextItem("OSVersion"),
			EnrolledBy:      sql.NullString{user.Upn, true},
			AzureDid:        sql.NullString{authClaims.MicrosoftSpecificAuthClaims.DeviceID, authClaims.MicrosoftSpecificAuthClaims.DeviceID != ""},
		}

		var certStore = "User"
		var clientCertSubject = pkix.Name{
			OrganizationalUnit: []string{"WinMDM"},
		}
		if cmd.GetAdditionalContextItem("EnrollmentType") == "Device" {
			certStore = "System"
			device.EnrollmentType = db.EnrollmentTypeDevice
			clientCertSubject.CommonName = cmd.GetAdditionalContextItem("DeviceID")
		} else {
			device.EnrollmentType = db.EnrollmentTypeUser
			clientCertSubject.CommonName = user.Upn
		}

		var deviceID int32
		if existingDevice.ID == 0 {
			if deviceID, err = srv.DB.NewDevice(r.Context(), device); err != nil {
				log.Error().Err(err).Msg("error creating new device")
				var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
				soap.Respond(res, w)
				return
			}
		} else {
			deviceID = existingDevice.ID
			if err := srv.DB.NewDeviceReplacingExisting(r.Context(), db.NewDeviceReplacingExistingParams(device)); err != nil {
				log.Error().Err(err).Msg("error updating existing device as new device")
				var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
				soap.Respond(res, w)
				return
			}
			if err := srv.DB.NewDeviceReplacingExistingResetCache(r.Context(), existingDevice.ID); err != nil {
				log.Error().Err(err).Msg("error resetting cache for device reenrollment")
				var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
				soap.Respond(res, w)
				return
			}
			if err := srv.DB.NewDeviceReplacingExistingResetInventory(r.Context(), existingDevice.ID); err != nil {
				log.Error().Err(err).Msg("error resetting inventory for device reenrollment")
				var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
				soap.Respond(res, w)
				return
			}
		}

		identityCertificate, signedClientCertificate, rawSignedClientCertificate, err := srv.Cert.IdentitySignCSR(csr, clientCertSubject)
		if err != nil {
			log.Error().Err(err).Msg("error creating client certificate")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		}

		var DMCLientProviderParameters = []wap.Parameter{
			{
				Name:     "EntDeviceName",
				Value:    cmd.GetAdditionalContextItem("DeviceName"),
				DataType: "string",
			},
			{
				Name:     "EntDMID",
				Value:    fmt.Sprintf("%d", deviceID),
				DataType: "string",
			},
			{
				Name:     "UPN",
				Value:    user.Upn,
				DataType: "string",
			},
		}

		if authClaims.MicrosoftSpecificAuthClaims.DeviceID != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, wap.Parameter{
				Name:     "AADResourceID",
				Value:    authClaims.MicrosoftSpecificAuthClaims.DeviceID,
				DataType: "string",
			})
		}

		if settings.TenantEmail != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, wap.Parameter{
				Name:     "HelpEmailAddress",
				Value:    settings.TenantEmail,
				DataType: "string",
			})

			var node = "./Vendor/MSFT/DMClient/Provider/" + ProviderID + "/HelpEmailAddress"
			if err := srv.DB.UpdateDeviceInventoryNode(r.Context(), db.UpdateDeviceInventoryNodeParams{
				DeviceID: deviceID,
				Uri:      node,
				Format:   "chr",
				Value:    settings.TenantEmail,
			}); err != nil {
				log.Error().Err(err).Str("node", node).Msg("Error updating device inventory node")
			}
		}

		if settings.TenantWebsite != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, wap.Parameter{
				Name:     "HelpWebsite",
				Value:    settings.TenantWebsite,
				DataType: "string",
			})

			var node = "./Vendor/MSFT/DMClient/Provider/" + ProviderID + "/HelpWebsite"
			if err := srv.DB.UpdateDeviceInventoryNode(r.Context(), db.UpdateDeviceInventoryNodeParams{
				DeviceID: deviceID,
				Uri:      node,
				Format:   "chr",
				Value:    settings.TenantWebsite,
			}); err != nil {
				log.Error().Err(err).Str("node", node).Msg("Error updating device inventory node")
			}
		}

		if settings.TenantPhone != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, wap.Parameter{
				Name:     "HelpPhoneNumber",
				Value:    settings.TenantPhone,
				DataType: "string",
			})

			var node = "./Vendor/MSFT/DMClient/Provider/" + ProviderID + "/HelpPhoneNumber"
			if err := srv.DB.UpdateDeviceInventoryNode(r.Context(), db.UpdateDeviceInventoryNodeParams{
				DeviceID: deviceID,
				Uri:      node,
				Format:   "chr",
				Value:    settings.TenantPhone,
			}); err != nil {
				log.Error().Err(err).Str("node", node).Msg("Error updating device inventory node")
			}
		}

		var node = "./Vendor/MSFT/DMClient/Provider/" + ProviderID + "/ManagementServiceAddress"
		if err := srv.DB.UpdateDeviceInventoryNode(r.Context(), db.UpdateDeviceInventoryNodeParams{
			DeviceID: deviceID,
			Uri:      node,
			Format:   "chr",
			Value:    managementServiceURL,
		}); err != nil {
			log.Error().Err(err).Str("node", node).Msg("Error updating device inventory node")
		}

		var wapProvisioningDoc = wap.NewProvisioningDoc()
		wapProvisioningDoc.NewCertStore(identityCertificate, certStore, rawSignedClientCertificate)
		wapProvisioningDoc.NewW7Application(ProviderID, settings.TenantName, managementServiceURL, certStore, signedClientCertificate.Subject.String())
		wapProvisioningDoc.NewDMClient(ProviderID, DMCLientProviderParameters, []wap.Characteristic{
			wap.DefaultPollCharacteristic,
			{
				Type: "CustomEnrollmentCompletePage",
				Params: []wap.Parameter{
					{
						Name:     "Title",
						Value:    "Mattrax Enrollment Complete",
						DataType: "string",
					},
					{
						Name:     "BodyText",
						Value:    "Welcome " + user.Fullname + ", Your device is now being managed by '" + settings.TenantName + "'. Please contact your IT administrators for support if you have any problems.",
						DataType: "string",
					},
				},
			},
		})

		rawProvisioningProfile, err := xml.Marshal(wapProvisioningDoc)
		if err != nil {
			log.Error().Err(err).Msg("error marshalling wap provisioning profile")
			var res = soap.NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			soap.Respond(res, w)
			return
		}

		fmt.Println(string(rawProvisioningProfile))

		soap.Respond(soap.NewEnrollmentResponse(cmd.Header.MessageID, rawProvisioningProfile), w)
	}
}
