package gowork

import (
	"math/rand"
	"strings"
	"time"
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

