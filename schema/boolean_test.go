package schema

import "testing"

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
		b, err := castBoolean(d.Value, d.TrueValues, d.FalseValues)
		if err != nil {
			t.Errorf("[%s] err want:nil got:%q", d.Desc, err)
		}
		if b != d.Expected {
			t.Errorf("[%s] val want:%v got:%v", d.Desc, d.Expected, b)
		}
	}
}

func TestCastBoolean_Error(t *testing.T) {
	_, err := castBoolean("foo", defaultTrueValues, defaultFalseValues)
	if err == nil {
		t.Errorf("want:err got:nil")
	}
}

func TestEncodeBoolean(t *testing.T) {
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
				got, err := encodeBoolean(d.value, d.trueValues, d.falseValues)
				if err != nil {
					t.Fatalf("err want:nil got:%q", err)
				}
				if d.want != got {
					t.Errorf("val want:%s got:%s", d.want, got)
				}
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
				_, err := encodeBoolean(d.value, d.trueValues, d.falseValues)
				if err == nil {
					t.Fatalf("err want:err got:nil")
				}
			})
		}
	})
}
