package gowork

import (
	"fmt"
	"strconv"
	"time"
	"log"
)

// TODO:Stan look at wrapping https://github.com/jinzhu/now.
var offset time.Duration = 0

func ClearOffset() {
	log.Printf("WARNING: ClearOffset")
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
	offset = date.Sub(CurrentTime())
	log.Printf("WARNING: SetOffsetTime %s (%s)", date, offset)
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

func CeilingDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 24, 0, -1, 0, t.Location()) //Ignore Nanoseconds
}

func FloorWeek(t time.Time) time.Time {
	for t.Weekday() != time.Sunday {
		t = t.Add(time.Hour * -24)
	}
	return FloorDay(t)
}

func CeilingWeek(t time.Time) time.Time {
	for t.Weekday() != time.Saturday {
		t = t.Add(time.Hour * 24)
	}
	return CeilingDay(t)
}

func FloorMonth(t time.Time) time.Time {
	return FloorDay(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()))
}

func CeilingMonth(t time.Time) time.Time {
	return CeilingDay(time.Date(t.Year(), t.Month() + 1, 0, 0, 0, 0, 0, t.Location()))
}

func FloorYear(t time.Time) time.Time {
	return FloorDay(time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()))
}

func CeilingYear(t time.Time) time.Time {
	return CeilingDay(time.Date(t.Year(), 12, 31, 0, 0, 0, 0, t.Location()))
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

	// JavaScript toISOString() format
	if dt, err := time.Parse("2006-01-02T15:04:05-07:00", datestr); err == nil {
		return dt.Local(), nil
	}

	// Another ISO8061 format
	if dt, err := time.Parse("2006-01-02 15:04:05-0700", datestr); err == nil {
		return dt.Local(), nil
	}

	// Another ISO8061 format without TZ
	if dt, err := time.ParseInLocation("2006-01-02 15:04:05", datestr, time.Local); err == nil {
		return dt, nil
	}

	// Default timeless dates to Local.
	if dt, err := time.ParseInLocation("2006-01-02", datestr, time.Local); err == nil {
		return dt, nil
	}

	return time.Time{}, fmt.Errorf("Could not parse datestr: %s (%s)", datestr, CurrentFunctionName(3))
}

// SubtractMonth subtracts one month from the provided value like AddDate, but unlike AddDate it will not normalize dates.
// i.e. With AddDate: "December 31st".AddDate(0, -1, 0) results in December 1st because there is no November 31st.
// With SubtractMonth "December 31st" becomes November 30.
func SubtractMonth(t time.Time) time.Time {
	nt := t.AddDate(0, -1, 0)
	if t.Month() == nt.Month() {
		nt = SubtractMonth(t.AddDate(0, 0, -1))
	}
	return nt
}