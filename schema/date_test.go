package schema

import (
	"testing"

	"github.com/matryer/is"
)

func TestCastDate(t *testing.T) {
	t.Run("ValidMaximum", func(t *testing.T) {
		is := is.New(t)
		_, err := castDate("2006-01-02", "2006-01-02", Constraints{Maximum: "2007-01-02"})
		is.NoErr(err)
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		is := is.New(t)
		_, err := castDate("2006-01-02", "2007-01-02", Constraints{Minimum: "2006-01-02"})
		is.NoErr(err)
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc        string
			date        string
			constraints Constraints
		}{
			{"InvalidDate", "foo", Constraints{}},
			{"DateBiggerThanMaximum", "2006-01-02", Constraints{Maximum: "2005-01-02"}},
			{"InvalidMaximum", "2006-01-02", Constraints{Maximum: "boo"}},
			{"DateSmallerThanMinimum", "2005-01-02", Constraints{Minimum: "2006-01-02"}},
			{"InvalidMinimum", "2006-01-02", Constraints{Minimum: "boo"}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := castDate("2006-01-02", d.date, d.constraints)
				is.True(err != nil)
			})
		}
	})
}
