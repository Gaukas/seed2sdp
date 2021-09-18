package seed2sdp

import (
	"crypto/sha256"

	randutil "github.com/Gaukas/randutil_kai"
	"github.com/pion/webrtc/v3"
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

func GetUfrag(hkdfParams *HKDFParams) (string, error) {
	uFragReader := hkdf.New(sha256.New, hkdfParams.secret, hkdfParams.salt, append(hkdfParams.infoPrefix, []byte("IceUfrag")...))
	return randutil.GenerateReaderCryptoRandomString(lenUFrag, runesAlpha, uFragReader)
}

func GetPwd(hkdfParams *HKDFParams) (string, error) {
	pwdReader := hkdf.New(sha256.New, hkdfParams.secret, hkdfParams.salt, append(hkdfParams.infoPrefix, []byte("IcePwd")...))
	return randutil.GenerateReaderCryptoRandomString(lenPwd, runesAlpha, pwdReader)
}

func PredictIceParameters(hkdfParams *HKDFParams) (ICEParameters, error) {
	ufrag, err := GetUfrag(hkdfParams)
	if err != nil {
		return ICEParameters{}, err
	}
	pwd, err := GetPwd(hkdfParams)
	if err != nil {
		return ICEParameters{}, err
	}
	return ICEParameters{
		UsernameFragment: ufrag,
		Password:         pwd,
		ICELite:          false,
	}, nil
}

func (i ICEParameters) Equal(d ICEParameters) bool {
	return i.UsernameFragment == d.UsernameFragment && i.Password == d.Password && i.ICELite == d.ICELite
}

func (i *ICEParameters) UpdateSettingEngine(se *webrtc.SettingEngine) error {
	if i.UsernameFragment == "" || i.Password == "" {
		return ErrMalformedICEParameters
	}
	se.SetICECredentials(i.UsernameFragment, i.Password)
	return nil
}
