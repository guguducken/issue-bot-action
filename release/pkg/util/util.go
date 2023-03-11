package util

import "time"

func ExpiredCheck(CreatedAt, UpdatedAt time.Time, expired int64) bool {
	loc, _ := time.LoadLocation(`UTC`)
	now := time.Now().In(loc).UnixMilli()
	CreatedAt = CreatedAt.In(loc)
	UpdatedAt = UpdatedAt.In(loc)
	lastTime := CreatedAt.UnixMilli()
	if UpdatedAt.UnixMilli() > lastTime {
		lastTime = UpdatedAt.UnixMilli()
	}
	if now-lastTime >= expired {
		return true
	}
	return false
}
