package utils

import (
	"sync"
	"time"
)

const (
	currentLocation string = "Asia/Singapore"
)

var (
	timeUtil     *TimeUtil
	timeUtilOnce sync.Once
)

type TimeUtil struct{}

func GetTimeUtil() *TimeUtil {
	if timeUtil == nil {
		timeUtilOnce.Do(func() {
			timeUtil = &TimeUtil{}
		})
	}
	return timeUtil
}

func (u *TimeUtil) GetCurrentLocation() (*time.Location, error) {
	return time.LoadLocation(currentLocation)
}

func (u *TimeUtil) Now() (int64, error) {
	loc, err := time.LoadLocation(currentLocation)

	if err != nil {
		return 0, err
	}

	return time.Now().In(loc).Unix(), nil
}
