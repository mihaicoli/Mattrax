package soap

import (
	"crypto/x509"
	"encoding/base64"
	"time"

	"github.com/fullsailor/pkcs7"
	"github.com/mattrax/Mattrax/pkg"
	"github.com/mattrax/xml"
)

// EnrollmentRequest contains the device information and identity certificate CSR
type EnrollmentRequest struct {
	XMLName xml.Name      `xml:"s:Envelope"`
	Header  RequestHeader `xml:"s:Header"`
	Body    struct {
		TokenType           string              `xml:"wst:TokenType"`
		RequestType         string              `xml:"wst:RequestType"`
		BinarySecurityToken BinarySecurityToken `xml:"wsse:BinarySecurityToken"`
		AdditionalContext   []ContextItem       `xml:"ac:AdditionalContext>ac:ContextItem"`
	} `xml:"s:Body>wst:RequestSecurityToken"`
}

// BinarySecurityToken contains the CSR for the request and wap-provisioning payload for the response
type BinarySecurityToken struct {
	ValueType    string `xml:"ValueType,attr"`
	EncodingType string `xml:"EncodingType,attr"`
	Value        string `xml:",chardata"`
}

// ParseVerifyCSR parses and verifies the signer (if necessary) of the Binary Security Token.
func (bst BinarySecurityToken) ParseVerifyCSR(verifyIssuer func(*x509.Certificate) error) (csr *x509.CertificateRequest, err error) {
	var decodedCertificateSigningRequest []byte
	if bst.EncodingType == "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd#base64binary" {
		if decodedCertificateSigningRequest, err = base64.StdEncoding.DecodeString(bst.Value); err != nil {
			return nil, pkg.AdvancedError{
				FaultCauser: "s:Receiver",
				FaultType:   "s:MessageFormat",
				FaultReason: "The binary security token could not be decoded",
			}
		}
	} else {
		return nil, pkg.AdvancedError{
			FaultCauser: "s:Receiver",
			FaultType:   "s:MessageFormat",
			FaultReason: "The binary security token format is not supported",
		}
	}

	if bst.ValueType == "http://schemas.microsoft.com/windows/pki/2009/01/enrollment#PKCS7" {
		p7, err := pkcs7.Parse(decodedCertificateSigningRequest)
		if err != nil {
			return nil, pkg.AdvancedError{
				Err:                 err,
				InternalDescription: "error parsing binary security token certificate renewal request",
				FaultCauser:         "s:Receiver",
				FaultType:           "s:MessageFormat",
				FaultReason:         "The binary security token could not be parsed",
			}
		} else if err := p7.Verify(); err != nil {
			return nil, pkg.AdvancedError{
				Err:                 err,
				InternalDescription: "error verifying binary security token certificate renewal request",
				FaultCauser:         "s:Receiver",
				FaultType:           "s:MessageFormat",
				FaultReason:         "The binary security token could not be verified",
			}
		} else if now := time.Now(); now.Before(p7.GetOnlySigner().NotBefore) || now.After(p7.GetOnlySigner().NotAfter) {
			return nil, pkg.AdvancedError{
				Err:                 err,
				InternalDescription: "error expired binary security token certificate renewal request",
				FaultCauser:         "s:Receiver",
				FaultType:           "s:MessageFormat",
				FaultReason:         "The binary security token is expired",
			}
		} else if verifyIssuer != nil {
			if err := verifyIssuer(p7.GetOnlySigner()); err != nil {
				return nil, pkg.AdvancedError{
					Err:                 err,
					InternalDescription: "error invalid issuer for binary security token certificate renewal request",
					FaultCauser:         "s:Receiver",
					FaultType:           "s:MessageFormat",
					FaultReason:         "The binary security token signer could not be verified",
				}
			}
		}

		decodedCertificateSigningRequest = p7.Content
	} else if bst.ValueType != "http://schemas.microsoft.com/windows/pki/2009/01/enrollment#PKCS10" {
		return nil, pkg.AdvancedError{
			FaultCauser: "s:Receiver",
			FaultType:   "s:MessageFormat",
			FaultReason: "The binary security token type is not supported",
		}
	}

	csr, err = x509.ParseCertificateRequest(decodedCertificateSigningRequest)
	if err != nil {
		return nil, pkg.AdvancedError{
			Err:                 err,
			InternalDescription: "error parsing binary security token certificate signing request",
			FaultCauser:         "s:Receiver",
			FaultType:           "s:MessageFormat",
			FaultReason:         "The binary security token could not be parsed",
		}
	} else if err = csr.CheckSignature(); err != nil {
		return nil, pkg.AdvancedError{
			Err:                 err,
			InternalDescription: "error verifying binary security token signature",
			FaultCauser:         "s:Receiver",
			FaultType:           "s:MessageFormat",
			FaultReason:         "The binary security token could not be verified",
		}
	}

	return csr, nil
}

// ContextItem are key value pairs which contains information about the device being enrolled
type ContextItem struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"ac:Value"`
}

// GetAdditionalContextItem retrieves the first AdditionalContext item with the specified name
func (cmd EnrollmentRequest) GetAdditionalContextItem(name string) string {
	for _, contextItem := range cmd.Body.AdditionalContext {
		if contextItem.Name == name {
			return contextItem.Value
		}
	}
	return ""
}

// GetAdditionalContextItems retrieves all AdditionalContext items with a specified name
func (cmd EnrollmentRequest) GetAdditionalContextItems(name string) []string {
	var contextItems []string
	for _, contextItem := range cmd.Body.AdditionalContext {
		if contextItem.Name == name {
			contextItems = append(contextItems, contextItem.Value)
		}
	}
	return contextItems
}

// EnrollmentResponse contains the management client configuration and signed identity certificate
type EnrollmentResponse struct {
	XMLName             xml.Name            `xml:"http://docs.oasis-open.org/ws-sx/ws-trust/200512 RequestSecurityTokenResponseCollection"`
	TokenType           string              `xml:"RequestSecurityTokenResponse>TokenType"`
	DispositionMessage  string              `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollment RequestSecurityTokenResponse>DispositionMessage"` // TODO: Invalid type
	BinarySecurityToken BinarySecurityToken `xml:"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd RequestSecurityTokenResponse>RequestedSecurityToken>BinarySecurityToken"`
	RequestID           int                 `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollment RequestSecurityTokenResponse>RequestID"`
}

// NewEnrollmentResponse creates a generic enrollment response envelope
func NewEnrollmentResponse(relatesTo string, rawProvisioningProfile []byte) ResponseEnvelope {
	var res = ResponseEnvelope{
		Header: ResponseHeader{
			RelatesTo: relatesTo,
		},
		Body: ResponseEnvelopeBody{
			Body: EnrollmentResponse{
				TokenType:          "http://schemas.microsoft.com/5.0.0.0/ConfigurationManager/Enrollment/DeviceEnrollmentToken",
				DispositionMessage: "", // TODO: Wrong type + What does it do?
				BinarySecurityToken: BinarySecurityToken{
					ValueType:    "http://schemas.microsoft.com/5.0.0.0/ConfigurationManager/Enrollment/DeviceEnrollmentProvisionDoc",
					EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd#base64binary",
					Value:        base64.StdEncoding.EncodeToString(rawProvisioningProfile),
				},
				RequestID: 0,
			},
		},
	}
	res.Populate("http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RSTRC/wstep")
	return res
}
