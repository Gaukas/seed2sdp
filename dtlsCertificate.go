package seed2sdp

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/pion/webrtc/v3"
	"golang.org/x/crypto/hkdf"
)

// getPrivkey creates ECDSA private key used in DTLS Certificates
func getPrivkey(hkdfParams *HKDFParams) (*ecdsa.PrivateKey, error) {
	pkReader := hkdf.New(sha256.New, hkdfParams.secret, hkdfParams.secret, append(hkdfParams.infoPrefix, []byte("ECDSAKey")...))
	privkey, err := ecdsa.GenerateKey(elliptic.P256(), pkReader)
	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}
	return privkey, nil
}

// getX509Tpl creates x509 template for x509 Certificates generation used in DTLS Certificates.
func getX509Tpl(hkdfParams *HKDFParams) (*x509.Certificate, error) {
	tplReader := hkdf.New(sha256.New, hkdfParams.secret, hkdfParams.salt, append(hkdfParams.infoPrefix, []byte("X509tpl")...))
	maxBigInt := new(big.Int)
	maxBigInt.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(maxBigInt, big.NewInt(1))
	serialNumber, err := rand.Int(tplReader, maxBigInt)
	if err != nil {
		return &x509.Certificate{}, err
	}

	// Make the Certificate valid from UTC today till next month.
	utcNow := time.Now().UTC()
	validFrom := time.Date(utcNow.Year(), utcNow.Month(), utcNow.Day(), 0, 0, 0, 0, time.UTC)
	validUntil := validFrom.AddDate(0, 1, 0)

	return &x509.Certificate{
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		NotBefore:             validFrom,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:              validUntil,
		SerialNumber:          serialNumber,
		Version:               2,
		Subject:               pkix.Name{CommonName: hex.EncodeToString(hkdfParams.infoPrefix)},
		IsCA:                  true,
	}, nil
}

// NewCertificate() might be ambiguous: we have notices 2 possible version of Certificates.
// Use GetCertificate() instead.
func NewCertificate(hkdfParams *HKDFParams) (webrtc.Certificate, error) {
	privkey, err := getPrivkey(hkdfParams)
	if err != nil {
		return webrtc.Certificate{}, err
	}

	tpl, err := getX509Tpl(hkdfParams)
	if err != nil {
		return webrtc.Certificate{}, err
	}

	x509Reader := hkdf.New(sha256.New, hkdfParams.secret, hkdfParams.salt, append(hkdfParams.infoPrefix, []byte("X509Cert")...))
	tpl.SignatureAlgorithm = x509.ECDSAWithSHA256
	certDER, err := x509.CreateCertificate(x509Reader, tpl, tpl, privkey.Public(), privkey)
	if err != nil {
		return webrtc.Certificate{}, err
	}
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return webrtc.Certificate{}, err
	}

	return webrtc.CertificateFromX509(privkey, cert), nil
}

// GetCertificate() generates DTLS Certificate used for webrtc.
func GetCertificate(hkdfParams *HKDFParams) (webrtc.Certificate, error) {
	firstCert, err := NewCertificate(hkdfParams)
	if err != nil {
		return webrtc.Certificate{}, err
	}
	firstFps, err := firstCert.GetFingerprints()
	if err != nil {
		return webrtc.Certificate{}, err
	}

	for {
		time.Sleep(1000000 * time.Nanosecond) // 1 cert/ms
		nextCert, err := NewCertificate(hkdfParams)
		if err != nil {
			return webrtc.Certificate{}, err
		}
		nextFps, err := nextCert.GetFingerprints()
		if err != nil {
			return webrtc.Certificate{}, err
		}

		if nextFps[0].Value > firstFps[0].Value {
			return nextCert, nil
		} else if nextFps[0].Value < firstFps[0].Value {
			return firstCert, nil
		}
	}

	// return webrtc.Certificate{}, errors.New("while(1) broken")
}
