package seed2sdp

import (
	"testing"
)

const (
	newSDPConnDataStr    = "IN IP4 0.0.0.0"
	customSDPConnDataStr = "IN IP6 ::1"
)

func TestSDPConnectionData(t *testing.T) {
	// Time: unlimited
	// lim := test.TimeOut(time.Second * 10)
	// defer lim.Stop()

	NewSDPConnData := NewSDPConnectionData()
	if NewSDPConnData.Direction != "IN" || NewSDPConnData.IPType != "IP4" || NewSDPConnData.Hostname != "0.0.0.0" {
		t.Error(errNewConnectionDataNotMatch)
	}

	if NewSDPConnData.String() != newSDPConnDataStr {
		t.Error(errNewConnectionDataBadString)
	}

	CustomSDPConnData := SDPConnectionData{
		Direction: "IN",
		IPType:    "IP6",
		Hostname:  "::1",
	}

	if CustomSDPConnData.String() != customSDPConnDataStr {
		t.Error(errCustomConnectionDataBadString)
	}

}
