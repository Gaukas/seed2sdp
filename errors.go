package seed2sdp

import "errors"

var (
	ErrBadSDPDeflated         = errors.New("seed2sdp: bad sdpDeflated input")
	ErrMalformedICEParameters = errors.New("seed2sdp: malformed iceParameters")
	ErrInvalidIP              = errors.New("seed2sdp: invalid IP address")
	ErrInvalidSDPType         = errors.New("seed2sdp: invalid SDP type")
)

// Internal
var (
	errNewConnectionDataNotMatch        = errors.New("seed2sdp: unexpected default connection data")
	errNewConnectionDataBadString       = errors.New("seed2sdp: unexpected default connection data to string")
	errCustomConnectionDataBadString    = errors.New("seed2sdp: unexpected custom connection data to string")
	errSDPDeflatedFromStringError       = errors.New("seed2sdp: SDPDeflatedFromString() returned error")
	errSDPDeflatedFromStringUnexptected = errors.New("seed2sdp: unexpected object from SDPDeflatedFromString()")
	errRecoverIPAddrError               = errors.New("seed2sdp: RecoverIPAddr() returned error")
	errRecoverIPAddrFail                = errors.New("seed2sdp: RecoverIPAddr() failed to recover the IP")

// errNewConnectionDataBadString    = errors.New("The New Connection Data String does not match expectation.")
// errCustomConnectionDataBadString = errors.New("The Custom Connection Data String does not match expectation.")
)
