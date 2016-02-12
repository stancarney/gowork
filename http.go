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

	if dt, err := time.Parse("2006-01-02T15:04:05-07:00", date); err == nil {
		local := dt.Local()
		
		// If somebody requests a given date in their TZ we don't want to roll the date ahead of forward to align it with the system TZ.
		if dt.Day() != local.Day() {
			local = local.AddDate(0, 0, dt.Day() - local.Day())
		}

		return FloorDay(local)
	} else {
		log.Printf("Could not parse date %s (%s)\n", err, GetCurrentFunctionName(3))
	}

	// Default timeless dates to Local.
	if dt, err := time.ParseInLocation("2006-01-02", date, time.Local); err == nil {
		return FloorDay(dt)
	} else {
		log.Printf("Could not parse date %s (%s)\n", err, GetCurrentFunctionName(3))
	}

	return FloorDay(CurrentTime())
}