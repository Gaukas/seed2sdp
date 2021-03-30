package main

import (
	"fmt"

	"github.com/Gaukas/seed2sdp"
)

func main() {
	fp, _ := seed2sdp.PredictDTLSFingerprint([]byte("ExampleSecret"), []byte("ExampleSeed"), []byte("ExamplePrefix"))
	fmt.Println("Fingerprint: ", fp)
}
