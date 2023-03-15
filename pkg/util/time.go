package util

import (
	"time"
)

var loc *time.Location
var day_one int64 = 86400000 //ms

func init() {
	loc = time.FixedZone(`UTC`, 0)
}

func ParseTime(t string) (time.Time, error) {
	ti, err := time.ParseInLocation(time.RFC3339, t, loc)
	if err != nil {
		return ti, err
	}
	return ti, nil
}

func GetNow() int64 {
	return time.Now().In(loc).UnixMilli()
}

type Time_m struct {
	Days            int64
	Hours           int64
	Minute          int64
	Second          int64
	MillSecond      int64
	TotalMillSecond int64
}

func MStoTime(pass int64) (Time_m, error) {
	ti := Time_m{}
	if pass < 0 {
		return ti, EndLittleStart
	}
	ti.TotalMillSecond = pass
	ti.MillSecond = pass % 1000
	pass /= 1000

	ti.Second = pass % 60
	pass /= 60

	ti.Minute = pass % 60
	pass /= 60

	ti.Hours = pass % 24
	pass /= 24

	ti.Days = pass

	return ti, nil
}

func GetPassedTimeWithoutWeekend(t time.Time) (work_t Time_m, holiday_t Time_m, err error) {
	ti, mill_ti := t, t.UnixMilli()

	now := time.Now().In(loc)
	mill_now := now.UnixMilli()
	pass, err := MStoTime(mill_now - mill_ti)
	if err != nil {
		return
	}
	var (
		work    int64 = 0
		holiday int64 = 0
	)

	weeks := pass.Days / 7
	if weeks >= 1 {
		holiday += weeks * 2 * day_one
		work += weeks * 5 * day_one
	}

	day_start := int64(ti.Weekday())
	if day_start == 0 {
		day_start = 7
	}
	day_end := int64(now.Weekday())
	if day_end == 0 {
		day_end = 7
	}

	dura_one_start := (day_start-1)*day_one + mill_ti - time.Date(ti.Year(), ti.Month(), ti.Day(), 0, 0, 0, 0, loc).UnixMilli()
	dura_one_end := (day_end-1)*day_one + mill_now - time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).UnixMilli()

	if dura_one_end > dura_one_start {
		if day_end >= 6 {
			if day_start >= 6 {
				holiday += dura_one_end - dura_one_start
			} else {
				holiday += dura_one_end - day_one*5
				work += day_one*5 - dura_one_start
			}
		} else {
			work += dura_one_end - dura_one_start
		}
	} else {
		work_start := day_one*5 - dura_one_start
		if work_start > 0 {
			work += work_start
			dura_one_start = day_one * 5
		}

		holiday += day_one*7 - dura_one_start

		holiday_end := dura_one_end - day_one*5
		if holiday_end >= 0 {
			holiday += holiday_end
			dura_one_end = day_one * 5
		}
		work += dura_one_end
	}
	work_t, err = MStoTime(work)
	if err != nil {
		return
	}
	holiday_t, err = MStoTime(holiday)
	return
}

func GetPassedTimeWithoutWeekend_String(t string) (Time_m, Time_m, error) {
	ti, err := ParseTime(t)
	if err != nil {
		return Time_m{}, Time_m{}, err
	}
	return GetPassedTimeWithoutWeekend(ti)
}
