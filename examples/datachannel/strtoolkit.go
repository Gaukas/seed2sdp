package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	randutil "github.com/Gaukas/randutil_kai"
)

func RandSeq(n int) string {
	val, err := randutil.GenerateCryptoRandomString(n, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if err != nil {
		panic(err)
	}

	return val
}

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

func MustReadStdin() string {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				panic(err)
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}

	fmt.Println("")

	return in
}
