package seed2sdp

import (
	"fmt"

	webrtc "github.com/Gaukas/webrtc_kai/v3"
)

type Sdp struct {
	SDPType       string // "offer", "answer"
	GlobalLines   SdpGlobal
	Payload       string
	Fingerprint   webrtc.DTLSFingerprint
	IceParams     ICEParameters
	IceCandidates []ICECandidate
}

func (s *Sdp) String() string {
	strsdp := fmt.Sprintf(`{"type":"%s","sdp":"v=0\r\n`, s.SDPType)
	strsdp += s.GlobalLines.String()
	strsdp += s.Payload
	strsdp += fmt.Sprintf(`a=fingerprint:%s %s\r\n`, s.Fingerprint.Algorithm, s.Fingerprint.Value)
	strsdp += fmt.Sprintf(`a=ice-ufrag:%s\r\n`, s.IceParams.UsernameFragment)
	strsdp += fmt.Sprintf(`a=ice-pwd:%s\r\n`, s.IceParams.Password)
	strsdp += candidatesToString(s.IceCandidates)
	strsdp += `"}`
	return strsdp
}
