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
