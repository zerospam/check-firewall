package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"time"
)

type CertificateGenerator struct {
	rootTemplate x509.Certificate
	rootKey      *ecdsa.PrivateKey
	notBefore    time.Time
	notAfter     time.Time

	RootKeyPair tls.Certificate
}

func (cg *CertificateGenerator) NewClient(notBefore time.Time, ttl time.Duration) *CertificateGenerator {
	client := &CertificateGenerator{
		notBefore: notBefore,
		notAfter:  notBefore.Add(ttl),
	}
	client.RootKeyPair = client.generateCA()
	return client
}

func (cg *CertificateGenerator) generateCA() tls.Certificate {

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}

	cg.rootKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	cg.rootTemplate = x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
			CommonName:   "Root CA",
		},
		NotBefore:             cg.notBefore,
		NotAfter:              cg.notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &cg.rootTemplate, &cg.rootTemplate, &cg.rootKey.PublicKey, cg.rootKey)
	if err != nil {
		panic(err)
	}

	cert := tls.Certificate{
		PrivateKey: cg.rootKey,
	}

	cert.Certificate = append(cert.Certificate, derBytes)

	return cert
}

func (cg *CertificateGenerator) GenerateClient(commonName string) tls.Certificate {

	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	clientTemplate := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(4),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
			CommonName:   commonName,
		},
		NotBefore:             cg.notBefore,
		NotAfter:              cg.notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA: false,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &clientTemplate, &cg.rootTemplate, &clientKey.PublicKey, cg.rootKey)
	if err != nil {
		panic(err)
	}

	cert := tls.Certificate{
		PrivateKey: clientKey,
	}

	cert.Certificate = append(cert.Certificate, derBytes)

	return cert
}
