package soap

import (
	"github.com/mattrax/xml"
)

// PolicyRequest contains the user authentication information
type PolicyRequest struct {
	XMLName xml.Name      `xml:"s:Envelope"`
	Header  RequestHeader `xml:"s:Header"`
}

// PolicyResponse contains the policy defining how the identity certificate and key will be generated
type PolicyResponse struct {
	XMLName  xml.Name           `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy GetPoliciesResponse"`
	Response PolicyXCEPResponse `xml:"response"`
	CAs      string             `xml:"cAs"`
	OIDs     []XCEPoID          `xml:"oIDs"`
}

// NillableField can be used to specific no value (hence use the default)
var NillableField = struct {
	Nil string `xml:"xsi:nil,attr"`
}{
	Nil: "true",
}

// PolicyXCEPResponse contains the policy information and its policies
type PolicyXCEPResponse struct {
	PolicyID           string       `xml:"policyID"`
	PolicyFriendlyName string       `xml:"policyFriendlyName"`
	NextUpdateHours    interface{}  `xml:"nextUpdateHours"`
	PoliciesNotChanged interface{}  `xml:"policiesNotChanged"`
	Policies           []XCEPPolicy `xml:"policies"`
}

// XCEPPolicy contains the policies configuration
type XCEPPolicy struct {
	XMLName        xml.Name       `xml:"policy"`
	OIDReferenceID int            `xml:"policyOIDReference"`
	CAs            interface{}    `xml:"cAs"`
	Attributes     XCEPAttributes `xml:"attributes"`
}

// XCEPAttributes contains the policies attributes
type XCEPAttributes struct {
	PolicySchema              int                      `xml:"policySchema"`
	PrivateKeyAttributes      XCEPPrivateKeyAttributes `xml:"privateKeyAttributes"`
	SupersededPolicies        interface{}              `xml:"supersededPolicies"`
	PrivateKeyFlags           interface{}              `xml:"privateKeyFlags"`
	SubjectNameFlags          interface{}              `xml:"subjectNameFlags"`
	EnrollmentFlags           interface{}              `xml:"enrollmentFlags"`
	GeneralFlags              interface{}              `xml:"generalFlags"`
	HashAlgorithmOIDReference int                      `xml:"hashAlgorithmOIDReference"`
	RARequirements            interface{}              `xml:"rARequirements"`
	KeyArchivalAttributes     interface{}              `xml:"keyArchivalAttributes"`
	Extensions                interface{}              `xml:"extensions"`
}

// XCEPPrivateKeyAttributes contains attributes for the private key's generation and usage
type XCEPPrivateKeyAttributes struct {
	MinimalKeyLength      int         `xml:"minimalKeyLength"`
	KeySpec               interface{} `xml:"keySpec"`
	KeyUsageProperty      interface{} `xml:"keyUsageProperty"`
	Permissions           interface{} `xml:"permissions"`
	AlgorithmOIDReference interface{} `xml:"algorithmOIDReference"`
	CryptoProviders       interface{} `xml:"cryptoProviders"`
}

// XCEPoID contains and OID value which can be referenced in a policy
type XCEPoID struct {
	OIDReferenceID int    `xml:"policyOIDReference"`
	DefaultName    string `xml:"defaultName"`
	Group          int    `xml:"group"`
	Value          string `xml:"value"`
}

// NewPolicyResponse creates a generic policy response envelope
func NewPolicyResponse(relatesTo, policyID, policyFriendlyName string) ResponseEnvelope {
	var res = ResponseEnvelope{
		Header: ResponseHeader{
			RelatesTo: relatesTo,
		},
		Body: ResponseEnvelopeBody{
			Body: PolicyResponse{
				Response: PolicyXCEPResponse{
					PolicyID:           policyID,
					PolicyFriendlyName: policyFriendlyName,
					NextUpdateHours:    NillableField,
					PoliciesNotChanged: NillableField,
					Policies: []XCEPPolicy{
						{
							OIDReferenceID: 0, // References to OID defined in OIDs section
							CAs:            NillableField,
							Attributes: XCEPAttributes{
								PolicySchema: 3,
								PrivateKeyAttributes: XCEPPrivateKeyAttributes{
									MinimalKeyLength:      4096,
									KeySpec:               NillableField,
									KeyUsageProperty:      NillableField,
									Permissions:           NillableField,
									AlgorithmOIDReference: NillableField,
									CryptoProviders:       NillableField,
								},
								SupersededPolicies:        NillableField,
								PrivateKeyFlags:           NillableField,
								SubjectNameFlags:          NillableField,
								EnrollmentFlags:           NillableField,
								GeneralFlags:              NillableField,
								HashAlgorithmOIDReference: 0,
								RARequirements:            NillableField,
								KeyArchivalAttributes:     NillableField,
								Extensions:                NillableField,
							},
						},
					},
				},
				OIDs: []XCEPoID{
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
	return res
}
