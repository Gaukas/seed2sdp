package seed2sdp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	SDPOffer     uint8  = 1
	SDPAnswer    uint8  = 2
	SDPOfferStr  string = "offer"
	SDPAnswerStr string = "answer"
)

// SDPDeflated represents the minimal info need to be exchanged for SDP
type SDPDeflated struct {
	SDPType    uint8
	Candidates []DeflatedICECandidate
}

// String() return SDPDeflated in string format (old version compatibility)
func (SD *SDPDeflated) String() string {
	var iceCandidateString string
	for _, candidate := range SD.Candidates {
		iceCandidateString += fmt.Sprintf("%d,%d,%d,", candidate.IPUpper64, candidate.IPLower64, candidate.Composed32)
	}
	// truncate the trailing comma
	if lastidx := len(iceCandidateString) - 1; lastidx >= 0 && iceCandidateString[lastidx] == ',' {
		iceCandidateString = iceCandidateString[:lastidx]
	} else {
		return "ERR_NO_ICE_CANDIDATE" // Bad SDP: no ICE candidate
	}

	return fmt.Sprintf("%d,%s", SD.SDPType, iceCandidateString)
}

func SDPDeflatedFromString(SDS string) (SDPDeflated, error) {
	ParsedSDPD := SDPDeflated{}

	s := strings.Split(string(SDS), ",")

	type64, err := strconv.ParseUint(s[0], 10, 8)
	if err != nil {
		return SDPDeflated{}, err
	}
	ParsedSDPD.SDPType = uint8(type64)

	// Candidates
	s = s[1:]
	for len(s) >= 3 {
		current := s[:3]
		// Parse IPUpper64, IPLower64, Composed32
		IPUpper64, err := strconv.ParseUint(current[0], 10, 64)
		if err != nil {
			continue // permissive to skip invalid candidate
		}
		IPLower64, err := strconv.ParseUint(current[1], 10, 64)
		if err != nil {
			continue // permissive to skip invalid candidate
		}
		Composed64, err := strconv.ParseUint(current[2], 10, 32)
		if err != nil {
			continue // permissive to skip invalid candidate
		}
		Composed32 := uint32(Composed64)

		candidate := DeflatedICECandidate{
			IPUpper64:  IPUpper64,
			IPLower64:  IPLower64,
			Composed32: Composed32,
		}

		ParsedSDPD.Candidates = append(ParsedSDPD.Candidates, candidate)

		s = s[3:]
	}

	return ParsedSDPD, nil
}

func (sd SDPDeflated) Inflate() (*SDP, error) {
	// SDPType
	sdp := &SDP{
		IceCandidates: []ICECandidate{},
	}
	if sd.SDPType == SDPOffer {
		sdp.SDPType = SDPOfferStr
	} else if sd.SDPType == SDPAnswer {
		sdp.SDPType = SDPAnswerStr
	} else {
		return nil, ErrInvalidSDPType
	}

	for _, candidate := range sd.Candidates {
		candidateInflated := candidate.Inflate()
		sdp.IceCandidates = append(sdp.IceCandidates, candidateInflated)
	}

	return sdp, nil
}

func (S *SDP) Deflate(candidateIP []net.IP) *SDPDeflated {
	var sdpDeflated = &SDPDeflated{}

	if S.SDPType == SDPOfferStr {
		sdpDeflated.SDPType = SDPOffer
	} else if S.SDPType == SDPAnswerStr {
		sdpDeflated.SDPType = SDPAnswer
	} else {
		return nil
	}

	var properCandidates []ICECandidate = []ICECandidate{}

	if len(candidateIP) != 0 { // If specified at least one IP, find those IPs
		for _, c := range S.IceCandidates {
			for idx, UseIP := range candidateIP {
				if c.ipAddr.Equal(UseIP) {
					properCandidates = append(properCandidates, c)

					// Remove the matched IP from the slice (as we use bundle, RTP vs RTCP doesn't matter)
					candidateIP = append(candidateIP[:idx], candidateIP[idx+1:]...)

					break // Out of this current loop interating over UseIPs. Stay in the loop iterating over IceCandidates.
				}
			}
		}
	} else { // Otherwise, extract the all non-internal IPs
		var exclusiveMap = map[string]bool{} // keep all IPs appear once only, as we use bundle (RTP vs RTCP doesn't matter)
		for _, c := range S.IceCandidates {
			if !isPrivateIP(c.ipAddr) { // not private
				if _, ok := exclusiveMap[c.ipAddr.String()]; !ok { // never seen before
					exclusiveMap[c.ipAddr.String()] = true
					properCandidates = append(properCandidates, c)
				}
			}
		}
	}

	// Now process the properCandidates
	for _, properCandidate := range properCandidates {
		deflatedCandidate := properCandidate.Deflate()
		sdpDeflated.Candidates = append(sdpDeflated.Candidates, deflatedCandidate)
	}

	return sdpDeflated
}
