package seed2sdp

import (
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

	// Check if valid IP
	if RecoveredIP.To16() == nil {
		return nil, ErrInvalidIP
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
	return nil, ErrInvalidSDPType
}

// GroupInflate() should be used when you want to put multiple SDPDeflated into one single SDP
// but why?
func GroupInflate(sds []SDPDeflated) (*SDP, error) {
	var firstSDPType string // get type from first one. will error if there are inconsistent types
	if sds[0].SDPType == SDPOffer {
		firstSDPType = SDPOfferStr
	} else if sds[0].SDPType == SDPAnswer {
		firstSDPType = SDPAnswerStr
	} else {
		return nil, ErrInvalidSDPType
	}

	allICECandidates := []ICECandidate{
		InflateICECandidateFromSD(sds[0]),
	}
	for _, sd := range sds[1:] {
		if sd.SDPType != sds[0].SDPType { // All SDPType must be consistent
			return nil, ErrInvalidSDPType
		}
		allICECandidates = append(allICECandidates, InflateICECandidateFromSD(sd))
	}

	return &SDP{
		SDPType:       firstSDPType,
		IceCandidates: allICECandidates,
	}, nil
}

func (S *SDP) GroupDeflate(UseIPs []net.IP) []SDPDeflated {
	var sliceSdpDeflated = []SDPDeflated{}
	var sdpType uint8

	if S.SDPType == SDPOfferStr {
		sdpType = SDPOffer
	} else if S.SDPType == SDPAnswerStr {
		sdpType = SDPAnswer
	} else {
		return []SDPDeflated{}
	}

	var properCandidates []ICECandidate = []ICECandidate{}

	if len(UseIPs) != 0 { // If specified at least one IP, find those IPs
		for _, c := range S.IceCandidates {
			for idx, UseIP := range UseIPs {
				if c.ipAddr.Equal(UseIP) {
					properCandidates = append(properCandidates, c)

					// Remove the matched IP from the slice (as we use bundle, RTP vs RTCP doesn't matter)
					UseIPs = append(UseIPs[:idx], UseIPs[idx+1:]...)

					break // Out of this current loop interating over UseIPs. Stay in the loop iterating over IceCandidates.
				}
			}
		}
	} else { // Otherwise, extract the all non-internal IPs
		var exclusiveMap = map[string]bool{} // keep all IPs appear once only, as we use bundle (RTP vs RTCP doesn't matter)
		for _, c := range S.IceCandidates {
			if !isPrivateIP(c.ipAddr) { // not private
				if _, ok := exclusiveMap[c.ipAddr.String()]; !ok { // never seen before
					exclusiveMap[c.ipAddr.String()] = true
					properCandidates = append(properCandidates, c)
				}
			}
		}
	}

	// Now process the properCandidates
	for _, properCandidate := range properCandidates {
		sdpDeflated := SDPDeflated{
			SDPType: sdpType,
		}

		// IPUpperUint64, IPLowerUint64
		IP := properCandidate.ipAddr.To16()

		IPUpper := uint64(0)
		IPLower := uint64(0)
		for i := 7; i >= 0; i-- {
			IPUpper += uint64((IP[i+8])&0xFF) << (i * 8)
			IPLower += uint64((IP[i])&0xFF) << (i * 8)
		}
		sdpDeflated.IPUpper64 = IPUpper
		sdpDeflated.IPLower64 = IPLower

		ComposedUint32 := uint32(0)
		ComposedUint32 += uint32(properCandidate.port) << 16         // ComposedUint32[16..31]: ICECandidate.port * (1<<16)
		ComposedUint32 += uint32(properCandidate.tcpType) << 4       // ComposedUint32[4..5]: ICECandidate.tcpType * (1<<4)
		ComposedUint32 += uint32(properCandidate.candidateType) << 2 // ComposedUint32[2..3]: ICECandidate.candidateType * (1<<2)
		ComposedUint32 += uint32(properCandidate.protocol) << 1      // ComposedUint32[1]: ICECandidate.protocol * (1<<1)
		ComposedUint32 += uint32(properCandidate.component) - 1      // ComposedUint32[0]: ICECandidate.ICEComponent-1
		sdpDeflated.Composed32 = ComposedUint32

		sliceSdpDeflated = append(sliceSdpDeflated, sdpDeflated)
	}

	return sliceSdpDeflated
}

// Deflate() will create the SDPDeflated corresponding to the UseIP.
// when UseIP == nil, automatically extrace the SDPDeflated corresponding to the FIRST non-private IP (not recommended due to lack of priority checking)
func (S *SDP) Deflate(UseIP net.IP) SDPDeflated {
	sdpDefs := S.GroupDeflate([]net.IP{UseIP})
	if len(sdpDefs) > 0 { // not empty
		return sdpDefs[0]
	} else {
		return SDPDeflated{}
	}
}
