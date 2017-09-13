package schema

import "testing"

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
				got, err := castInt(d.bn, d.number)
				if err != nil {
					t.Fatalf("err want:nil got:%q", err)
				}
				if d.want != got {
					t.Fatalf("val want:%d got:%d", d.want, got)
				}
			})
		}
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc   string
			number string
			bn     bool
		}{
			{"InvalidIntToStrip_TooManyNumbers", "+10++10", notBareInt},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				if _, err := castInt(defaultBareNumber, d.number); err == nil {
					t.Fatalf("err want:err got:nil")
				}
			})
		}
	})
}
