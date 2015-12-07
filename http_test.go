package gowork

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestGetLimit_Valid(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost?num=10")
	request := http.Request{URL: u}

	//execute
	r := GetLimit(&request)

	//verify
	require.Equal(t, 10, r)
}

func TestGetLimit_Nothing(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost")
	request := http.Request{URL: u}

	//execute
	r := GetLimit(&request)

	//verify
	require.Equal(t, -1, r)
}

func TestGetLimit_Invalid(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost?num=aa")
	request := http.Request{URL: u}

	//execute
	r := GetLimit(&request)

	//verify
	require.Equal(t, -1, r)
}

func TestGetDate_Valid(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost?date=2015-01-01")
	request := http.Request{URL: u}

	//execute
	r := GetDate(&request)

	//verify
	require.Equal(t, "2015-01-01", r)
}

func TestGetDate_Nothing(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost")
	request := http.Request{URL: u}

	//execute
	r := GetDate(&request)

	//verify
	require.Equal(t, MarshalDate(CurrentTime()), r)
}

func TestGetIp_XForwardedFor(t *testing.T) {

	u, _ := url.ParseRequestURI("http://localhost")
	h := http.Header{}
	h.Add("X-Forwarded-For", "10.0.0.1")
	request := http.Request{URL: u, Header: h, RemoteAddr: "10.0.0.2:999"}

	//execute
	r := GetIp(&request)

	//verify
	require.Equal(t, "10.0.0.1", r)
}

func TestGetIp_XRealIp(t *testing.T) {

	u, _ := url.ParseRequestURI("http://localhost")
	h := http.Header{}
	h.Add("X-Real-Ip", "10.0.0.1")
	request := http.Request{URL: u, Header: h, RemoteAddr: "10.0.0.2:999"}

	//execute
	r := GetIp(&request)

	//verify
	require.Equal(t, "10.0.0.1", r)
}

func TestGetIp_Default(t *testing.T) {

	u, _ := url.ParseRequestURI("http://localhost")
	h := http.Header{}
	request := http.Request{URL: u, Header: h, RemoteAddr: "10.0.0.2:999"}

	//execute
	r := GetIp(&request)

	//verify
	require.Equal(t, "10.0.0.2", r)
}

