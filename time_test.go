package gowork

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestFloorDay_Match(t *testing.T) {

	//setup
	today := time.Now()

	//execute
	r := FloorDay(today)

	//verify
	require.Equal(t, today.Year(), r.Year())
	require.Equal(t, today.Month(), r.Month())
	require.Equal(t, today.Day(), r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, today.Location(), r.Location())
}

func TestCeilingDay_Match(t *testing.T) {

	//setup
	today := time.Now()

	//execute
	r := CeilingDay(today)

	//verify
	require.Equal(t, today.Year(), r.Year())
	require.Equal(t, today.Month(), r.Month())
	require.Equal(t, today.Day(), r.Day())
	require.Equal(t, 23, r.Hour())
	require.Equal(t, 59, r.Minute())
	require.Equal(t, 59, r.Second())
	require.Equal(t, 999999999, r.Nanosecond())
	require.Equal(t, today.Location(), r.Location())
}

func TestDayBefore_Success(t *testing.T) {

	//setup
	today := time.Now()

	//execute
	r := DayBefore(today)

	//verify
	require.Equal(t, today.AddDate(0, 0, -1), r)
}

func TestDaysBefore_0Days(t *testing.T) {

	//setup
	today := time.Now()

	//execute
	r := DaysBefore(today, 0)

	//verify
	require.Equal(t, today, r)
}

func TestDaysBefore_NDays(t *testing.T) {

	//setup
	today := time.Now()

	//execute
	r := DaysBefore(today, 1)

	//verify
	require.Equal(t, today.AddDate(0, 0, -1), r)
}

func TestDaysBefore_1Week(t *testing.T) {

	//setup
	date := time.Date(2015, 3, 23, 0, 0, 0, 0, time.Local)

	//execute
	r := DaysBefore(date, 7)

	//verify
	require.Equal(t, 2015, r.Year())
	require.Equal(t, time.March, r.Month())
	require.Equal(t, 17, r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, date.Location(), r.Location())
}

func TestToJulianDate_Success(t *testing.T) {

	//setup
	date := time.Date(2015, 3, 23, 0, 0, 0, 0, time.Local)

	//execute
	r := ToJulianDate(date)

	//verify
	require.Equal(t, "15082", r)
}

func TestFromJulianDate_Success(t *testing.T) {

	//setup
	jd := "15082"

	//execute
	r := FromJulianDate(jd)

	//verify
	require.Equal(t, 2015, r.Year())
	require.Equal(t, 82, r.YearDay())
}

func TestFloorWeek_Success(t *testing.T) {

	//setup
	d, err := UnMarshalDate("2016-04-13 06:30:29")
	assert.NoError(t, err)

	//execute
	r := FloorWeek(d)

	//verify
	require.Equal(t, 2016, r.Year())
	require.Equal(t, time.April, r.Month())
	require.Equal(t, 10, r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, d.Location(), r.Location())
}

func TestCeilingWeek_Success(t *testing.T) {

	//setup
	d, err := UnMarshalDate("2016-04-16 06:30:29")
	assert.NoError(t, err)

	//execute
	r := CeilingWeek(d)

	//verify
	require.Equal(t, 2016, r.Year())
	require.Equal(t, time.April, r.Month())
	require.Equal(t, 16, r.Day())
	require.Equal(t, 23, r.Hour())
	require.Equal(t, 59, r.Minute())
	require.Equal(t, 59, r.Second())
	require.Equal(t, 999999999, r.Nanosecond())
	require.Equal(t, r.Location(), r.Location())
}

func TestCeilingMonth_Success(t *testing.T) {

	//setup
	d, err := UnMarshalDate("2016-04-16 06:30:29")
	assert.NoError(t, err)

	//execute
	r := CeilingMonth(d)

	//verify
	require.Equal(t, 2016, r.Year())
	require.Equal(t, time.April, r.Month())
	require.Equal(t, 30, r.Day())
	require.Equal(t, 23, r.Hour())
	require.Equal(t, 59, r.Minute())
	require.Equal(t, 59, r.Second())
	require.Equal(t, 999999999, r.Nanosecond())
	require.Equal(t, r.Location(), r.Location())
}

