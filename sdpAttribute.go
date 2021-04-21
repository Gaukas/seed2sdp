package seed2sdp

import (
	"fmt"
)

// The lines that does not impact
type SDPAttribute struct {
	Key   string
	Value string
}

func (sa *SDPAttribute) String() string {
	if sa.Key != "" && sa.Value != "" {
		return fmt.Sprintf(`%s:%s`, sa.Key, sa.Value) // Key:Value pair
	} else {
		return fmt.Sprintf(`%s%s`, sa.Key, sa.Value) // either just sa.Key or sa.Value (as at least 1 empty)
	}
}
