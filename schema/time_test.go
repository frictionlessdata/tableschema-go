package schema

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestCastTime(t *testing.T) {
	t.Run("ValidMaximum", func(t *testing.T) {
		is := is.New(t)
		_, err := castTime(defaultFieldFormat, "11:45:00", Constraints{Maximum: "11:45:01"})
		is.NoErr(err)
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		is := is.New(t)
		_, err := castTime(defaultFieldFormat, "11:45:00", Constraints{Minimum: "11:44:59"})
		is.NoErr(err)
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc        string
			time        string
			constraints Constraints
		}{
			{"InvalidYear", "foo", Constraints{}},
			{"BiggerThanMaximum", "11:45:00", Constraints{Maximum: "11:44:59"}},
			{"InvalidMaximum", "11:45:00", Constraints{Maximum: "boo"}},
			{"SmallerThanMinimum", "11:45:00", Constraints{Minimum: "11:45:01"}},
			{"InvalidMinimum", "11:45:00", Constraints{Minimum: "boo"}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := castTime(defaultFieldFormat, d.time, d.constraints)
				is.True(err != nil)
			})
		}
	})
}

func TestUncastTime(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []struct {
			desc  string
			value time.Time
			want  string
		}{
			{"SimpleDate", time.Unix(1, 0), "1970-01-01T00:00:01Z"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				got, err := uncastTime(d.value)
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
			{"InvalidType", "Boo"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := uncastTime(d.value)
				is.True(err != nil)
			})
		}
	})
}
