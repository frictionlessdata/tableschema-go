package schema

import (
	"testing"

	"github.com/matryer/is"
)

func TestCastBoolean(t *testing.T) {
	data := []struct {
		Desc        string
		TrueValues  []string
		FalseValues []string
		Value       string
		Expected    bool
	}{
		{"simple true value", []string{"1"}, []string{"0"}, "1", true},
		{"simple false value", []string{"1"}, []string{"0"}, "0", false},
		{"duplicate value, true wins", []string{"1"}, []string{"1"}, "1", true},
	}
	for _, d := range data {
		t.Run(d.Desc, func(t *testing.T) {
			is := is.New(t)
			b, err := castBoolean(d.Value, d.TrueValues, d.FalseValues)
			is.NoErr(err)
			is.Equal(b, d.Expected)
		})
	}
}

func TestCastBoolean_Error(t *testing.T) {
	is := is.New(t)
	_, err := castBoolean("foo", defaultTrueValues, defaultFalseValues)
	is.True(err != nil)
}

func TestUncastBoolean(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []struct {
			desc        string
			value       interface{}
			want        string
			trueValues  []string
			falseValues []string
		}{
			{"True", true, "true", []string{}, []string{}},
			{"False", false, "false", []string{}, []string{}},
			{"TrueFromTrueValues", "0", "0", []string{"0"}, []string{}},
			{"FalseFromFalseValues", "1", "1", []string{}, []string{"1"}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				got, err := uncastBoolean(d.value, d.trueValues, d.falseValues)
				is.NoErr(err)
				is.Equal(d.want, got)
			})
		}
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc        string
			value       interface{}
			trueValues  []string
			falseValues []string
		}{
			{"InvalidType", 10, []string{}, []string{}},
			{"NotInTrueOrFalseValues", "1", []string{}, []string{}},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				is := is.New(t)
				_, err := uncastBoolean(d.value, d.trueValues, d.falseValues)
				is.True(err != nil)
			})
		}
	})
}
