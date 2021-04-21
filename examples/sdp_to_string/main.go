package main

import (
	"fmt"
	"net"

	s2s "github.com/Gaukas/seed2sdp"
	"github.com/pion/webrtc/v3"
)

func main() {
	testCandidate := s2s.ICECandidate{}

	testCandidate.SetComponent(s2s.ICEComponentRTP).
		SetProtocol(s2s.UDP).
		SetIpAddr(net.ParseIP("73.243.1.11")).
		SetPort(63447).
		SetCandidateType(s2s.Srflx)

	sampleSDP := s2s.SDP{
		SDPType:    "offer",
		Malleables: s2s.NewSDPMalleables(),
		Medias: []s2s.SDPMedia{
			{
				MediaType:   "application",
				Description: "9 UDP/DTLS/SCTP webrtc-datachannel",
			},
		},
		Attributes: []s2s.SDPAttribute{
			{
				Key:   "group",
				Value: "BUNDLE 0",
			},
			{
				Key:   "setup",
				Value: "actpass",
			},
			{
				Key:   "mid",
				Value: "0",
			},
			{
				Value: "sendrecv", // Transceivers
			},
			{
				Key:   "sctp-port",
				Value: "5000",
			},
		},
		Fingerprint: webrtc.DTLSFingerprint{
			Algorithm: "sha-256",
			Value:     `70:A8:3B:77:1C:7F:A5:EB:DB:D3:57:D7:7F:54:CF:0F:E0:45:F0:7D:60:25:7A:D2:38:64:C3:71:F2:A3:76:A1`,
		},
		IceParams: s2s.ICEParameters{
			UsernameFragment: "bmcKQgPzKgnCMKPL",
			Password:         "BPwAktFHQCXCbPTjyQJrXnsBUMRgxDUT",
		},
		IceCandidates: []s2s.ICECandidate{
			testCandidate,
		},
	}

	fmt.Println(sampleSDP.String())
}
