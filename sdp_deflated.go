package seed2sdp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	ice "github.com/Gaukas/ice_kai/v2"
	webrtc "github.com/Gaukas/webrtc_kai/v3"
)

// SdpDeflated = String(int(SDPType))+","+String(IPUpperUint64)+","+String(IPLowerUint64)+","+String(ComposedUint32)
// ComposedUint32[16..31]: ICECandidate.port * (1<<16)
// ComposedUint32[6..15]: Reserved
// ComposedUint32[4..5]: ICECandidate.tcpType * (1<<4)
// ComposedUint32[2..3]: ICECandidate.candidateType * (1<<2)
// ComposedUint32[1]: ICECandidate.protocol * (1<<1)
// ComposedUint32[0]: ICECandidate.ICEComponent-1
type SdpDeflated string

func RecoverIPAddr(IPUpper uint64, IPLower uint64) (net.IP, error) {
	byteIP := make([]byte, 16)

	for i := 7; i >= 0; i-- {
		byteIP[i+8] = uint8(IPUpper >> (8 * i) & 0xFF)
		byteIP[i] = uint8(IPLower >> (8 * i) & 0xFF)
	}

	RecoveredIP := net.IP(byteIP)

	// RecoveredIP := net.IP{uint8((IPUpper >> 56) & 0xFF),
	// 	uint8((IPUpper >> 48) & 0xFF),
	// 	uint8((IPUpper >> 40) & 0xFF),
	// 	uint8((IPUpper >> 32) & 0xFF),
	// 	uint8((IPUpper >> 24) & 0xFF),
	// 	uint8((IPUpper >> 16) & 0xFF),
	// 	uint8((IPUpper >> 8) & 0xFF),
	// 	uint8((IPUpper >> 0) & 0xFF),
	// 	uint8((IPLower >> 56) & 0xFF),
	// 	uint8((IPLower >> 48) & 0xFF),
	// 	uint8((IPLower >> 40) & 0xFF),
	// 	uint8((IPLower >> 32) & 0xFF),
	// 	uint8((IPLower >> 24) & 0xFF),
	// 	uint8((IPLower >> 16) & 0xFF),
	// 	uint8((IPLower >> 8) & 0xFF),
	// 	uint8((IPLower >> 0) & 0xFF),
	// }

	// Check if valid IP
	if RecoveredIP.To16() == nil {
		return nil, errors.New("Invalid IP")
	}

	return RecoveredIP, nil
}

func InflateICECandidate(IPUpper uint64, IPLower uint64, ComposedUint32 uint32) ICECandidate {
	inflatedIC := ICECandidate{}

	// Recover IP
	recIP, errIP := RecoverIPAddr(IPUpper, IPLower)
	if errIP != nil {
		return ICECandidate{}
	}

	component := ICEComponent((ComposedUint32 >> 0 & 0x01) + 1)   // 1 bit, 1/2
	protocol := ICENetworkProtocol(ComposedUint32 >> 1 & 0x01)    // 1 bit, 0/1
	candidateType := ICECandidateType(ComposedUint32 >> 2 & 0x03) // 2 bits, 0/1/2/3
	tcpType := ice.TCPType(ComposedUint32 >> 4 & 0x03)            // 2 bits, 0/1/2/3
	port := uint16(ComposedUint32 >> 16 & 0xFFFF)                 // 16 bits 0~65535

	inflatedIC.SetComponent(component).
		SetProtocol(protocol).
		SetIpAddr(recIP).
		SetPort(port).
		SetCandidateType(candidateType).
		SetTcpType(tcpType)

	return inflatedIC
}

