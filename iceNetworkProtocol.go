package seed2sdp

type ICENetworkProtocol uint8

// 1-bits to exchange
const (
	UDP ICENetworkProtocol = iota // 0
	TCP                           // 1 // 1-bit
	BADNETWORKPROTOCOL
)

func (inp ICENetworkProtocol) String() string {
	switch inp {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	default:
		return "bad"
	}
}
