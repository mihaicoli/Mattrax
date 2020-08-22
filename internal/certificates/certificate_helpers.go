package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/mattrax/Mattrax/internal/db"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// LoadOrGenerate retrieves a certificate by id and if it is not found generates a new one
func LoadOrGenerate(ctx context.Context, q *db.Queries, id string, subject pkix.Name) (cert *x509.Certificate, key *rsa.PrivateKey, err error) {
	rawCert, err := q.GetRawCert(ctx, id)
	if err == sql.ErrNoRows {
		cert, certRaw, key, keyRaw, err := GenerateCertificate(subject)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error generating new certificate")
		}

		if err := q.CreateRawCert(ctx, db.CreateRawCertParams{
			ID:   id,
			Cert: certRaw,
			Key:  keyRaw,
		}); err != nil {
			return nil, nil, err
		}

		log.Info().Str("id", id).Msg("Generated new certificate")
		return cert, key, nil
	} else if err != nil {
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

// GenerateCertificate takes care of generating a new CA certificate
func GenerateCertificate(subject pkix.Name) (cert *x509.Certificate, certRaw []byte, key *rsa.PrivateKey, keyRaw []byte, err error) {
	key, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	publicKeyBytes := x509.MarshalPKCS1PublicKey(&key.PublicKey)
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

	certRaw, err = x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return cert, certRaw, key, x509.MarshalPKCS1PrivateKey(key), nil
}
