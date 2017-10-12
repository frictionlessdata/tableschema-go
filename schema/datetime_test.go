package schema

import (
	"testing"

	"github.com/matryer/is"
)

func TestDecodeDatetime(t *testing.T) {
	t.Run("ValidMaximum", func(t *testing.T) {
		is := is.New(t)
		_, err := decodeDateTime("2013-01-24T22:01:00+07:00", Constraints{Maximum: "2014-01-24T22:01:00Z"})
		is.NoErr(err)
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		is := is.New(t)
		_, err := decodeDateTime("2013-01-24T22:01:00Z", Constraints{Minimum: "2012-01-24T22:01:00Z"})
		is.NoErr(err)
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc        string
			datetime    string
			constraints Constraints
		}{
			{
				"InvalidDateTime",
				"foo",
				Constraints{},
			},
			{
				"DateTimeBiggerThanMaximum",
				"2013-01-24T22:01:00Z",
				Constraints{Maximum: "2013-01-24T01:01:00Z"},
			},
			{
				"InvalidMaximum",
				"2013-01-24T22:01:00Z",
				Constraints{Maximum: "boo"},
			},
			{
				"DateTimeSmallerThanMinimum",
				"2013-01-24T22:01:00Z",
				Constraints{Minimum: "2013-01-24T22:01:01Z"},
			},
			{
				"InvalidMinimum",
				"2013-01-24T22:01:00Z",
				Constraints{Minimum: "boo"},
			},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := decodeDateTime(d.datetime, d.constraints)
				is.True(err != nil)
			})
		}
	})
}

func TestDecodeYear(t *testing.T) {
	t.Run("ValidMaximum", func(t *testing.T) {
		is := is.New(t)
		_, err := decodeYear("2006", Constraints{Maximum: "2007"})
		is.NoErr(err)
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		is := is.New(t)
		_, err := decodeYear("2007", Constraints{Minimum: "2006"})
		is.NoErr(err)
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc        string
			year        string
			constraints Constraints
		}{
			{"InvalidYear", "foo", Constraints{}},
			{"YearBiggerThanMaximum", "2006", Constraints{Maximum: "2005"}},
			{"InvalidMaximum", "2005", Constraints{Maximum: "boo"}},
			{"YearSmallerThanMinimum", "2005", Constraints{Minimum: "2006"}},
			{"InvalidMinimum", "2005", Constraints{Minimum: "boo"}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := decodeYear(d.year, d.constraints)
				is.True(err != nil)
			})
		}
	})
}

func TestDecodeYearMonth(t *testing.T) {
	t.Run("ValidMaximum", func(t *testing.T) {
		is := is.New(t)
		_, err := decodeYearMonth("2006-02", Constraints{Maximum: "2006-03"})
		is.NoErr(err)
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		is := is.New(t)
		_, err := decodeYearMonth("2006-03", Constraints{Minimum: "2006-02"})
		is.NoErr(err)
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc        string
			year        string
			constraints Constraints
		}{
			{"InvalidYear", "foo", Constraints{}},
			{"YearBiggerThanMaximum", "2006-02", Constraints{Maximum: "2006-01"}},
			{"InvalidMaximum", "2005-02", Constraints{Maximum: "boo"}},
			{"YearSmallerThanMinimum", "2006-02", Constraints{Minimum: "2006-03"}},
			{"InvalidMinimum", "2005-02", Constraints{Minimum: "boo"}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := decodeYearMonth(d.year, d.constraints)
				is.True(err != nil)
			})
		}
	})
}
