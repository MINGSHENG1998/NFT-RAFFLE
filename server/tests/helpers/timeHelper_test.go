package tests_helpers

import (
	"nft-raffle/helpers"
	"testing"
	"time"
)

var (
	timeHelper helpers.ITimeHelper = helpers.TimeHelper
)

func TestGetCurrentTimeSingapore(t *testing.T) {
	time, err := timeHelper.GetCurrentTimeSingapore()

	if err != nil {
		t.Error(err.Error())
	}

	t.Log(time)
}

func TestGetCurrentTimeSingaporeWithAdditionalDuration(t *testing.T) {
	time, err := timeHelper.GetCurrentTimeSingaporeWithAdditionalDuration(time.Hour * time.Duration(1))

	if err != nil {
		t.Error(err.Error())
	}

	t.Log(time)
}

func TestGetCurrentDateSingapore(t *testing.T) {
	year, month, day, err := timeHelper.GetCurrentDateSingapore()

	if err != nil {
		t.Error(err.Error())
	}

	t.Log(year, month, day)
}

func TestConvertTimeToSingaporeTime(t *testing.T) {
	utcTime, err := time.Parse(time.RFC3339, "2022-04-25T12:44:24.000+00:00")

	if err != nil {
		t.Error(err.Error())
	}

	singaporeTime, err := timeHelper.ConvertTimeToSingaporeTime(utcTime)

	if err != nil {
		t.Error(err.Error())
	}

	if singaporeTime.Format(time.RFC3339) != "2022-04-25T20:44:24+08:00" {
		t.Error("not matching")
	}

	t.Log(singaporeTime)
}

func TestConvertDateStringToSingaporeDate(t *testing.T) {
	s := "2022-10-30 23:59:59"

	dateTime, err := timeHelper.ConvertDateTimeStringToSingaporeTime(s)

	if err != nil {
		t.Error(err.Error())
	}

	if dateTime.Format(time.RFC3339) != "2022-10-30T23:59:59+08:00" {
		t.Error("not matching")
	}
}
