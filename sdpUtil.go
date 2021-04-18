package seed2sdp

import (
	"fmt"
	"net"
)

type SdpNetworkType string
type SdpIpaddrType string

const (
	IN SdpNetworkType = "IN"
)

const (
	IP4 SdpIpaddrType = "IP4"
	IP6 SdpIpaddrType = "IP6"
)

// Global Lines
type SdpGlobal struct {
	SessionId    uint64
	SessionVer   uint32
	NetworkType  SdpNetworkType
	IpaddrType   SdpIpaddrType
	UnicastAddr  net.IP
	SessionName  string
	StartingTime uint64
	EndingTime   uint64
	GroupBundle  []string
	Payload      string // Everything that is yet to be supported
}

func (sg *SdpGlobal) String() string {
	strsg := fmt.Sprintf(`o=- %d %d %s %s %s\r\n`, sg.SessionId, sg.SessionVer, sg.NetworkType, sg.IpaddrType, sg.UnicastAddr.String())
	strsg += fmt.Sprintf(`s=- %s\r\n`, sg.SessionName)
	strsg += fmt.Sprintf(`t=%d %d\r\n`, sg.StartingTime, sg.EndingTime)
	if len(sg.GroupBundle) > 0 {
		strsg += fmt.Sprintf("a=group:BUNDLE")
		for _, b := range sg.GroupBundle {
			strsg += fmt.Sprintf(" %s", b)
		}
		strsg += `\r\n`
	}
	strsg += sg.Payload
	return strsg
}
