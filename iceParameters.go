package seed2sdp

import (
	ice "github.com/Gaukas/ice_kai/v2"
)

type ICEParameters struct {
	UsernameFragment string `json:"usernameFragment"`
	Password         string `json:"password"`
	ICELite          bool   `json:"iceLite"` // Always false for now
}

func PredictIceParameters(Secret []byte, Salt []byte, InfoPrefix []byte) (ICEParameters, error) {
	iceParamArr, err := ice.GetHKDFUfragPwd(Secret, Salt, InfoPrefix)
	if err != nil {
		return ICEParameters{}, err
	}
	return ICEParameters{
		UsernameFragment: iceParamArr[0],
		Password:         iceParamArr[1],
		ICELite:          false,
	}, nil
}
