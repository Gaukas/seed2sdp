package main

import (
	"fmt"
	"net"
	"os"

	s2s "github.com/Gaukas/seed2sdp"
	"github.com/pion/webrtc/v3"
)

func usage() {
	fmt.Println("Usage: ./datachannel [seed] [offer/answer]")
}

func main() {
	if len(os.Args) != 3 {
		usage()
		return
	}

	if len(os.Args[1]) < 16 {
		fmt.Println("Mininum seed length: 16 Bytes")
		return
	}

	secret := []byte("By Gaukas Wang")
	seed := []byte(os.Args[1])
	prefix_offer := []byte("Offer_")
	prefix_answer := []byte("Answer_")

	offerFp, _ := s2s.PredictDTLSFingerprint(secret, seed, prefix_offer)
	answerFp, _ := s2s.PredictDTLSFingerprint(secret, seed, prefix_answer)

	fmt.Println("Offer Fp:", offerFp)
	fmt.Println("Answer Fp:", answerFp)

	offerICE, _ := s2s.PredictIceParameters(secret, seed, prefix_offer)
	answerICE, _ := s2s.PredictIceParameters(secret, seed, prefix_answer)

	// fmt.Println("Offer Fp:", offerFp)
	// fmt.Println("Offer ICE:", offerICE)

	if os.Args[2] == "offer" {
		// Used in inflation
		OfferGlobalLines := s2s.SdpGlobal{
			SessionId:   7821628436479802472,
			SessionVer:  1617173148,
			NetworkType: s2s.IN,
			IpaddrType:  s2s.IP4,
			UnicastAddr: net.IPv4(0, 0, 0, 0),
			// SessionName: "",
			// StartingTime: 0,
			// EndingTime: 0,
			GroupBundle: []string{"0"},
			// Payload: "",
		}
		OfferPayload := `m=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:active\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000\r\n`
		config := webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
			HKDFConfig: webrtc.HKDFConfig{
				HkdfSecret:     secret,
				HkdfSalt:       seed,
				HkdfInfoPrefix: prefix_offer,
			},
		}
		offerer(config, OfferGlobalLines, OfferPayload, answerFp, answerICE)
	} else if os.Args[2] == "answer" {
		// Used in inflation
		OfferGlobalLines := s2s.SdpGlobal{
			SessionId:   5615412156857050866,
			SessionVer:  1614192136,
			NetworkType: s2s.IN,
			IpaddrType:  s2s.IP4,
			UnicastAddr: net.IPv4(0, 0, 0, 0),
			// SessionName: "",
			// StartingTime: 0,
			// EndingTime: 0,
			GroupBundle: []string{"0"},
			Payload:     "",
		}
		OfferPayload := `m=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000\r\n`
		config := webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
			HKDFConfig: webrtc.HKDFConfig{
				HkdfSecret:     secret,
				HkdfSalt:       seed,
				HkdfInfoPrefix: prefix_answer,
			},
		}
		answerer(config, OfferGlobalLines, OfferPayload, offerFp, offerICE)
	} else {
		usage()
	}

}
