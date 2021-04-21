package seed2sdp

import (
	"fmt"
)

// The lines that does not impact
type SDPOrigin struct {
	SessionId      uint64
	SessionVer     uint32
	ConnectionData SDPConnectionData
}

func NewSDPOrigin() SDPOrigin {
	return SDPOrigin{
		SessionId:      0,
		SessionVer:     0,
		ConnectionData: NewSDPConnectionData(),
	}
}

func (so *SDPOrigin) String() string {
	return fmt.Sprintf(`- %d %d %s`, so.SessionId, so.SessionVer, so.ConnectionData.String())
}
