package xtime

import (
	"context"
	"database/sql/driver"
	"strconv"
	"time"

	"ascale/pkg/conf/env"

	"github.com/golang-module/carbon"

	"github.com/hooklift/gowsdl/soap"
)

var (
	delta time.Duration
)

func NowCarbon() carbon.Carbon {
	return carbon.Time2Carbon(Now())
}

func Now() time.Time {
	if env.DeployEnv == env.DeployEnvProd {
		return time.Now()
	} else {
		return time.Now().Add(delta)
	}
}

func NowUnix() int64 {
	if env.DeployEnv == env.DeployEnvProd {
		return time.Now().Unix()
	} else {
		return time.Now().Add(delta).Unix()
	}
}

func BeginingOfDayUnix(tz string) int64 {
	loc, _ := time.LoadLocation(tz)
	if env.DeployEnv == env.DeployEnvProd {
		now := time.Now().In(loc)
		year, month, day := now.Date()
		beginningOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		return beginningOfDay.Unix()
	} else {
		now := time.Now().Add(delta).In(loc)
		year, month, day := now.Date()
		beginningOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		return beginningOfDay.Unix()
	}
}

func BeginingOfMonthUnix(tz string) int64 {
	loc, _ := time.LoadLocation(tz)
	if env.DeployEnv == env.DeployEnvProd {
		now := time.Now().In(loc)
		year, month, _ := now.Date()
		beginningOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		return beginningOfMonth.Unix()
	} else {
		now := time.Now().Add(delta).In(loc)
		year, month, _ := now.Date()
		beginningOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		return beginningOfMonth.Unix()
	}
}

func EndOfMonthUnix(tz string) int64 {
	loc, _ := time.LoadLocation(tz)
	if env.DeployEnv == env.DeployEnvProd {
		now := time.Now().In(loc)
		year, month, _ := now.Date()
		endOfMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, now.Location())
		return endOfMonth.Unix() - 1
	} else {
		now := time.Now().Add(delta).In(loc)
		year, month, _ := now.Date()
		endOfMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, now.Location())
		return endOfMonth.Unix() - 1
	}
}

func BeginningOfDayUnixWithTimezone(timeUnix int64, tz string) int64 {
	loc, _ := time.LoadLocation(tz)
	locTime := time.Unix(timeUnix, 0).In(loc)
	year, month, day := locTime.Date()
	beginningOfDay := time.Date(year, month, day, 0, 0, 0, 0, locTime.Location())
	return beginningOfDay.Unix()
}

func ResetSystemTimeDelta() {
	if env.DeployEnv == env.DeployEnvProd {
		return
	}

	delta = 0
}

func SetSystemTimeDelta(dt time.Time) {
	if env.DeployEnv == env.DeployEnvProd {
		return
	}

	delta = dt.Sub(time.Now())
}

// Time be used to MySql timestamp converting.
type Time int64

// Scan scan time.
func (jt *Time) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*jt = Time(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = Time(i)
	}
	return
}

// Value get time value.
func (jt Time) Value() (driver.Value, error) {
	return time.Unix(int64(jt), 0), nil
}

// Time get time.
func (jt Time) Time() time.Time {
	return time.Unix(int64(jt), 0)
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration time.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < time.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, time.Duration(d))
	return d, ctx, cancel
}

// IsUnder18YearsOld check birthdate under 18 years old
func IsUnder18YearsOld(birthdate int64) bool {
	return birthdate > time.Now().AddDate(-18, 0, 0).Unix()
}

func RxntGoUnix(t soap.XSDDateTime) (dt int64) {
	t.StripTz()
	d := t.ToGoTime()

	loc, _ := time.LoadLocation("EST")
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), loc).Unix()
}

func UnixToTimeString(unixTime int64) (timeString string) {
	t := time.Unix(unixTime, 0)
	layout := "2006-01-02T15:04:05"
	return t.UTC().Format(layout)
}

func UnixToTimeStringWithTimezone(unixTime int64, timeZone string) (timeString string, err error) {
	t := time.Unix(unixTime, 0)

	// Load the time zone location
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return "", err
	}

	// Convert the time to the specified time zone
	t = t.In(loc)

	// Format the time
	formattedTime := t.Format("2006-01-02 15:04:05 MST")

	return formattedTime, nil
}

func DayStartOfUnixAndLoc(unixTime int64, loc *time.Location) (b time.Time) {
	d := time.Unix(unixTime, 0).In(loc)
	b = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
	return
}
