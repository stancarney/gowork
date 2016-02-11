package gowork

import (
	"fmt"
	"strconv"
	"time"
	"log"
)

var offset time.Duration = 0

func ClearOffset() {
	offset = 0
}

// SetOffset is only used for testing. It is used when you want to move the state of the application ahead of forward by the provided duration
// in order to test time sensitive code. Jobs, etc...
func SetOffset(val time.Duration) {
	log.Printf("WARNING: SetOffset %s", val)
	offset = val
}

// SetOffsetTime is only used for testing. It is used when you want to move the state of the application ahead of forward by the provided duration
// in order to test time sensitive code. Jobs, etc...
func SetOffsetTime(date time.Time) {
	log.Printf("WARNING: SetOffsetTime %s", date)
	offset = time.Since(date)
}

// CurrentTime should be used throughout the entire system when the current time is needed. This allows the testing infrastructure to use offsets.
func CurrentTime() time.Time {
	return time.Now().Add(offset)
}

func DayBefore(t time.Time) time.Time {
	return t.Add(time.Hour * -24)
}

func DaysBefore(t time.Time, days int) time.Time {
	return t.Add(time.Hour * time.Duration(-24 * days))
}

func FloorDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func ToJulianDate(t time.Time) (date string) {
	year := strconv.Itoa(t.Year())
	date = fmt.Sprintf("%s%03v", year[2:len(year)], t.YearDay())
	return
}

func FromJulianDate(str string) time.Time {
	year, _ := strconv.Atoi(fmt.Sprintf("20%v", str[:2])) //TODO:Stan this will break long after I'm dead.
	days, _ := strconv.Atoi(string(str[2:]))

	date := time.Date(year, 1, 0, 0, 0, 0, 0, time.Local)
	return date.AddDate(0, 0, days)
}

func MarshalDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05-0700")
}

func UnMarshalDate(datestr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05-0700", datestr, time.Local)
}