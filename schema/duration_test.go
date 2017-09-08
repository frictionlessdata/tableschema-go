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

func TestEncodeDuration(t *testing.T) {
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
				got, err := encodeDuration(d.value)
				if err != nil {
					t.Fatalf("err want:nil got:%q", err)
				}
				if d.want != got {
					t.Errorf("val want:%s got:%s", d.want, got)
				}
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
				_, err := encodeDuration(d.value)
				if err == nil {
					t.Fatalf("err want:err got:nil")
				}
			})
		}
	})
}
