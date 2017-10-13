package event

import (
	"time"
)

type timer struct {
	now    time.Time
	second *Hub
	minute *Hub
	day    *Hub
}

func newTimer() *timer {
	return &timer{
		now:    time.Now(),
		second: NewHub(LargeHub),
		minute: NewHub(MediumHub),
		day:    NewHub(SmallHub),
	}
}

func (t *timer) loop() {
	last := t.now
	now := time.Now()
	t.now = now
	t.second.Emit(now)
	if last.Minute() != now.Minute() {
		t.minute.Emit(now)
		// change day
		if last.Day() != now.Day() {
			t.day.Emit(now)
		}
	}
}

func (t *timer) UnixNow() int64 {
	return t.now.Unix()
}

func (t *timer) EverySecond() (Line, Cancel) {
	return t.second.NewLine()
}

func (t *timer) MinuteBegin() (Line, Cancel) {
	return t.minute.NewLine()
}

func (t *timer) DayBegin() (Line, Cancel) {
	return t.day.NewLine()
}

var Timer = newTimer()

func (t *timer) launch() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		t.loop()
	}
}

func init() {
	go Timer.launch()
}
