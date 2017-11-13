package query

import (
	"http/web"
	"time"
	"util/timer"
)

const dateFormat = "2006-01-02"

type DateRange struct {
	start timer.Day
	end   timer.Day
}

func (d *DateRange) StartUnix() int64 {
	return d.start.Begin().Unix()
}

func (d *DateRange) EndUnix() int64 {
	return d.end.End().Unix()
}

func (q Query) GetDay(key string) (*timer.Day, error) {
	value := q.Get(key)
	if value == "" {
		return nil, web.BadRequest("missing " + key)
	}
	ti, err := time.Parse(dateFormat, value)
	if err != nil {
		return nil, web.WrapBadRequest(err, "date format for "+key+" must be "+dateFormat)
	}
	d := timer.NewDay(ti)
	return &d, nil
}

func (q Query) GetDateRange(startKey string, endKey string) (*DateRange, error) {
	start, err := q.GetDay(startKey)
	if err != nil {
		return nil, err
	}
	end, err := q.GetDay(endKey)
	if err != nil {
		return nil, err
	}
	return &DateRange{start: *start, end: *end}, nil
}

func (q Query) MustGetDateRange(startKey string, endKey string) *DateRange {
	dateRange, err := q.GetDateRange(startKey, endKey)
	if err != nil {
		panic(err)
	}
	return dateRange
}
