package schema

import "testing"

func TestDecodeYear(t *testing.T) {
	t.Run("ValidMaximum", func(t *testing.T) {
		if _, err := decodeYear("2006", Constraints{Maximum: "2007"}); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		if _, err := decodeYear("2007", Constraints{Minimum: "2006"}); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
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
				if _, err := decodeYear(d.year, d.constraints); err == nil {
					t.Fatalf("err want:err got:nil")
				}
			})
		}
	})
}

func TestDecodeYearMonth(t *testing.T) {
	t.Run("ValidMaximum", func(t *testing.T) {
		if _, err := decodeYearMonth("2006-02", Constraints{Maximum: "2006-03"}); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		if _, err := decodeYearMonth("2006-03", Constraints{Minimum: "2006-02"}); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
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
				if _, err := decodeYearMonth(d.year, d.constraints); err == nil {
					t.Fatalf("err want:err got:nil")
				}
			})
		}
	})
}
