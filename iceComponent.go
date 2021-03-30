package seed2sdp

type ICEComponent uint8

// 1-bit to exchange
const (
	ICEComponentUnknown ICEComponent = iota // 0
	ICEComponentRTP                         // 1
	ICEComponentRTCP                        // 2
)
