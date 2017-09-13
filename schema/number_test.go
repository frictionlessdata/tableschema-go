package schema

import (
	"math"
	"testing"
)

func TestCastNumber(t *testing.T) {
	t.Run("Common_Cases", func(t *testing.T) {
		data := []struct {
			desc   string
			number string
			want   float64
		}{
			{"Positive_WithPlus", "+10.10", 10.10},
			{"Positive_WithoutPlus", "10.10", 10.10},
			{"Negative_WithPlus", "-10.10", -10.10},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				got, err := castNumber(d.number)
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
		got, err := castNumber("NaN")
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if !math.IsNaN(got) {
			t.Fatalf("val want:NaN got:%f", got)
		}
	})
	t.Run("INF", func(t *testing.T) {
		got, err := castNumber("INF")
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if !math.IsInf(got, 1) {
			t.Fatalf("val want:+Inf got:%f", got)
		}
	})
	t.Run("NegativeINF", func(t *testing.T) {
		got, err := castNumber("-INF")
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if !math.IsInf(got, -1) {
			t.Fatalf("val want:-Inf got:%f", got)
		}
	})
}
