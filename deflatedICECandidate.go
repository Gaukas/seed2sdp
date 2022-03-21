package seed2sdp

import (
	"net"

	"github.com/pion/ice/v2"
)

type DeflatedICECandidate struct {
	IPUpper64  uint64
	IPLower64  uint64
	Composed32 uint32 // [16..31]: ICECandidate.port * (1<<16), [4..5]: ICECandidate.tcpType * (1<<4), [2..3]: ICECandidate.candidateType * (1<<2), [1]: ICECandidate.protocol * (1<<1), [0]: ICECandidate.ICEComponent-1
}

func (c *DeflatedICECandidate) Inflate() ICECandidate {
	inflatedIC := ICECandidate{}

	// Recover IP
	recIP, errIP := c.IPAddr()
	if errIP != nil {
		return ICECandidate{}
	}

	component := ICEComponent((c.Composed32 >> 0 & 0x01) + 1) // 1 bit, 1/2
	// fmt.Println("Parsed component:", component)
	protocol := ICENetworkProtocol(c.Composed32 >> 1 & 0x01)    // 1 bit, 0/1
	candidateType := ICECandidateType(c.Composed32 >> 2 & 0x03) // 2 bits, 0/1/2/3
	tcpType := ice.TCPType(c.Composed32 >> 4 & 0x03)            // 2 bits, 0/1/2/3
	port := uint16(c.Composed32 >> 16 & 0xFFFF)                 // 16 bits 0~65535

	inflatedIC.SetComponent(component).
		SetProtocol(protocol).
		SetIpAddr(recIP).
		SetPort(port).
		SetCandidateType(candidateType).
		SetTcpType(tcpType)

	return inflatedIC
}

func (c *DeflatedICECandidate) IPAddr() (net.IP, error) {
	byteIP := make([]byte, 16)

	for i := 7; i >= 0; i-- {
		byteIP[i+8] = uint8(c.IPUpper64 >> (8 * i) & 0xFF)
		byteIP[i] = uint8(c.IPLower64 >> (8 * i) & 0xFF)
	}

	RecoveredIP := net.IP(byteIP)

	// Check if valid IP
	if RecoveredIP.To16() == nil {
		return nil, ErrInvalidIP
	}

	return RecoveredIP, nil
}
