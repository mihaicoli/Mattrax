package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	mathrand "math/rand"
	"sync"
	"time"

	"github.com/mattrax/Mattrax/internal/db"
)

// Service handles certificate generation, retrieval and signing on behalf of the rest of the server.
type Service struct {
	authenticationPrivateKey *rsa.PrivateKey
	authenticationLock       sync.RWMutex

	identityCertificate *x509.Certificate
	identityPrivateKey  *rsa.PrivateKey
	identityLock        sync.RWMutex
}

// IsIssuerIdentity verifies if the certificate was issued by the Identity certificate
func (s *Service) IsIssuerIdentity(cert *x509.Certificate) error {
	signerVerificationOpts := x509.VerifyOptions{
		Roots: x509.NewCertPool(),
	}

	s.identityLock.RLock()
	signerVerificationOpts.Roots.AddCert(s.identityCertificate)
	s.identityLock.RUnlock()

	_, err := cert.Verify(signerVerificationOpts)
	return err
}

// IdentitySignCSR will sign a csr with the Identity certificate
func (s *Service) IdentitySignCSR(csr *x509.CertificateRequest, subject pkix.Name) (*x509.Certificate, *x509.Certificate, []byte, error) {
	s.identityLock.RLock()
	var identityCertificate = s.identityCertificate
	var identityCertificateKey = s.identityPrivateKey
	s.identityLock.RUnlock()

	var notBefore = time.Now().Add(time.Duration(mathrand.Int31n(120)) * -time.Minute)
	clientCertificate := &x509.Certificate{
		Version:            csr.Version,
		Signature:          csr.Signature,
		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKey:          csr.PublicKey,
		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		Subject:            subject,
		Extensions:         csr.Extensions,
		ExtraExtensions:    csr.ExtraExtensions,
		DNSNames:           csr.DNSNames,
		EmailAddresses:     csr.EmailAddresses,
		IPAddresses:        csr.IPAddresses,
		URIs:               csr.URIs,

		SerialNumber:          big.NewInt(2), // TODO: Increasing (Should be unqiue for CA)
		Issuer:                identityCertificate.Issuer,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	rawSignedCert, err := x509.CreateCertificate(rand.Reader, clientCertificate, identityCertificate, csr.PublicKey, identityCertificateKey)
	return identityCertificate, clientCertificate, rawSignedCert, err
}

// AuthenticationKey returns the private key used for authentication
func (s *Service) AuthenticationKey() *rsa.PrivateKey {
	s.authenticationLock.RLock()
	var pk = s.authenticationPrivateKey
	s.authenticationLock.RUnlock()
	return pk
}

// New initialises a new certificate service
func New(q *db.Queries) (s *Service, err error) {
	s = &Service{}
	if _, s.authenticationPrivateKey, err = LoadOrGenerate(context.Background(), q, "authentication", pkix.Name{
		CommonName: "Mattrax Authentication",
	}); err != nil {
		return nil, err
	}
	if s.identityCertificate, s.identityPrivateKey, err = LoadOrGenerate(context.Background(), q, "identity", pkix.Name{
		CommonName: "Mattrax Identity",
	}); err != nil {
		return nil, err
	}

	return s, nil
}
