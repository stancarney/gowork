package gowork

import (
	"net/http"
	"strconv"
	"time"
	"log"
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

	if dt, err := UnMarshalDate(date); err == nil {
		return FloorDay(dt.Local())
	} else {
		log.Printf("Could not parse date: %s (%s)\n", err, CurrentFunctionName(3))
	}

	return FloorDay(CurrentTime())
}