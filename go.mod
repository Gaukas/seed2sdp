module github.com/Gaukas/seed2sdp

go 1.15

require (
	github.com/Gaukas/ice_kai/v2 v2.0.16-a
	github.com/Gaukas/webrtc_kai/v3 v3.0.19-a
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
)

// replace github.com/Gaukas/webrtc_kai/v3 => ../../../pion_kai/webrtc_kai

// replace github.com/Gaukas/ice_kai/v2 => ../../../pion_kai/ice_kai

// replace github.com/Gaukas/randutil_kai => ../../../pion_kai/randutil_kai
