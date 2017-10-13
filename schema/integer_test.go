package schema

import (
	"testing"

	"github.com/matryer/is"
)

const notBareInt = false

func TestCastInt(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []struct {
			desc   string
			number string
			want   int64
			bn     bool
		}{
			{"Positive_WithSignal", "+10", 10, defaultBareNumber},
			{"Positive_WithoutSignal", "10", 10, defaultBareNumber},
			{"Negative", "-10", -10, defaultBareNumber},
			{"BareNumber", "€95", 95, notBareInt},
			{"BareNumber_TrailingAtBeginning", "€95", 95, notBareInt},
			{"BareNumber_TrailingAtBeginningSpace", "EUR 95", 95, notBareInt},
			{"BareNumber_TrailingAtEnd", "95%", 95, notBareInt},
			{"BareNumber_TrailingAtEndSpace", "95 %", 95, notBareInt},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				got, err := castInt(d.bn, d.number, Constraints{})
				is.NoErr(err)
				is.Equal(d.want, got)
			})
		}
	})
	t.Run("ValidMaximum", func(t *testing.T) {
		is := is.New(t)
		_, err := castInt(defaultBareNumber, "2", Constraints{Maximum: "2"})
		is.NoErr(err)
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		is := is.New(t)
		_, err := castInt(defaultBareNumber, "2", Constraints{Minimum: "1"})
		is.NoErr(err)
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc        string
			number      string
			constraints Constraints
		}{
			{"InvalidIntToStrip_TooManyNumbers", "+10++10", Constraints{}},
			{"NumBiggerThanMaximum", "3", Constraints{Maximum: "2"}},
			{"InvalidMaximum", "1", Constraints{Maximum: "boo"}},
			{"NumSmallerThanMinimum", "1", Constraints{Minimum: "2"}},
			{"InvalidMinimum", "1", Constraints{Minimum: "boo"}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := castInt(defaultBareNumber, d.number, d.constraints)
				is.True(err != nil)
			})
		}
	})
}
