package timer

import (
	"time"
)

const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * 60
	secondsPerDay    = 24 * secondsPerHour
	secondsPerWeek   = 7 * secondsPerDay
	daysPer400Years  = 365*400 + 97
	daysPer100Years  = 365*100 + 24
	daysPer4Years    = 365*4 + 1
)

var offset int64

func init() {
	var _, d = time.Now().Zone()
	offset = int64(d)
}

func LocalStartOfToday() int64 {
	var times int64
	var nowGTM = time.Now().Unix()
	var nowLocal = time.Now().Unix() - offset
	if timeToDay(nowGTM) != timeToDay(nowLocal) {
		if offset > 0 {
			var ellapsed = nowGTM % secondsPerDay
			times = nowGTM - ((ellapsed + offset) % secondsPerDay) + secondsPerDay
		} else {
			var ellapsed = nowGTM % secondsPerDay
			times = nowGTM - ((ellapsed + offset) % secondsPerDay) - secondsPerDay
		}

	} else {
		var ellapsed = nowGTM % secondsPerDay
		times = nowGTM - ((ellapsed + offset) % secondsPerDay)
	}
	return times

}

func TimeZone() int64 {
	_, t := time.Now().Zone()
	return int64(t)
}
func timeToDay(ctime int64) string {
	var t = time.Unix(ctime, 0)
	return t.Format("2006-01-02")
}

type Day struct {
	t time.Time
}

func NewDay(t time.Time) Day {
	return Day{t}
}

func DayStart(t time.Time) time.Time {
	yy, mm, dd := t.Date()
	t = time.Date(yy, mm, dd, 0, 0, 0, 0, t.Location())
	return t
}

const dayFormat = "2006-01-02"

func (d Day) String() string {
	return d.t.Format(dayFormat)
}

func (d Day) Begin() time.Time {
	return d.t
}

func (d Day) End() time.Time {
	return d.Begin().Add(time.Hour*24 - 1)
}
