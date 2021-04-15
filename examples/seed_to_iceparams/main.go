package main

import (
	"fmt"

	s2s "github.com/Gaukas/seed2sdp"
)

func main() {
	hkdfParams := s2s.NewHKDFParams().SetSecret("ExampleSecret").SetSalt("ExampleSeed").SetInfoPrefix("ExamplePrefix")
	iceParams, _ := s2s.PredictIceParameters(hkdfParams)
	fmt.Println("ice-ufrag: ", iceParams.UsernameFragment)
	fmt.Println("ice-pwd: ", iceParams.Password)
	fmt.Println("ice-lite: ", iceParams.ICELite)
}
