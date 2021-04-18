package main

import (
	"fmt"
	"os"

	s2s "github.com/Gaukas/seed2sdp"
)

const (
	offerHKDFPrefix  string = "0xFEED"
	answerHKDFPrefix string = "0xCAFE"
	offerPayload     string = `m=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000\r\n`
	answerPayload    string = `m=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:active\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000\r\n`
)

func usage() {
	fmt.Println("Usage: ./data-channel-create [seed] [offer/answer]")
	fmt.Println("Min seed length: 6 Bytes") // Magic Number 6.
}

func main() {
	if len(os.Args) != 3 || len(os.Args[1]) < 6 {
		usage()
		return
	}

	programSecret := "0xFEEDCAFEC0DEDEADBEEFBABEBAAD"
	offerHKDFParams := s2s.NewHKDFParams().SetSecret(programSecret).SetSalt(os.Args[1]).SetInfoPrefix(offerHKDFPrefix)
	answerHKDFParams := s2s.NewHKDFParams().SetSecret(programSecret).SetSalt(os.Args[1]).SetInfoPrefix(answerHKDFPrefix)

	// fmt.Println("Offer Fp:", offerFp)
	// fmt.Println("Offer ICE:", offerICE)

	// offerer(offerHKDFParams, answerHKDFParams)

	if os.Args[2] == "offer" {
		offerer(offerHKDFParams, answerHKDFParams)
	} else if os.Args[2] == "answer" {
		answerer(offerHKDFParams, answerHKDFParams)
	} else {
		usage()
	}

}
