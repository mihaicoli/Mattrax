package wap

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// NewProvisioningDoc returns a new empty WAP Provisioning Document
func NewProvisioningDoc() ProvisioningDoc {
	return ProvisioningDoc{
		Version:        "1.1",
		Characteristic: []Characteristic{},
	}
}

// NewCertStore creates a new "CertificateStore" characteristic on the document
func (doc *ProvisioningDoc) NewCertStore(identityRootCertificate *x509.Certificate, certStore string, clientIssuedCertificateRaw []byte) {
	doc.Characteristic = append(doc.Characteristic, Characteristic{
		Type: "CertificateStore",
		Characteristics: []Characteristic{
			{
				Type: "Root",
				Characteristics: []Characteristic{
					{
						Type: "System",
						Characteristics: []Characteristic{
							{
								Type: strings.ToUpper(fmt.Sprintf("%x", sha1.Sum(identityRootCertificate.Raw))),
								Params: []Parameter{
									{
										Name:  "EncodedCertificate",
										Value: base64.StdEncoding.EncodeToString(identityRootCertificate.Raw),
									},
								},
							},
						},
					},
				},
			},
			{
				Type: "My",
				Characteristics: []Characteristic{
					{
						Type: certStore,
						Characteristics: []Characteristic{
							{
								Type: strings.ToUpper(fmt.Sprintf("%x", sha1.Sum(clientIssuedCertificateRaw))),
								Params: []Parameter{
									{
										Name:  "EncodedCertificate",
										Value: base64.StdEncoding.EncodeToString(clientIssuedCertificateRaw),
									},
								},
							},
							{
								Type: "PrivateKeyContainer",
								Params: []Parameter{
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
				Characteristics: []Characteristic{
					{
						Type: "WSTEP",
						Characteristics: []Characteristic{
							{
								Type: "Renew",
								Params: []Parameter{
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
	})
}

// NewW7Application creates a new "w7 APPLICATION" characteristic on the document
func (doc *ProvisioningDoc) NewW7Application(providerID, tenantName, managementServiceURL, certStore, clientSubject string) {
	doc.Characteristic = append(doc.Characteristic, Characteristic{
		Type: "APPLICATION",
		Params: []Parameter{
			{
				Name:  "APPID",
				Value: "w7",
			},
			{
				Name:  "PROVIDER-ID",
				Value: providerID,
			},
			{
				Name:  "ADDR",
				Value: managementServiceURL,
			},
			{
				Name:  "NAME",
				Value: tenantName,
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
				Value: TimeInMiliseconds(30 * time.Second),
			},
			{
				Name:  "MAXBACKOFFTIME",
				Value: TimeInMiliseconds(120 * time.Second),
			},
			{
				Name:  "SSLCLIENTCERTSEARCHCRITERIA",
				Value: "Subject=" + url.QueryEscape(clientSubject) + "&Stores=MY%5C" + certStore,
			},
		},
		Characteristics: []Characteristic{
			{
				Type: "APPAUTH",
				Params: []Parameter{
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
				Params: []Parameter{
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
	})
}

// NewDMClient creates a new "DMClient" characteristic on the document
func (doc *ProvisioningDoc) NewDMClient(providerID string, providerParameters []Parameter, providerCharacteristics []Characteristic) {
	doc.Characteristic = append(doc.Characteristic, Characteristic{
		Type: "DMClient",
		Characteristics: []Characteristic{
			{
				Type: "Provider",
				Characteristics: []Characteristic{
					{
						Type:            providerID,
						Params:          providerParameters,
						Characteristics: providerCharacteristics,
					},
				},
			},
		},
	})
}

// DefaultPollCharacteristic is a default "Poll" characteristic with default parameter set.
// This characteristic is for use with the DMClient characteristics.
var DefaultPollCharacteristic = Characteristic{
	Type: "Poll",
	Params: []Parameter{
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
}

// TimeInMiliseconds converts a duration into a time string
func TimeInMiliseconds(d time.Duration) string {
	return fmt.Sprintf("%d", d/time.Millisecond)
}
