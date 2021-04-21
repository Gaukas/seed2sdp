package seed2sdp

import (
	"fmt"
)

// The lines that does not impact
type SDPConnectionData struct {
	Direction string
	IPType    string
	Hostname  string
}

func NewSDPConnectionData() SDPConnectionData {
	return SDPConnectionData{
		Direction: "IN",
		IPType:    "IP4",
		Hostname:  "0.0.0.0",
	}
}

func (cp *SDPConnectionData) String() string {
	return fmt.Sprintf(`%s %s %s`, cp.Direction, cp.IPType, cp.Hostname)
}
