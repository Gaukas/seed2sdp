package main

import (
	"fmt"

	s2s "github.com/Gaukas/seed2sdp"
)

func main() {
	hkdfParams := s2s.NewHKDFParams().SetSecret("ExampleSecret").SetSalt("ExampleSeed").SetInfoPrefix("ExamplePrefix")
	fp, _ := s2s.PredictDTLSFingerprint(hkdfParams)
	fmt.Println("Fingerprint: ", fp)
}
