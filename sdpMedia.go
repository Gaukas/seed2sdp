package seed2sdp

import (
	"fmt"
)

// The lines that does not impact
type SDPMedia struct {
	MediaType   string
	Description string
}

func (sm *SDPMedia) String() string {
	return fmt.Sprintf(`%s %s`, sm.MediaType, sm.Description)
}
