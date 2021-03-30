package main

import (
	"fmt"
	"net"

	s2s "github.com/Gaukas/seed2sdp"
	webrtc "github.com/Gaukas/webrtc_kai/v3"
)

func main() {
	// To parse into SDP and deflate, then inflate.
	originalSDP := `{"type":"offer","sdp":"v=0\r\no=- 5615412156857050866 1614192136 IN IP4 0.0.0.0\r\ns=- \r\nt=0 0\r\na=group:BUNDLE 0\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000a=fingerprint:sha-256 70:A8:3B:77:1C:7F:A5:EB:DB:D3:57:D7:7F:54:CF:0F:E0:45:F0:7D:60:25:7A:D2:38:64:C3:71:F2:A3:76:A1\r\na=ice-ufrag:bmcKQgPzKgnCMKPL\r\na=ice-pwd:BPwAktFHQCXCbPTjyQJrXnsBUMRgxDUT\r\na=candidate:940760967 1 udp 1694498815 73.243.1.11 63447 typ srflx\r\na=end-of-candidates\r\n"}`

	// Used in inflation
	GlobalLines := s2s.SdpGlobal{
		SessionId:   5615412156857050866,
		SessionVer:  1614192136,
		NetworkType: s2s.IN,
		IpaddrType:  s2s.IP4,
		UnicastAddr: net.IPv4(0, 0, 0, 0),
		// SessionName: "",
		// StartingTime: 0,
		// EndingTime: 0,
		GroupBundle: []string{"0"},
		// Payload: "",
	}
	Payload := `m=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000`
	Fingerprint := webrtc.DTLSFingerprint{Algorithm: "sha-256", Value: `70:A8:3B:77:1C:7F:A5:EB:DB:D3:57:D7:7F:54:CF:0F:E0:45:F0:7D:60:25:7A:D2:38:64:C3:71:F2:A3:76:A1`}
	IceParams := s2s.ICEParameters{
		UsernameFragment: "bmcKQgPzKgnCMKPL",
		Password:         "BPwAktFHQCXCbPTjyQJrXnsBUMRgxDUT",
	}

	parsedSDP := s2s.ParseSDP(originalSDP)
	// fmt.Println(parsedSDP.String())

	defParsedSDP := parsedSDP.Deflate(nil)
	// fmt.Println("Deflated SDP:", defParsedSDP)

	infDefParsedSDP := *(defParsedSDP.Inflate(GlobalLines, Payload, Fingerprint, IceParams))
	fmt.Println(infDefParsedSDP.String())
	if len(infDefParsedSDP.String()) > 0 {
		return
	} else {
		return
	}
}
