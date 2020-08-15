package authentication

import "gopkg.in/square/go-jose.v2/jwt"

// BasicClaims contains the generic JWT claims. These are shared between internal and externally issued tokens.
type BasicClaims struct {
	Issuer   string           `json:"iss"`
	IssuedAt *jwt.NumericDate `json:"iat"`
	Expiry   *jwt.NumericDate `json:"exp"`
}

// AuthClaims contains the JWT claims for the authentication token issued by Mattrax's internal authentication
type AuthClaims struct {
	BasicClaims
	MicrosoftSpecificAuthClaims
	Subject            string `json:"sub,omitempty"`
	FullName           string `json:"name,omitempty"`
	Authenticated      bool   `json:"authed,omitempty"`    // This value is set true once any extra authentication has been completed (such as MFA or forced password change). If false DO NOT USE!
	AuthenticationOnly bool   `json:"auth_only,omitempty"` // Set true for Windows MDM enrollment. If true DO NOT USE!
}

// MicrosoftSpecificAuthClaims has the claims for Microsoft AzureAD authentication tokens.
type MicrosoftSpecificAuthClaims struct {
	Audience          string `json:"aud"`
	ObjectID          string `json:"oid"`
	UserPrincipalName string `json:"upn"`
	TenantID          string `json:"tid"`
	Name              string `json:"name"`
	DeviceID          string `json:"deviceid"`
}
