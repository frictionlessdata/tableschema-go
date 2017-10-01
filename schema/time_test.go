package schema

import (
	"testing"
	"time"
)

func TestDecodeTime(t *testing.T) {
	t.Run("ValidMaximum", func(t *testing.T) {
		if _, err := decodeTime(defaultFieldFormat, "11:45:00", Constraints{Maximum: "11:45:01"}); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
	})
	t.Run("ValidMinimum", func(t *testing.T) {
		if _, err := decodeTime(defaultFieldFormat, "11:45:00", Constraints{Minimum: "11:44:59"}); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
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
				if _, err := decodeTime(defaultFieldFormat, d.time, d.constraints); err == nil {
					t.Fatalf("err want:err got:nil")
				}
			})
		}
	})
}

func TestEncodeTime(t *testing.T) {
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
				got, err := encodeTime(d.value)
				if err != nil {
					t.Fatalf("err want:nil got:%q", err)
				}
				if d.want != got {
					t.Fatalf("val want:%s got:%s", d.want, got)
				}
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
				_, err := encodeTime(d.value)
				if err == nil {
					t.Fatalf("err want:err got:nil")
				}
			})
		}
	})
}
