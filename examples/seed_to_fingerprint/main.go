package main

import (
	"fmt"

	s2s "github.com/Gaukas/seed2sdp"
)

func main() {
	hkdfParams, _ := s2s.NewHKDFParams([]byte("ExampleSecret"), []byte("ExampleSeed"), []byte("ExamplePrefix"))
	fp, _ := s2s.PredictDTLSFingerprint(hkdfParams)
	fmt.Println("Fingerprint: ", fp)
}
