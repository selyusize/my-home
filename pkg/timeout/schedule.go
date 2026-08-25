package timeout

import (
	"time"
	_ "time/tzdata"

	"github.com/selyusize/my-home/pkg/bitrix"
)

const (
	openHour       = 9
	closeHour      = 18
	khabarovskZone = "Asia/Vladivostok"
)

func khabarovsk() (*time.Location, error) {
	return time.LoadLocation(khabarovskZone)
}

func inOpenWindow(now time.Time, loc *time.Location) bool {
	hour := now.In(loc).Hour()
	return hour >= openHour && hour < closeHour
}

func pastClose(now time.Time, loc *time.Location) bool {
	return now.In(loc).Hour() >= closeHour
}

func nextOpen(now time.Time, loc *time.Location) time.Time {
	return nextAt(now, loc, openHour)
}

func nextClose(now time.Time, loc *time.Location) time.Time {
	return nextAt(now, loc, closeHour)
}

func nextAt(now time.Time, loc *time.Location, hour int) time.Time {
	local := now.In(loc)
	today := time.Date(local.Year(), local.Month(), local.Day(), hour, 0, 0, 0, loc)
	if local.Before(today) {
		return today
	}
	return today.Add(24 * time.Hour)
}

func needsOpen(status bitrix.TimeManStatus) bool {
	switch status {
	case bitrix.TimeManClosed, bitrix.TimeManExpired:
		return true
	default:
		return false
	}
}

func needsClose(status bitrix.TimeManStatus) bool {
	switch status {
	case bitrix.TimeManOpened, bitrix.TimeManPaused, bitrix.TimeManExpired:
		return true
	default:
		return false
	}
}
