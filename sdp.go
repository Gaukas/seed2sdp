package seed2sdp

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	webrtc "github.com/Gaukas/webrtc_kai/v3"
)

type Sdp struct {
	SDPType       string // "offer", "answer"
	GlobalLines   SdpGlobal
	Payload       string
	Fingerprint   webrtc.DTLSFingerprint
	IceParams     ICEParameters
	IceCandidates []ICECandidate
}

func (s *Sdp) String() string {
	strsdp := fmt.Sprintf(`{"type":"%s","sdp":"v=0\r\n`, s.SDPType)
	strsdp += s.GlobalLines.String()
	strsdp += s.Payload
	strsdp += fmt.Sprintf(`a=fingerprint:%s %s\r\n`, s.Fingerprint.Algorithm, s.Fingerprint.Value)
	strsdp += fmt.Sprintf(`a=ice-ufrag:%s\r\n`, s.IceParams.UsernameFragment)
	strsdp += fmt.Sprintf(`a=ice-pwd:%s\r\n`, s.IceParams.Password)
	strsdp += candidatesToString(s.IceCandidates)
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

func ParseSDP(sdp_text string) Sdp {
	S := Sdp{}
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
