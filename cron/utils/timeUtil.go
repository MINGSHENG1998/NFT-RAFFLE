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

func (s *TimeUtil) GetCurrentLocation() (*time.Location, error) {
	return time.LoadLocation(currentLocation)
}
