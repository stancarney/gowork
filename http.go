package gowork

import (
	"net/http"
	"strconv"
	"time"
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

func GetDate(r *http.Request) time.Time {
	vals := r.URL.Query()
	date := vals.Get("date")

	if dt, err := time.Parse("2006-01-02T15:04:05-07:00", date); err == nil {
		return FloorDay(dt)
	}

	if dt, err := time.ParseInLocation("2006-01-02", date, time.Local); err == nil {
		return FloorDay(dt)
	}

	return FloorDay(CurrentTime())
}