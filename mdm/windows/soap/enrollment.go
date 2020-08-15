package soap

import "github.com/mattrax/xml"

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
