package schema

import (
	"testing"
	"time"
)

func TestCastDuration_Success(t *testing.T) {
	data := []struct {
		desc  string
		value string
		want  time.Duration
	}{
		{"OnlyP", "P", time.Duration(0)},
		{"OnlyHour", "P2H", time.Duration(2) * time.Hour},
		{"HourDefaultZero", "PH", time.Duration(0) * time.Hour},
		{"OnlyPeriod", "P3Y6M4D", 3*hoursInYear + 6*hoursInMonth + 4*hoursInDay},
		{"OnlyTime", "PT12H30M5S", 12*time.Hour + 30*time.Minute + 5*time.Second},
		{"Complex", "P3Y6M4DT12H30M5S", 3*hoursInYear + 6*hoursInMonth + 4*hoursInDay + 12*time.Hour + 30*time.Minute + 5*time.Second},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			got, err := castDuration(d.value)
			if err != nil {
				t.Errorf("want:nil got:%q", err)
			}
			if got != d.want {
				t.Errorf("want:%s got:%s", d.want, got)
			}
		})
	}
}

func TestCastDuration_Error(t *testing.T) {
	data := []struct {
		desc  string
		value string
	}{
		{"WrongStartChar", "C2H"},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			_, err := castDuration(d.value)
			if err == nil {
				t.Errorf("want:err got:nil")
			}
		})
	}
}
