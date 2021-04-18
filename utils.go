package seed2sdp

import "encoding/json"

func ToJSON(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func FromJSON(in string, obj interface{}) {
	err := json.Unmarshal([]byte(in), obj)
	if err != nil {
		panic(err)
	}
}
