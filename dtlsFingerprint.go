package seed2sdp

import webrtc "github.com/Gaukas/webrtc_kai/v3"

// type DTLSFingerprint string

func PredictDTLSFingerprint(Secret []byte, Salt []byte, InfoPrefix []byte) (webrtc.DTLSFingerprint, error) {
	privkey, _ := getPrivkey(Secret, Salt, InfoPrefix)
	certificate, errCert := GetCertificate(Secret, Salt, InfoPrefix, privkey)
	if errCert != nil {
		return webrtc.DTLSFingerprint{}, errCert
	}

	DTLSFPS, errFp := certificate.GetFingerprints()
	if errFp != nil {
		return webrtc.DTLSFingerprint{}, errFp
	}
	return DTLSFPS[0], nil
}
