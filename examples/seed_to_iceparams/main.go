package main

import (
	"fmt"

	"github.com/Gaukas/seed2sdp"
)

func main() {
	iceParams, _ := seed2sdp.PredictIceParameters([]byte("ExampleSecret"), []byte("ExampleSeed"), []byte("ExamplePrefix"))
	fmt.Println("ice-ufrag: ", iceParams.UsernameFragment)
	fmt.Println("ice-pwd: ", iceParams.Password)
	fmt.Println("ice-lite: ", iceParams.ICELite)
}
