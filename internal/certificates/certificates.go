package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/mattrax/Mattrax/internal/db"
	"github.com/rs/zerolog/log"
)

// RsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type RsaPublicKey struct {
	N *big.Int
	E int
}

// CertificateService handles certificate generation and parsing. It uses a certificate store to persist the certificates and/or keys
type Service struct {
	*db.Queries
}

// Get retrieves and parses a certificates
func (cs Service) Get(ctx context.Context, id string) (cert *x509.Certificate, key *rsa.PrivateKey, err error) {
	rawCert, err := cs.GetRawCert(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if len(rawCert.Cert) != 0 {
		cert, err = x509.ParseCertificate(rawCert.Cert)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(rawCert.Key) != 0 {
		key, err = x509.ParsePKCS1PrivateKey(rawCert.Key)
		if err != nil {
			return nil, nil, err
		}
	}

	return cert, key, nil
}

// Create generates and saves a new certificate
func (cs Service) Create(ctx context.Context, id string, keyOnly bool, subject pkix.Name) (cert *x509.Certificate, key *rsa.PrivateKey, err error) {
	key, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	if !keyOnly {
		publicKeyBytes, err := asn1.Marshal(RsaPublicKey{
			N: key.PublicKey.N,
			E: key.PublicKey.E,
		})
		if err != nil {
			return nil, nil, err
		}

		subjectKeyIDRaw := sha1.Sum(publicKeyBytes)
		notBefore := time.Now().Add(time.Duration(mathrand.Int31n(120)) * -time.Minute) // This randomises the creation time for added security
		cert = &x509.Certificate{
			SerialNumber:                big.NewInt(1),
			Subject:                     subject,
			NotBefore:                   notBefore,
			NotAfter:                    notBefore.Add(365 * 24 * time.Hour),
			KeyUsage:                    x509.KeyUsageCertSign | x509.KeyUsageCRLSign, // TODO: Are they required
			ExtKeyUsage:                 nil,                                          // TODO: What does it do
			UnknownExtKeyUsage:          nil,                                          // TODO: What does it do
			BasicConstraintsValid:       true,                                         // TODO: What does it do
			IsCA:                        true,
			MaxPathLen:                  0, // TODO: What does it do
			SubjectKeyId:                subjectKeyIDRaw[:],
			PermittedDNSDomainsCritical: false, // TODO: What does it do
			PermittedDNSDomains:         nil,   // TODO: What does it do
		}

		certRaw, err := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
		if err != nil {
			return nil, nil, err
		}

		if err := cs.CreateRawCert(ctx, db.CreateRawCertParams{
			ID:   id,
			Cert: certRaw,
			Key:  x509.MarshalPKCS1PrivateKey(key),
		}); err != nil {
			return nil, nil, err
		}

		log.Info().Str("id", id).Msg("Generated new certificate")
	} else {
		if err := cs.CreateRawCert(ctx, db.CreateRawCertParams{
			ID:   id,
			Cert: nil,
			Key:  x509.MarshalPKCS1PrivateKey(key),
		}); err != nil {
			return nil, nil, err
		}

		log.Info().Str("id", id).Msg("Generated new key")
	}

	return cert, key, nil
}
