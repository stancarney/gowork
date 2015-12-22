package gowork

import (
	"net/http"
	"strconv"
	"net"
)

func GetLimit(r *http.Request) (limit int) {
	vals := r.URL.Query()
	num := vals.Get("num")
	if num != "" {
		limit, err := strconv.Atoi(num)
		if err != nil {
			return 0
		}
		return limit
	}
	return 0
}

func GetDate(r *http.Request) (date string) {
	vals := r.URL.Query()
	date = vals.Get("date")
	if date == "" {
		date = MarshalDate(CurrentTime())
		return
	}
	return
}

func GetIp(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

