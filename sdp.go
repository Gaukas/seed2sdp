package seed2sdp

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/pion/webrtc/v3"
)

type SDP struct {
	SDPType       string                 // value of "type" key
	Malleables    SDPMalleables          // v, o, s, t, c lines in "sdp" key's value
	Medias        []SDPMedia             // m lines in "sdp" key's value
	Attributes    []SDPAttribute         // a lines in "sdp" key's value
	Fingerprint   webrtc.DTLSFingerprint // Also an attribute, but calculated
	IceParams     ICEParameters          // Also 2 attribute, but calculated
	IceCandidates []ICECandidate         // Also attributes, but calculated
}

func (s *SDP) String() string {
	// Fixed line
	strsdp := fmt.Sprintf(`{"type":"%s","sdp":"`, s.SDPType)

	// v, o, s, t, c lines which does not impact negotiation but may be used in renegotiation
	strsdp += s.Malleables.String()

	// Application Specific Attributes, not generated from seed
	for _, m := range s.Medias {
		strsdp += fmt.Sprintf(`m=%s\r\n`, m.String())
	}

	// Application Specific Attributes, not generated from seed
	for _, a := range s.Attributes {
		strsdp += fmt.Sprintf(`a=%s\r\n`, a.String())
	}

	// fingerprint, as a specific, seed-generated Attribute
	fingerprintAttr := SDPAttribute{
		Key:   "fingerprint",
		Value: fmt.Sprintf(`%s %s`, s.Fingerprint.Algorithm, s.Fingerprint.Value),
	}
	strsdp += fmt.Sprintf(`a=%s\r\n`, fingerprintAttr.String())

	// ice-ufrag, as a specific, seed-generated Attribute
	iceUfrag := SDPAttribute{
		Key:   "ice-ufrag",
		Value: s.IceParams.UsernameFragment,
	}
	strsdp += fmt.Sprintf(`a=%s\r\n`, iceUfrag.String())

	// ice-pwd, as a specific, seed-generated Attribute
	icePwd := SDPAttribute{
		Key:   "ice-pwd",
		Value: s.IceParams.Password,
	}
	strsdp += fmt.Sprintf(`a=%s\r\n`, icePwd.String())

	candidatesAsAttr := candidatesToAttributes(s.IceCandidates)

	for _, c := range candidatesAsAttr {
		strsdp += fmt.Sprintf(`a=%s\r\n`, c.String())
	}
	strsdp += `"}`
	return strsdp
}

// Parse Foundation, Component, Protocol, Priority, IP, Port, CandidateType from 1 single candidate line from SDP
// Todo: Parse TcpType (No Example Text)
func parseCandidate(candidate_text string) ICECandidate {
	trimC := strings.ReplaceAll(strings.ReplaceAll(candidate_text, `\r\n`, ""), `a=candidate:`, "")

	// fmt.Println("Candidate:", trimC)

	splitC := strings.Split(trimC, " ")

	// fmt.Println("Candidate splited into", len(splitC), "parts.")

	foundation_64, _ := strconv.ParseUint(splitC[0], 10, 32)
	component_64, _ := strconv.ParseUint(splitC[1], 10, 8)
	protocol := func() ICENetworkProtocol {
		switch strings.ToLower(splitC[2]) {
		case "udp":
			return UDP
		case "tcp":
			return TCP
		}
		return BADNETWORKPROTOCOL
	}()
	priority_64, _ := strconv.ParseUint(splitC[3], 10, 32)
	port_64, _ := strconv.ParseUint(splitC[5], 10, 16)
	candidateType := func() ICECandidateType {
		switch strings.ToLower(splitC[7]) {
		case "host":
			return Host
		case "srflx":
			return Srflx
		case "prflx":
			return Prflx
		case "relay":
			return Relay
		}
		return Unknown
	}()

	return ICECandidate{
		foundation:    uint32(foundation_64),
		component:     ICEComponent(component_64),
		protocol:      protocol,
		priority:      uint32(priority_64),
		ipAddr:        net.ParseIP(splitC[4]),
		port:          uint16(port_64),
		candidateType: candidateType,
		// tcpType: ,
	}
}

func ParseSDP(sdp_text string) SDP {
	S := SDP{}
	isOffer, _ := regexp.MatchString(`"type":"offer"`, sdp_text)
	isAnswer, _ := regexp.MatchString(`"type":"answer"`, sdp_text)
	if isOffer {
		S.SDPType = "offer" // 0 for offer
	} else if isAnswer {
		S.SDPType = "answer" // 1 for answer
	}

	// Global Lines, Payload, DTLSFingerprint, ICEParams won't be parsed.
	// They are not that helpful for seed-based SDP afterall.

	// Extract all candidates
	reAllCandidate := regexp.MustCompile(`a=candidate:.*?\\r\\n`)
	candidates := reAllCandidate.FindAllString(sdp_text, -1)

	for _, candidate := range candidates {
		S.IceCandidates = append(S.IceCandidates, parseCandidate(candidate))
	}

	return S
}

func (s *SDP) SetMalleables(newval SDPMalleables) {
	s.Malleables = newval
}

func (s *SDP) AddMedia(newval SDPMedia) {
	if s.Medias == nil {
		s.Medias = []SDPMedia{}
	}
	s.Medias = append(s.Medias, newval)
}

func (s *SDP) AddAttrs(newval SDPAttribute) {
	if s.Attributes == nil {
		s.Attributes = []SDPAttribute{}
	}
	s.Attributes = append(s.Attributes, newval)
}

func (s *SDP) SetFingerprint(newval webrtc.DTLSFingerprint) {
	s.Fingerprint = newval
}

func (s *SDP) SetIceParams(newval ICEParameters) {
	s.IceParams = newval
}

func (s *SDP) AddIceCandidates(newval ICECandidate) {
	if s.IceCandidates == nil {
		s.IceCandidates = []ICECandidate{}
	}
	s.IceCandidates = append(s.IceCandidates, newval)
}
