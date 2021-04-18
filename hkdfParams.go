package seed2sdp

const (
	MINSECRETLEN     int = 8
	MINSALTLEN       int = 8
	MININFOPREFIXLEN int = 8
)

type HKDFParams struct {
	secret     []byte
	salt       []byte
	infoPrefix []byte
}

func NewHKDFParams() *HKDFParams {
	// Add-on: Security measures
	return &HKDFParams{
		secret:     []byte("DefaultHKDFSecret"),
		salt:       []byte("DefaultHKDFSalt"),
		infoPrefix: []byte("DefaultHKDFPrefix"),
	}
}

func (p *HKDFParams) SetSecret(secret string) *HKDFParams {
	if p != nil {
		p.secret = []byte(secret)
	}
	return p
}

func (p *HKDFParams) SetSalt(salt string) *HKDFParams {
	if p != nil {
		p.salt = []byte(salt)
	}
	return p
}

func (p *HKDFParams) SetInfoPrefix(infoPrefix string) *HKDFParams {
	if p != nil {
		p.infoPrefix = []byte(infoPrefix)
	}
	return p
}

// A valid HKDFParams has all 3 []byte with proper length.
func (p *HKDFParams) IsValid() bool {
	return (len(p.secret) >= MINSECRETLEN && len(p.salt) >= MINSALTLEN && len(p.infoPrefix) >= MININFOPREFIXLEN)
}
