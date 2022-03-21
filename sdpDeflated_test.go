package seed2sdp

import (
	"testing"
)

const (
	offerDefString  = "1,9518543359329632256,0,3351379968"
	answerDefString = "2,9518543359329632256,0,1771503616"
)

func TestSDPDeflated(t *testing.T) {
	// Time: unlimited
	// lim := test.TimeOut(time.Second * 10)
	// defer lim.Stop()
	var err error

	offerSDPDefl, err := SDPDeflatedFromString(offerDefString)
	if err != nil {
		t.Error(errSDPDeflatedFromStringError)
	}
	if offerSDPDefl.SDPType != 1 || len(offerSDPDefl.Candidates) != 1 || offerSDPDefl.Candidates[0].IPUpper64 != 9518543359329632256 || offerSDPDefl.Candidates[0].IPLower64 != 0 || offerSDPDefl.Candidates[0].Composed32 != 3351379968 {
		t.Error(errSDPDeflatedFromStringUnexptected)
	}

	answerSDPDefl, err := SDPDeflatedFromString(answerDefString)
	if err != nil {
		t.Error(errSDPDeflatedFromStringError)
	}
	if answerSDPDefl.SDPType != 2 || len(answerSDPDefl.Candidates) != 1 || answerSDPDefl.Candidates[0].IPUpper64 != 9518543359329632256 || answerSDPDefl.Candidates[0].IPLower64 != 0 || answerSDPDefl.Candidates[0].Composed32 != 1771503616 {
		t.Error(errSDPDeflatedFromStringUnexptected)
	}

	recoveredIP, err := answerSDPDefl.Candidates[0].IPAddr()
	if err != nil {
		t.Error(errRecoverIPAddrError)
	}
	if recoveredIP.String() != "192.168.24.132" {
		t.Error(errRecoverIPAddrFail)
	}
}
