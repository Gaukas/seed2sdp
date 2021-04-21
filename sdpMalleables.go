package seed2sdp

import (
	"fmt"
)

// SDPMalleable includes v, o, s, c, t which "do not affect the WebRTC session". (WebRTC For The Curious, page 13)
type SDPMalleables struct {
	Version        uint32            // v=0
	Origin         SDPOrigin         // o=- 0 0 IN IP4 0.0.0.0
	SessionName    string            // s=-
	ConnectionData SDPConnectionData // c=IN IP4 0.0.0.0
	Timing         SDPTiming         // t=0 0
}

func NewSDPMalleables() SDPMalleables {
	return SDPMalleables{
		Version:        0,
		Origin:         NewSDPOrigin(),
		SessionName:    "-",
		ConnectionData: NewSDPConnectionData(),
		Timing:         NewSDPTiming(),
	}
}

// To-Do: Figure out Origin?
func SDPMalleablesFromSeed(HkdfParams *HKDFParams) SDPMalleables {
	return NewSDPMalleables()
}

// v=0
// o=- 0 0 IN IP4 0.0.0.0
// s=-
// c=IN IP4 0.0.0.0
// t=0 0
func (sm *SDPMalleables) String() string {
	strsm := fmt.Sprintf(`v=%d\r\n`, sm.Version)
	strsm += fmt.Sprintf(`o=%s\r\n`, sm.Origin.String())
	strsm += fmt.Sprintf(`s=%s\r\n`, sm.SessionName)
	strsm += fmt.Sprintf(`c=%s\r\n`, sm.ConnectionData.String())
	strsm += fmt.Sprintf(`t=%s\r\n`, sm.Timing.String())
	return strsm
}
