package seed2sdp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pion/ice/v2"
)

const (
	SDPOffer     uint8  = 1
	SDPAnswer    uint8  = 2
	SDPOfferStr  string = "offer"
	SDPAnswerStr string = "answer"
)

// SDPDeflated represents the minimal info need to be exchanged for SDP
type SDPDeflated struct {
	SDPType    uint8
	IPUpper64  uint64
	IPLower64  uint64
	Composed32 uint32 // [16..31]: ICECandidate.port * (1<<16), [4..5]: ICECandidate.tcpType * (1<<4), [2..3]: ICECandidate.candidateType * (1<<2), [1]: ICECandidate.protocol * (1<<1), [0]: ICECandidate.ICEComponent-1
}

// String() return SDPDeflated in string format (old version compatibility)
func (SD *SDPDeflated) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", SD.SDPType, SD.IPUpper64, SD.IPLower64, SD.Composed32)
}

func SDPDeflatedFromString(SDS string) (SDPDeflated, error) {
	ParsedSDPD := SDPDeflated{}

	s := strings.Split(string(SDS), ",")

	type64, err := strconv.ParseUint(s[0], 10, 8)
	if err != nil {
		return SDPDeflated{}, err
	}
	ParsedSDPD.SDPType = uint8(type64)

	ParsedSDPD.IPUpper64, err = strconv.ParseUint(s[1], 10, 64)
	if err != nil {
		return SDPDeflated{}, err
	}
	ParsedSDPD.IPLower64, err = strconv.ParseUint(s[2], 10, 64)
	if err != nil {
		return SDPDeflated{}, err
	}

	composed64, err := strconv.ParseUint(s[3], 10, 32)
	if err != nil {
		return SDPDeflated{}, err
	}
	ParsedSDPD.Composed32 = uint32(composed64)

	return ParsedSDPD, nil
}

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

	component := ICEComponent((ComposedUint32 >> 0 & 0x01) + 1) // 1 bit, 1/2
	// fmt.Println("Parsed component:", component)
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

func InflateICECandidateFromSD(SD SDPDeflated) ICECandidate {
	return InflateICECandidate(SD.IPUpper64, SD.IPLower64, SD.Composed32)
}

func (sd SDPDeflated) Inflate() (*SDP, error) {
	// SDPType
	if sd.SDPType == SDPOffer {
		return &SDP{
			SDPType: SDPOfferStr,
			IceCandidates: []ICECandidate{
				InflateICECandidateFromSD(sd),
			},
		}, nil
	} else if sd.SDPType == SDPAnswer {
		return &SDP{
			SDPType: SDPAnswerStr,
			IceCandidates: []ICECandidate{
				InflateICECandidateFromSD(sd),
			},
		}, nil
	}
	return nil, errors.New("")
}

func (S *SDP) Deflate(UseIP net.IP) SDPDeflated {
	builtSDPD := SDPDeflated{}

	if S.SDPType == SDPOfferStr {
		builtSDPD.SDPType = SDPOffer
	} else if S.SDPType == SDPAnswerStr {
		builtSDPD.SDPType = SDPAnswer
	} else {
		return SDPDeflated{}
	}

	candidatePtr := (*ICECandidate)(nil)

	if UseIP != nil { // Specified Candidate IP to use, usually it is the Public IP
		for _, c := range S.IceCandidates {
			if c.ipAddr.Equal(UseIP) {
				candidatePtr = &c
				break
			}
		}
		if candidatePtr == nil { // not found
			return SDPDeflated{}
		}
	} else { // Otherwise, extract the first IP
		if len(S.IceCandidates) > 0 {
			candidatePtr = &S.IceCandidates[0]
		} else {
			// Bad SDP
			return SDPDeflated{}
		}
	}

	// IPUpperUint64, IPLowerUint64
	IPFound := (*candidatePtr).ipAddr.To16()

	IPUpper := uint64(0)
	IPLower := uint64(0)
	for i := 7; i >= 0; i-- {
		IPUpper += uint64((IPFound[i+8])&0xFF) << (i * 8)
		IPLower += uint64((IPFound[i])&0xFF) << (i * 8)
	}

	builtSDPD.IPUpper64 = IPUpper
	builtSDPD.IPLower64 = IPLower

	ComposedUint32 := uint32(0)
	ComposedUint32 += uint32((*candidatePtr).port) << 16         // ComposedUint32[16..31]: ICECandidate.port * (1<<16)
	ComposedUint32 += uint32((*candidatePtr).tcpType) << 4       // ComposedUint32[4..5]: ICECandidate.tcpType * (1<<4)
	ComposedUint32 += uint32((*candidatePtr).candidateType) << 2 // ComposedUint32[2..3]: ICECandidate.candidateType * (1<<2)
	ComposedUint32 += uint32((*candidatePtr).protocol) << 1      // ComposedUint32[1]: ICECandidate.protocol * (1<<1)
	ComposedUint32 += uint32((*candidatePtr).component) - 1      // ComposedUint32[0]: ICECandidate.ICEComponent-1

	builtSDPD.Composed32 = ComposedUint32

	return builtSDPD
}
