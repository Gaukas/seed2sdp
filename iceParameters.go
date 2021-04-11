package seed2sdp

import (
	"crypto/sha256"

	randutil "github.com/Gaukas/randutil_kai"
	"golang.org/x/crypto/hkdf"
)

// Copied from ice/rand.go
const (
	runesAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	runesDigit = "0123456789"
	// runesCandidateIDFoundation = runesAlpha + runesDigit + "+/"

	lenUFrag = 16
	lenPwd   = 32
)

type ICEParameters struct {
	UsernameFragment string // 16-char
	Password         string // 32-char
	ICELite          bool   // Always false for now
}

func GetUfrag(secret []byte, salt []byte, infoPrefix []byte) (string, error) {
	uFragReader := hkdf.New(sha256.New, Secret, Salt, append(InfoPrefix, []byte("IceUfrag")...))
	return randutil.GenerateReaderCryptoRandomString(lenUFrag, runesAlpha, uFragReader)
}

func GetPwd(secret []byte, salt []byte, infoPrefix []byte) (string, error) {
	pwdReader := hkdf.New(sha256.New, Secret, Salt, append(InfoPrefix, []byte("IcePwd")...))
	return randutil.GenerateReaderCryptoRandomString(lenPwd, runesAlpha, pwdReader)
}

func PredictIceParameters(secret []byte, salt []byte, infoPrefix []byte) (ICEParameters, error) {
	ufrag, err := GetUfrag(Secret, Salt, InfoPrefix)
	if err != nil {
		return ICEParameters{}, err
	}
	pwd, err := GetPwd(Secret, Salt, InfoPrefix)
	if err != nil {
		return ICEParameters{}, err
	}
	return ICEParameters{
		UsernameFragment: ufrag,
		Password:         pwd,
		ICELite:          false,
	}, nil
}
