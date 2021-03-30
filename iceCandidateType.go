package seed2sdp

type ICECandidateType uint8

// 3-bits to exchange
const (
	Host  ICECandidateType = iota // 0
	Srflx                         // 1
	Prflx                         // 2
	Relay                         // 3 // 2-bits
	Unknown
)
