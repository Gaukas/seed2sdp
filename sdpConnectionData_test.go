package seed2sdp

import (
	"errors"
	"testing"
)

const (
	newSDPConnDataStr    = "IN IP4 0.0.0.0"
	customSDPConnDataStr = "IN IP6 ::1"
)

var (
	errNewConnectionDataNotMatch     = errors.New("The Default New Connection Data does not match expectation.")
	errNewConnectionDataBadString    = errors.New("The New Connection Data String does not match expectation.")
	errCustomConnectionDataBadString = errors.New("The Custom Connection Data String does not match expectation.")
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
