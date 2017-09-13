package schema

import (
	"math"
	"testing"
)

const notBareNumber = false

func TestCastNumber(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []struct {
			desc   string
			number string
			want   float64
			dc     string
			gc     string
			bn     bool
		}{
			{"Positive_WithPlus", "+10.10", 10.10, defaultDecimalChar, defaultGroupChar, defaultBareNumber},
			{"Positive_WithoutPlus", "10.10", 10.10, defaultDecimalChar, defaultGroupChar, defaultBareNumber},
			{"Negative_WithPlus", "-10.10", -10.10, defaultDecimalChar, defaultGroupChar, defaultBareNumber},
			{"BareNumber", "€95", 95, defaultDecimalChar, defaultGroupChar, notBareNumber},
			{"BareNumber_TrailingAtBeginning", "€95", 95, defaultDecimalChar, defaultGroupChar, notBareNumber},
			{"BareNumber_TrailingAtBeginningSpace", "EUR 95", 95, defaultDecimalChar, defaultGroupChar, notBareNumber},
			{"BareNumber_TrailingAtEnd", "95%", 95, defaultDecimalChar, defaultGroupChar, notBareNumber},
			{"BareNumber_TrailingAtEndSpace", "95 %", 95, defaultDecimalChar, defaultGroupChar, notBareNumber},
			{"GroupChar", "100,000", 100000, defaultDecimalChar, defaultGroupChar, defaultBareNumber},
			{"DecimalChar", "95;10", 95.10, ";", defaultGroupChar, defaultBareNumber},
			{"Mix", "EUR 95;10", 95.10, ";", ";", notBareNumber},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				got, err := castNumber(d.dc, d.gc, d.bn, d.number)
				if err != nil {
					t.Fatalf("err want:nil got:%q", err)
				}
				if d.want != got {
					t.Fatalf("val want:%f got:%f", d.want, got)
				}
			})
		}
	})
	t.Run("NaN", func(t *testing.T) {
		got, err := castNumber(defaultDecimalChar, defaultGroupChar, defaultBareNumber, "NaN")
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if !math.IsNaN(got) {
			t.Fatalf("val want:NaN got:%f", got)
		}
	})
	t.Run("INF", func(t *testing.T) {
		got, err := castNumber(defaultDecimalChar, defaultGroupChar, defaultBareNumber, "INF")
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if !math.IsInf(got, 1) {
			t.Fatalf("val want:+Inf got:%f", got)
		}
	})
	t.Run("NegativeINF", func(t *testing.T) {
		got, err := castNumber(defaultDecimalChar, defaultGroupChar, defaultBareNumber, "-INF")
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if !math.IsInf(got, -1) {
			t.Fatalf("val want:-Inf got:%f", got)
		}
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc   string
			number string
			dc     string
			gc     string
			bn     bool
		}{
			{"InvalidNumberToStrip_TooManyNumbers", "+10.10++10", defaultDecimalChar, defaultGroupChar, notBareNumber},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				if _, err := castNumber(defaultDecimalChar, defaultGroupChar, defaultBareNumber, d.number); err == nil {
					t.Fatalf("err want:err got:nil")
				}
			})
		}
	})
}
