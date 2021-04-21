package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	s2s "github.com/Gaukas/seed2sdp"
	"github.com/pion/webrtc/v3"
)

var sendmsgcntr uint64 = 0
var recvmsgcntr uint64 = 0

const maxcntr uint64 = 1024

const (
	exampleSecret string = "0xDEADC0DEFADECAFE"
	exampleSalt   string = "0x0EEDFADE"

	offerHKDFPrefix  string = "0xFADE"
	answerHKDFPrefix string = "0xCAFE"
	actpassPayload   string = `m=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000\r\n`
	activePayload    string = `m=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:active\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000\r\n`
)

func usage() {
	fmt.Println("Usage: ./data-channel-client serverIP serverPort seed")
}

func main() {
	if len(os.Args) != 4 || len(os.Args[3]) < 6 { // Min seed length: 6
		usage()
		return
	}

	serverIP := net.ParseIP(os.Args[1])
	if serverIP == nil {
		panic("Not a valid IP!")
	}
	serverPort64, err := strconv.ParseUint(os.Args[2], 10, 16)
	if err != nil {
		panic(fmt.Sprintf("Not a valid Port: %v", err))
	}
	serverPort := uint16(serverPort64)

	clientHkdfParams := s2s.NewHKDFParams().SetSecret(exampleSecret).SetSalt(os.Args[3]).SetInfoPrefix(offerHKDFPrefix)
	serverHkdfParams := s2s.NewHKDFParams().SetSecret(exampleSecret).SetSalt(os.Args[3]).SetInfoPrefix(answerHKDFPrefix)

	dataChannel := s2s.DeclareDatachannel(
		&s2s.DataChannelConfig{
			Label:          "Client DataChannel",
			SelfSDPType:    "offer",
			SelfHkdfParams: clientHkdfParams,
			PeerSDPType:    "answer",
			PeerHkdfParams: serverHkdfParams,
			PeerMedias: []s2s.SDPMedia{
				{
					MediaType:   "application",
					Description: "9 UDP/DTLS/SCTP webrtc-datachannel",
				},
			},
			PeerAttributes: []s2s.SDPAttribute{
				{
					Key:   "group",
					Value: "BUNDLE 0",
				},
				{
					Key:   "setup",
					Value: "active",
				},
				{
					Key:   "mid",
					Value: "0",
				},
				{
					Value: "sendrecv", // Transceivers
				},
				{
					Key:   "sctp-port",
					Value: "5000",
				},
			},
		},
	)

	if dataChannel.Initialize() != nil {
		fmt.Println("[Fatal] Client failed to initialize a data channel instance.")
		panic("Fatal error.")
	}

	////   Set handlers here   ////

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	dataChannel.WebRTCPeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("[Info] Peer Connection changed state to: %s\n", connectionState.String())
		if connectionState.String() == "disconnected" || connectionState.String() == "closed" {
			fmt.Printf("[Info] Peer Connection disconnected\n")
			fmt.Printf("[Info] Shutting down...\n")
			os.Exit(0)
		}
	})

	dataChannel.WebRTCDataChannel.OnOpen(func() {
		fmt.Printf("[Info] Successfully opened Data Channel '%s'-'%d'. \n", dataChannel.WebRTCDataChannel.Label(), dataChannel.WebRTCDataChannel.ID())

		// for range time.NewTicker(5 * time.Second).C {
		// 	message := RandSeq(15)
		// 	fmt.Printf("Sending '%s'\n", message)

		// 	// Send the message as text
		// 	sendErr := dataChannel.WebRTCDataChannel.SendText(message)
		// 	if sendErr != nil {
		// 		panic(sendErr)
		// 	}
		// }
	})

	dataChannel.WebRTCDataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("[Comm] %s: '%s'\n", dataChannel.WebRTCDataChannel.Label(), string(msg.Data))
		recvmsgcntr++
	})

	dataChannel.WebRTCDataChannel.OnClose(func() {
		fmt.Printf("[Warning] Data Channel %s closed\n", dataChannel.WebRTCDataChannel.Label())

		if recvmsgcntr != sendmsgcntr {
			fmt.Printf("[Warning] Packet loss detected. recv: %d, send: %d\n", recvmsgcntr, sendmsgcntr)
		}

		fmt.Printf("[Info] Tearing down Peer Connection\n")
		dataChannel.WebRTCPeerConnection.Close()
	})

	dataChannel.WebRTCDataChannel.OnError(func(err error) {
		fmt.Printf("[Fatal] Data Channel %s errored: %v\n", dataChannel.WebRTCDataChannel.Label(), err)
		fmt.Printf("[Info] Tearing down Peer Connection\n")
		dataChannel.WebRTCPeerConnection.Close()
	})
	//// Stop setting handlers ///

	if dataChannel.CreateOffer() != nil {
		fmt.Println("[FATAL] Client failed to create SDP offer.")
		panic("Fatal error.")
	}

	JsonOffer := s2s.ToJSON(dataChannel.GetLocalDescription())
	ParsedOffer := s2s.ParseSDP(JsonOffer)
	DeflatedOffer := ParsedOffer.Deflate(MyPublicIP(v4))
	// DeflatedOffer := ParsedOffer.Deflate(net.ParseIP("192.168.24.132"))

	fmt.Println("====================== Offer ======================")
	fmt.Println(JsonOffer)
	fmt.Println("===================================================")
	fmt.Println("==================== Def Offer ====================")
	fmt.Println(DeflatedOffer)
	fmt.Println("===================================================")
	AnswerCandidateHost := s2s.ICECandidate{}
	AnswerCandidateHost.
		SetComponent(s2s.ICEComponentRTP).SetProtocol(s2s.UDP).
		SetIpAddr(serverIP).SetPort(serverPort).
		SetCandidateType(s2s.Host)

	AnswerCandidateSrflx := s2s.ICECandidate{}
	AnswerCandidateSrflx.
		SetComponent(s2s.ICEComponentRTP).SetProtocol(s2s.UDP).
		SetIpAddr(serverIP).SetPort(serverPort).
		SetCandidateType(s2s.Srflx)

	time.Sleep(10 * time.Second)

	err = dataChannel.SetAnswer([]s2s.ICECandidate{AnswerCandidateHost, AnswerCandidateSrflx})
	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)

	for sendmsgcntr = 0; sendmsgcntr < maxcntr; sendmsgcntr++ {
		for !dataChannel.ReadyToSend() {
			// fmt.Println("[Info] Data Channel not ready...")
		} // Always wait for ready to send
		fmt.Printf("[Comm] Sending '%d' via %s\n", sendmsgcntr, dataChannel.WebRTCDataChannel.Label())
		sendErr := dataChannel.Send([]byte(fmt.Sprintf("%d", sendmsgcntr)))
		if sendErr != nil {
			panic(sendErr)
		}
	}

	select {}
}
