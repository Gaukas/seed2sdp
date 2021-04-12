package seed2sdp

import "github.com/pion/webrtc/v3"

// type DTLSFingerprint string

func PredictDTLSFingerprint(hkdfParams *HKDFParams) (webrtc.DTLSFingerprint, error) {
	cert, err := GetCertificate(hkdfParams)
	if err != nil {
		return webrtc.DTLSFingerprint{}, err
	}
	DTLSFPS, err := cert.GetFingerprints()
	if err != nil {
		return webrtc.DTLSFingerprint{}, err
	}
	return DTLSFPS[0], nil
}
