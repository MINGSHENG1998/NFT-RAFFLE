package helpers

import "time"

const (
	currentLocation string = "Asia/Singapore"
	localLocation   string = "Local"
	dateTimeFormat  string = "2006-01-02 15:04:05"
)

var (
	TimeHelper ITimeHelper = NewTimeHelper()
)

type ITimeHelper interface {
	GetCurrentLocationTime() (time.Time, error)
	GetCurrentLocationTimeWithAdditionalDuration(d time.Duration) (time.Time, error)
	GetCurrentLocationDate() (int, time.Month, int, error)
	ConvertTimeToCurrentLocationTime(timeToConvert time.Time) (time.Time, error)
	ConvertDateTimeStringToCurrentLocationTime(s string) (time.Time, error)
}

type timeHelperStruct struct{}

func NewTimeHelper() ITimeHelper {
	return &timeHelperStruct{}
}

func (t *timeHelperStruct) GetCurrentLocationTime() (time.Time, error) {
	loc, err := time.LoadLocation(currentLocation)

	if err != nil {
		return time.Now(), err
	}

	currentTime, err := time.ParseInLocation(time.RFC3339, time.Now().In(loc).Format(time.RFC3339), loc)

	if err != nil {
		return time.Now(), err
	}

	return currentTime, nil
}

func (t *timeHelperStruct) GetCurrentLocationTimeWithAdditionalDuration(d time.Duration) (time.Time, error) {
	loc, err := time.LoadLocation(currentLocation)

	if err != nil {
		return time.Now(), err
	}

	newTime, err := time.ParseInLocation(time.RFC3339, time.Now().In(loc).Add(d).Format(time.RFC3339), loc)

	if err != nil {
		return time.Now(), err
	}

	return newTime, nil
}

func (t *timeHelperStruct) GetCurrentLocationDate() (int, time.Month, int, error) {
	var year, day int
	var month time.Month
	loc, err := time.LoadLocation(currentLocation)

	if err != nil {
		year, month, day = time.Now().Local().Date()
		return year, month, day, err
	}

	year, month, day = time.Now().In(loc).Date()
	return year, month, day, nil
}

func (t *timeHelperStruct) ConvertTimeToCurrentLocationTime(timeToConvert time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(currentLocation)

	if err != nil {
		return timeToConvert, err
	}

	return timeToConvert.In(loc), nil
}

func (t *timeHelperStruct) ConvertDateTimeStringToCurrentLocationTime(s string) (time.Time, error) {
	loc, err := time.LoadLocation(currentLocation)

	if err != nil {
		return time.Now().Local(), err
	}

	dateTime, err := time.ParseInLocation(dateTimeFormat, s, loc)

	if err != nil {
		return time.Now().In(loc), err
	}

	return dateTime, nil
}
