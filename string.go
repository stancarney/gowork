package gowork

import (
	"math/rand"
	"strings"
	"time"
	"io"
	"io/ioutil"
	"bytes"
	"log"
)

var Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func StringPrefixMatch(full string, prefixes []string) bool {
	for _, str := range prefixes {
		if strings.HasPrefix(full, str) {
			return true
		}
	}
	return false
}

func RandomStringSelection(num int, values []rune) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, num)
	for i := range b {
		b[i] = values[rand.Intn(len(values))]
	}
	return string(b)
}

func RandomString(num int) string {
	return RandomStringSelection(num, Letters)
}

func TrimAllSpace(str []string) {
	for i, s := range str {
		str[i] = strings.TrimSpace(s)
	}
}

func NewReadCloserFromString(string string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(string))
}

func StringFromReadCloser(r io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	if err := r.Close(); err != nil {
		log.Println("Could not close (%s)", CurrentFunctionName(2))
	}
	return buf.String()
}