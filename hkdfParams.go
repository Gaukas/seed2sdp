package seed2sdp

type HKDFParams struct {
	secret     []byte
	salt       []byte
	infoPrefix []byte
}

func NewHKDFParams(secret []byte, salt []byte, infoPrefix []byte) (*HKDFParams, error) {
	// TODO: Add security measures
	return &HKDFParams{
		secret:     secret,
		salt:       salt,
		infoPrefix: infoPrefix,
	}, nil
}
