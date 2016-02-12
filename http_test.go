package gowork

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"time"
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
	require.Equal(t, 0, r)
}

func TestGetLimit_Invalid(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost?num=aa")
	request := http.Request{URL: u}

	//execute
	r := GetLimit(&request)

	//verify
	require.Equal(t, 0, r)
}

func TestGetDate_Valid_DateOnly(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost?date=2015-01-01")
	request := http.Request{URL: u}

	date, _ := time.ParseInLocation("2006-01-02", "2015-01-01", time.Local)

	//execute
	r := GetDate(&request)

	//verify
	require.Equal(t, FloorDay(date), r)
}

func TestGetDate_Valid_ISO8061(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost?date=2015-01-01T15:00:00-07:00")
	request := http.Request{URL: u}

	date, _ := time.ParseInLocation("2006-01-02", "2015-01-01", time.Local)

	//execute
	r := GetDate(&request)

	//verify
	require.Equal(t, FloorDay(date), r)
}

func TestGetDate_Valid_ISO8061_Different_TZ(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost?date=2015-12-30T00:00:00-05:00") // EST
	request := http.Request{URL: u}
	
	date, _ := time.ParseInLocation("2006-01-02T15:04:05", "2015-12-30T00:00:00", time.Local) //This allows for this test to be run on desktops with different TZ's.

	//execute
	r := GetDate(&request)

	//verify
	require.Equal(t, FloorDay(date), r)
}

func TestGetDate_Valid_ISO8061_Different_Day_TZ(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost?date=2015-12-30T00:00:00%2B14:00") // Kiribati TZ
	request := http.Request{URL: u}

	date, _ := time.ParseInLocation("2006-01-02T15:04:05", "2015-12-30T00:00:00", time.Local) //This allows for this test to be run on desktops with different TZ's.

	//execute
	r := GetDate(&request)

	//verify
	require.Equal(t, FloorDay(date), r)
}

func TestGetDate_Nothing(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost")
	request := http.Request{URL: u}

	//execute
	r := GetDate(&request)

	//verify
	require.Equal(t, FloorDay(CurrentTime()), r)
}
