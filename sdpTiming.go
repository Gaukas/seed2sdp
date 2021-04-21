package seed2sdp

import (
	"fmt"
)

// The lines that does not impact
type SDPTiming struct {
	Start uint32
	End   uint32
}

func NewSDPTiming() SDPTiming {
	return SDPTiming{
		Start: 0,
		End:   0,
	}
}

func (st *SDPTiming) String() string {
	return fmt.Sprintf(`%d %d`, st.Start, st.End)
}