func (SD *SdpDeflated) Inflate(GlobalLinesOverride SdpGlobal, Payload string, Fp webrtc.DTLSFingerprint, IceParams ICEParameters) *Sdp {
	s := strings.Split(string(*SD), ",")
	IPUpper, _ := strconv.ParseUint(s[1], 10, 64)
	IPLower, _ := strconv.ParseUint(s[2], 10, 64)
	ComposedUint32_64, _ := strconv.ParseUint(s[3], 10, 32)

	// SDPType
	if s[0] == "1" {
		return &Sdp{
			SDPType:     "offer",
			GlobalLines: GlobalLinesOverride,
			Payload:     Payload,
			Fingerprint: Fp,
			IceParams:   IceParams,
			IceCandidates: []ICECandidate{
				InflateICECandidate(IPUpper, IPLower, uint32(ComposedUint32_64)),
			},
		}
	} else {
		return &Sdp{
			SDPType:     "answer",
			GlobalLines: GlobalLinesOverride,
			Payload:     Payload,
			Fingerprint: Fp,
			IceParams:   IceParams,
			IceCandidates: []ICECandidate{
				InflateICECandidate(IPUpper, IPLower, uint32(ComposedUint32_64)),
			},
		}
	}
}

// Abandoned. We don't want to extract those.
// func (S *Sdp) Deflate(UseIP net.IP, GlobalLinesExtracted *SdpGlobal, PayloadExtracted *string, FpExtracted *webrtc.DTLSFingerprint, IceParamsExtracted *ICEParameters) SdpDeflated {
// 	if GlobalLinesExtracted != nil {
// 		*GlobalLinesExtracted = S.GlobalLines
// 	}
// 	if PayloadExtracted != nil {
// 		*PayloadExtracted = S.Payload
// 	}
// 	if FpExtracted != nil {
// 		*FpExtracted = S.Fingerprint
// 	}
// 	if IceParamsExtracted != nil {
// 		*IceParamsExtracted = S.IceParams
// 	}

func (S *Sdp) Deflate(UseIP net.IP) SdpDeflated {
	sdp_d := ""

	if S.SDPType == "offer" {
		sdp_d += "1"
	} else {
		sdp_d += "2"
	}

	c_ptr := (*ICECandidate)(nil)

	if UseIP != nil {
		// Specified Candidate IP to use, usually it is the Public IP
		for _, c := range S.IceCandidates {
			if c.ipAddr.Equal(UseIP) {
				c_ptr = &c
			}
		}
		if c_ptr == nil { // not found
			return ""
		}
	} else {
		if len(S.IceCandidates) > 0 {
			c_ptr = &S.IceCandidates[0]
		} else {
			// Bad SDP
			return ""
		}
	}

	// IPUpperUint64, IPLowerUint64
	IPFound := (*c_ptr).ipAddr.To16()

	IPUpper := uint64(0)
	IPLower := uint64(0)
	for i := 7; i >= 0; i-- {
		IPUpper += uint64((IPFound[i+8])&0xFF) << (i * 8)
		IPLower += uint64((IPFound[i])&0xFF) << (i * 8)
	}

	sdp_d += fmt.Sprintf(",%d,%d", IPUpper, IPLower)

	// ComposedUint32[16..31]: ICECandidate.port * (1<<16)
	// ComposedUint32[6..15]: Reserved
	// ComposedUint32[4..5]: ICECandidate.tcpType * (1<<4)
	// ComposedUint32[2..3]: ICECandidate.candidateType * (1<<2)
	// ComposedUint32[1]: ICECandidate.protocol * (1<<1)
	// ComposedUint32[0]: ICECandidate.ICEComponent-1
	ComposedUint32 := uint32(0)
	ComposedUint32 += uint32((*c_ptr).port) << 16
	ComposedUint32 += uint32((*c_ptr).tcpType) << 4
	ComposedUint32 += uint32((*c_ptr).candidateType) << 2
	ComposedUint32 += uint32((*c_ptr).protocol) << 1
	ComposedUint32 += uint32((*c_ptr).component) - 1

	sdp_d += fmt.Sprintf(",%d", ComposedUint32)

	return SdpDeflated(sdp_d)
}
