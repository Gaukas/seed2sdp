module example.com/sdp_to_string

go 1.15

require (
	github.com/Gaukas/seed2sdp v0.0.0
	github.com/Gaukas/webrtc_kai/v3 v3.0.19-a
)

replace github.com/Gaukas/seed2sdp => ../../../seed2sdp

// replace github.com/Gaukas/webrtc_kai/v3 => ../../../../../pion_kai/webrtc_kai

// replace github.com/Gaukas/ice_kai/v2 => ../../../../../pion_kai/ice_kai

// replace github.com/Gaukas/randutil_kai => ../../../../../pion_kai/randutil_kai
