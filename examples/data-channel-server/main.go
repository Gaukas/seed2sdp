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
	fmt.Println("Usage: ./data-channel-server bindIP bindPort seed")
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
			Label:          "Server DataChannel",
			SelfSDPType:    "answer",
			SelfHkdfParams: serverHkdfParams,
			PeerSDPType:    "offer",
			PeerHkdfParams: clientHkdfParams,
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
					Value: "actpass", // Client should be actpass, so server active.
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

	//// Set IP/Port here ///

	dataChannel.
		SetIP([]string{serverIP.String()}, s2s.Host).
		SetPort(serverPort) //.SetNetworkTypes()

	/////////////////////////

	if dataChannel.Initialize() != nil {
		fmt.Println("[Fatal] Server failed to initialize a data channel instance.")
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

	dataChannel.WebRTCPeerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		dataChannel.WebRTCDataChannel = d

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
			fmt.Printf("[Comm] %s: '%s'. Sending reply.\n", dataChannel.WebRTCDataChannel.Label(), string(msg.Data))
			recvmsgcntr++
			for !dataChannel.ReadyToSend() {
				// fmt.Println("[Info] Data Channel not ready...")
			} // Always wait for ready to send
			dataChannel.Send(msg.Data)
			sendmsgcntr++
			if recvmsgcntr == maxcntr && sendmsgcntr == recvmsgcntr {
				time.Sleep(10 * time.Second) // Wait for all sent data to arrive
				dataChannel.Close()
			}
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
	})
	//// Stop setting handlers ///

	offerCandidate := s2s.InflateICECandidateFromSD(s2s.SDPDeflated(MustReadStdin()))

	err = dataChannel.SetOffer([]s2s.ICECandidate{offerCandidate})
	if err != nil {
		panic(err)
	}

	if dataChannel.CreateAnswer() != nil {
		fmt.Println("[FATAL] Server failed to create SDP answer.")
		panic("Fatal error.")
	}

	JsonAnswer := s2s.ToJSON(dataChannel.GetLocalDescription())
	ParsedAnswer := s2s.ParseSDP(JsonAnswer)
	DeflatedAnswer := ParsedAnswer.Deflate(serverIP)

	fmt.Println("====================== Answer ======================")
	fmt.Println(JsonAnswer)
	fmt.Println("====================================================")
	fmt.Println("==================== Def Answer ====================")
	fmt.Println(DeflatedAnswer)
	fmt.Println("====================================================")

	select {}
}
