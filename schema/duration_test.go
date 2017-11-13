package schema

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestCastDuration_Success(t *testing.T) {
	data := []struct {
		desc  string
		value string
		want  time.Duration
	}{
		{"OnlyP", "P", time.Duration(0)},
		{"OnlyHour", "P2H", time.Duration(2) * time.Hour},
		{"SecondsWithDecimal", "P22.519S", 22519 * time.Millisecond},
		{"HourDefaultZero", "PH", time.Duration(0) * time.Hour},
		{"OnlyPeriod", "P3Y6M4D", 3*hoursInYear + 6*hoursInMonth + 4*hoursInDay},
		{"OnlyTime", "PT12H30M5S", 12*time.Hour + 30*time.Minute + 5*time.Second},
		{"Complex", "P3Y6M4DT12H30M5S", 3*hoursInYear + 6*hoursInMonth + 4*hoursInDay + 12*time.Hour + 30*time.Minute + 5*time.Second},
		{"2Years", "P2Y", (2 * 360 * 24) * time.Hour},
		{"StringFieldsAreIgnored", "PfooHdddS", time.Duration(0)},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			is := is.New(t)
			got, err := castDuration(d.value)
			is.NoErr(err)
			is.Equal(got, d.want)
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
			is := is.New(t)
			_, err := castDuration(d.value)
			is.True(err != nil)
		})
	}
}

func TestUncastDuration(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []struct {
			desc  string
			value time.Duration
			want  string
		}{
			{"1Year", 1*hoursInYear + 1*hoursInMonth + 1*hoursInDay + 1*time.Hour + 1*time.Minute + 500*time.Millisecond, "P1Y1M1DT1H1M0.5S"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				got, err := uncastDuration(d.value)
				is.NoErr(err)
				is.Equal(d.want, got)
			})
		}
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc  string
			value interface{}
		}{
			{"InvalidType", 10},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := uncastDuration(d.value)
				is.True(err != nil)
			})
		}
	})
}
