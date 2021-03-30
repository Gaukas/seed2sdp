package seed2sdp

type ICENetworkType uint8

// Derived locally
const (
	UDP4 ICENetworkType = iota
	UDP6
	TCP4
	TCP6
	BADNETWORKTYPE
)

func (intype ICENetworkType) String() string {
	switch intype {
	case UDP4:
		return "udp4"
	case UDP6:
		return "udp6"
	case TCP4:
		return "tcp4"
	case TCP6:
		return "tcp6"
	default:
		return "bad"
	}
}
