package seed2sdp

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"time"

	webrtc "github.com/Gaukas/webrtc_kai/v3"
	"golang.org/x/crypto/hkdf"
)

func GetCertificate(Secret []byte, Salt []byte, InfoPrefix []byte, sk *ecdsa.PrivateKey) (webrtc.Certificate, error) {

	firstCertReader := hkdf.New(sha256.New, Secret, Salt, append(InfoPrefix, []byte("x509Cert")...)) // initConfiguration() from peerconnection.go in webrtc_kai
	firstCertificate, err := webrtc.GenerateCertificateWithReader(sk, firstCertReader)

	if err != nil {
		return webrtc.Certificate{}, err
	}

	firstFps, err := firstCertificate.GetFingerprints()

	for true {
		time.Sleep(1000000 * time.Nanosecond)                                                           // every 1ms
		nextCertReader := hkdf.New(sha256.New, Secret, Salt, append(InfoPrefix, []byte("x509Cert")...)) // initConfiguration() from peerconnection.go in webrtc_kai
		nextCertificate, err := webrtc.GenerateCertificateWithReader(sk, nextCertReader)

		if err != nil {
			return webrtc.Certificate{}, err
		}

		nextFps, err := nextCertificate.GetFingerprints()

		if nextFps[0].Value > firstFps[0].Value {
			return *nextCertificate, nil
		} else if nextFps[0].Value < firstFps[0].Value {
			return *firstCertificate, nil
		}
	}

	return webrtc.Certificate{}, errors.New("while(1) broken without return")
}

func getPrivkey(Secret []byte, Salt []byte, InfoPrefix []byte) (*ecdsa.PrivateKey, error) {
	KeyReader := hkdf.New(sha256.New, Secret, Salt, append(InfoPrefix, []byte("ECDSAKey")...)) // initConfiguration() from peerconnection.go in webrtc_kai
	sk, err := ecdsa.GenerateKey(elliptic.P256(), KeyReader)
	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}
	return sk, nil
}
