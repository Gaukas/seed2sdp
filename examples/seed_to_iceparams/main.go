package main

import (
	"fmt"

	s2s "github.com/Gaukas/seed2sdp"
)

func main() {
	hkdfParams, _ := s2s.NewHKDFParams([]byte("ExampleSecret"), []byte("ExampleSeed"), []byte("ExamplePrefix"))
	iceParams, _ := s2s.PredictIceParameters(hkdfParams)
	fmt.Println("ice-ufrag: ", iceParams.UsernameFragment)
	fmt.Println("ice-pwd: ", iceParams.Password)
	fmt.Println("ice-lite: ", iceParams.ICELite)
}
