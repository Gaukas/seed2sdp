package main

import (
	"fmt"
	"os"
	"time"

	s2s "github.com/Gaukas/seed2sdp"
	webrtc "github.com/Gaukas/webrtc_kai/v3"
)

func offerer(config webrtc.Configuration, GlobalLinesOverride s2s.SdpGlobal, Payload string, PeerFp webrtc.DTLSFingerprint, PeerIceParams s2s.ICEParameters) {
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
		if connectionState.String() == "disconnected" {
			os.Exit(0)
		}
	})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())

		for range time.NewTicker(5 * time.Second).C {
			message := RandSeq(15)
			fmt.Printf("Sending '%s'\n", message)

			// Send the message as text
			sendErr := dataChannel.SendText(message)
			if sendErr != nil {
				panic(sendErr)
			}
		}
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the offer in base64 so we can paste it in browser
	// fmt.Println(ToJSON(*peerConnection.LocalDescription()))
	// return

	// fmt.Println(ToJSON(*peerConnection.LocalDescription()))

	OfferSDP := s2s.ParseSDP(ToJSON(*peerConnection.LocalDescription()))
	DeflatedOffer := OfferSDP.Deflate(MyPublicIP(v4))
	fmt.Println(DeflatedOffer)

	InputAnswer := MustReadStdin()
	AnswerSDP := s2s.SdpDeflated(InputAnswer).Inflate(GlobalLinesOverride, Payload, PeerFp, PeerIceParams)
	inflated_answer := (*AnswerSDP).String()

	// fmt.Println(inflated_answer)

	// Wait for the answer to be pasted
	answer := webrtc.SessionDescription{}
	FromJSON(inflated_answer, &answer)

	// Apply the answer as the remote description
	err = peerConnection.SetRemoteDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block forever
	select {}
}

func answerer(config webrtc.Configuration, GlobalLinesOverride s2s.SdpGlobal, Payload string, PeerFp webrtc.DTLSFingerprint, PeerIceParams s2s.ICEParameters) {
	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
		if connectionState.String() == "disconnected" {
			os.Exit(0)
		}
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Second).C {
				message := RandSeq(15)
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendErr := d.SendText(message)
				if sendErr != nil {
					panic(sendErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	inflated_offer := s2s.SdpDeflated(MustReadStdin()).Inflate(GlobalLinesOverride, Payload, PeerFp, PeerIceParams).String()

	// fmt.Println(inflated_offer)
	// return
	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	FromJSON(inflated_offer, &offer)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	// fmt.Println(string(
	// 	s2s.ParseSDP(ToJSON(*peerConnection.LocalDescription())).
	// 		Deflate(MyPublicIP(v4))))

	AnswerSDP := s2s.ParseSDP(ToJSON(*peerConnection.LocalDescription()))
	// fmt.Println(ToJSON(*peerConnection.LocalDescription()))
	DeflatedAnswer := AnswerSDP.Deflate(MyPublicIP(v4))
	fmt.Println(DeflatedAnswer)

	// Block forever
	select {}
}