func TestFloorMonth_Success(t *testing.T) {

	//setup
	d, err := UnMarshalDate("2016-04-13 06:30:29")
	assert.NoError(t, err)

	//execute
	r := FloorMonth(d)

	//verify
	require.Equal(t, 2016, r.Year())
	require.Equal(t, time.April, r.Month())
	require.Equal(t, 1, r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, d.Location(), r.Location())
}

func TestCeilingYear_Success(t *testing.T) {

	//setup
	d, err := UnMarshalDate("2016-04-16 06:30:29")
	assert.NoError(t, err)

	//execute
	r := CeilingYear(d)

	//verify
	require.Equal(t, 2016, r.Year())
	require.Equal(t, time.December, r.Month())
	require.Equal(t, 31, r.Day())
	require.Equal(t, 23, r.Hour())
	require.Equal(t, 59, r.Minute())
	require.Equal(t, 59, r.Second())
	require.Equal(t, 999999999, r.Nanosecond())
	require.Equal(t, r.Location(), r.Location())
}

func TestFloorYear_Success(t *testing.T) {

	//setup
	d, err := UnMarshalDate("2016-04-13 06:30:29")
	assert.NoError(t, err)

	//execute
	r := FloorYear(d)

	//verify
	require.Equal(t, 2016, r.Year())
	require.Equal(t, time.January, r.Month())
	require.Equal(t, 1, r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, d.Location(), r.Location())
}

func TestCurrentTime_Success(t *testing.T) {

	//setup

	//execute
	r := CurrentTime()

	//verify
	require.NotNil(t, r)
	require.Equal(t, time.Now().Location(), r.Location())
}

func TestUnMarshalDate_ISO8061_Success(t *testing.T) {

	//setup
	date := CurrentTime()

	//execute
	r, err := UnMarshalDate(date.Format("2006-01-02 00:00:00-0700"))

	//verify
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Equal(t, date.Year(), r.Year())
	require.Equal(t, date.Month(), r.Month())
	require.Equal(t, date.Day(), r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, date.Location(), r.Location())
}

func TestUnMarshalDate_NO_TZ(t *testing.T) {

	//setup
	date := CurrentTime()

	//execute
	r, err := UnMarshalDate(date.Format("2006-01-02 00:00:00"))

	//verify
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Equal(t, date.Year(), r.Year())
	require.Equal(t, date.Month(), r.Month())
	require.Equal(t, date.Day(), r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, date.Location(), r.Location())
}

func TestUnMarshalDate_ISO8061_JS(t *testing.T) {

	//setup
	date := CurrentTime()

	//execute
	r, err := UnMarshalDate(date.Format("2006-01-02T00:00:00-07:00"))

	//verify
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Equal(t, date.Year(), r.Year())
	require.Equal(t, date.Month(), r.Month())
	require.Equal(t, date.Day(), r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, date.Location(), r.Location())
}

func TestUnMarshalDate_Timeless(t *testing.T) {

	//setup
	date := CurrentTime()

	//execute
	r, err := UnMarshalDate(date.Format("2006-01-02"))

	//verify
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Equal(t, date.Year(), r.Year())
	require.Equal(t, date.Month(), r.Month())
	require.Equal(t, date.Day(), r.Day())
	require.Equal(t, 0, r.Hour())
	require.Equal(t, 0, r.Minute())
	require.Equal(t, 0, r.Second())
	require.Equal(t, 0, r.Nanosecond())
	require.Equal(t, date.Location(), r.Location())
}

func TestUnMarshalDate_Invalid(t *testing.T) {

	//setup

	//execute
	_, err := UnMarshalDate("2015-AB-01")

	//verify
	require.NotNil(t, err)
	require.Equal(t, "Could not parse datestr: 2015-AB-01 (gowork.TestUnMarshalDate_Invalid)", err.Error())
}
