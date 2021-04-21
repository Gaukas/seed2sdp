package seed2sdp

type DataChannelConfig struct {
	Label          string
	SelfSDPType    string // "offer", "answer"
	SelfHkdfParams *HKDFParams
	PeerSDPType    string // "answer", "offer"
	PeerHkdfParams *HKDFParams
	PeerMedias     []SDPMedia
	PeerAttributes []SDPAttribute
}
