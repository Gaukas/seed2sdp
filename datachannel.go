package seed2sdp

import (
	"errors"

	"github.com/pion/webrtc/v3"
)

var DataChannelBufferBytesLim uint64 = 1024

type DataChannel struct {
	config               *DataChannelConfig
	WebRTCSettingEngine  webrtc.SettingEngine
	WebRTCConfiguration  webrtc.Configuration
	WebRTCPeerConnection *webrtc.PeerConnection
	WebRTCDataChannel    *webrtc.DataChannel
}

// DeclareDatachannel sets all the predetermined information needed to
func DeclareDatachannel(newconfig *DataChannelConfig) *DataChannel {
	if newconfig.TxBufferSize == 0 {
		newconfig.TxBufferSize = DataChannelBufferBytesLim // Default: 1K
	}

	// Initialize the struct
	dataChannel := DataChannel{
		config:              newconfig,
		WebRTCSettingEngine: webrtc.SettingEngine{},
		WebRTCConfiguration: webrtc.Configuration{},
	}

	// Setting up WebRTC Configuration
	dataChannel.WebRTCConfiguration.ICEServers = []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	}
	cert, _ := GetCertificate(dataChannel.config.SelfHkdfParams)
	dataChannel.WebRTCConfiguration.Certificates = []webrtc.Certificate{cert}

	// Set ICE Parameters for SettingEngine
	selfICE, _ := PredictIceParameters(dataChannel.config.SelfHkdfParams) // To Be Fixed
	selfICE.UpdateSettingEngine(&dataChannel.WebRTCSettingEngine)

	return &dataChannel
}

func (d *DataChannel) Initialize() error {
	var err error
	api := webrtc.NewAPI(webrtc.WithSettingEngine(d.WebRTCSettingEngine))
	d.WebRTCPeerConnection, err = api.NewPeerConnection(d.WebRTCConfiguration)
	if err != nil {
		return err
	}

	if d.config.SelfSDPType == "offer" {
		d.WebRTCDataChannel, err = d.WebRTCPeerConnection.CreateDataChannel(d.config.Label, nil)
	}
	return nil
}

func (d *DataChannel) CreateLocalDescription() error {
	var localDescription webrtc.SessionDescription
	var err error
	if d.config.SelfSDPType == "offer" {
		localDescription, err = d.WebRTCPeerConnection.CreateOffer(nil)
	} else if d.config.SelfSDPType == "answer" {
		localDescription, err = d.WebRTCPeerConnection.CreateAnswer(nil)
	}

	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(d.WebRTCPeerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = d.WebRTCPeerConnection.SetLocalDescription(localDescription)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
	return nil
}

func (d *DataChannel) CreateOffer() error {
	if d.config.SelfSDPType == "offer" {
		return d.CreateLocalDescription()
	}
	return errors.New("Mismatched SelfSDPType in config: " + d.config.SelfSDPType)
}

func (d *DataChannel) CreateAnswer() error {
	if d.config.SelfSDPType == "answer" {
		return d.CreateLocalDescription()
	}
	return errors.New("Mismatched SelfSDPType in config: " + d.config.SelfSDPType)
}

func (d *DataChannel) GetLocalDescription() *webrtc.SessionDescription {
	return d.WebRTCPeerConnection.LocalDescription()
}

func (d *DataChannel) SetRemoteDescription(remoteCandidates []ICECandidate) error {
	peerFp, _ := PredictDTLSFingerprint(d.config.PeerHkdfParams)
	peerICE, _ := PredictIceParameters(d.config.PeerHkdfParams)

	RemoteSDP := SDP{
		SDPType:       d.config.PeerSDPType,
		Malleables:    SDPMalleablesFromSeed(d.config.PeerHkdfParams),
		Medias:        d.config.PeerMedias,
		Attributes:    d.config.PeerAttributes,
		Fingerprint:   peerFp,
		IceParams:     peerICE,
		IceCandidates: remoteCandidates,
	}

	rdesc := webrtc.SessionDescription{}
	FromJSON(RemoteSDP.String(), &rdesc)

	err := d.WebRTCPeerConnection.SetRemoteDescription(rdesc)
	return err
}

func (d *DataChannel) SetOffer(remoteCandidates []ICECandidate) error {
	if d.config.SelfSDPType == "answer" {
		return d.SetRemoteDescription(remoteCandidates)
	}
	return errors.New("SelfSDPType in config: " + d.config.SelfSDPType + " can't set offer.")
}

func (d *DataChannel) SetAnswer(remoteCandidates []ICECandidate) error {
	if d.config.SelfSDPType == "offer" {
		return d.SetRemoteDescription(remoteCandidates)
	}
	return errors.New("SelfSDPType in config: " + d.config.SelfSDPType + " can't set answer.")
}

func (d *DataChannel) SetPort(port uint16) *DataChannel {
	if d.WebRTCSettingEngine.SetEphemeralUDPPortRange(port, port) != nil {
		return nil
	}
	return d
}

// ips: list of IP in string
// iptype: Host (if full 1-to-1 DNAT), Srflx (if behind a NAT)
func (d *DataChannel) SetIP(ips []string, iptype ICECandidateType) *DataChannel {
	switch iptype {
	case Host:
		d.WebRTCSettingEngine.SetNAT1To1IPs(ips, webrtc.ICECandidateTypeHost)
		break
	case Srflx:
		d.WebRTCSettingEngine.SetNAT1To1IPs(ips, webrtc.ICECandidateTypeSrflx)
		break
	default:
		return nil
	}
	return d
}

func (d *DataChannel) SetNetworkTypes(candidateTypes []webrtc.NetworkType) *DataChannel {
	d.WebRTCSettingEngine.SetNetworkTypes(candidateTypes)
	return d
}

// ReadyToSend() when Data Channal is opened and is not exceeding the bytes limit.
func (d *DataChannel) ReadyToSend() bool {
	return (d.WebRTCDataChannel.ReadyState() == webrtc.DataChannelStateOpen) && (d.WebRTCDataChannel.BufferedAmount() < d.config.TxBufferSize)
}

// Send []byte object via Data Channel
func (d *DataChannel) Send(data []byte) error {
	if d.WebRTCDataChannel.ReadyState() == webrtc.DataChannelStateOpen {
		if d.WebRTCDataChannel.BufferedAmount() >= DataChannelBufferBytesLim {
			return ErrDataChannelAtCapacity
		}
		return d.WebRTCDataChannel.Send(data)
	} else if d.WebRTCDataChannel.ReadyState() == webrtc.DataChannelStateConnecting {
		return ErrDatachannelNotReady
	} else {
		return ErrDataChannelClosed
	}
}

// Close() ends Data Channel and Peer Connection
func (d *DataChannel) Close() {
	d.WebRTCDataChannel.Close()
	d.WebRTCPeerConnection.Close()
}
