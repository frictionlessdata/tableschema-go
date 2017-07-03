package schema

import "testing"

func TestCastValue(t *testing.T) {
	data := []struct {
		Value    string
		Field    Field
		Expected interface{}
	}{
		{"42", Field{Type: IntegerType}, int64(42)},
		{"http:/frictionlessdata.io", Field{Type: StringType, Format: "uri"}, "http:/frictionlessdata.io"},
	}
	for _, d := range data {
		c, err := d.Field.CastValue(d.Value)
		if err != nil {
			t.Errorf("err want:nil got:%s", err)
		}
		if c != d.Expected {
			t.Errorf("val want:%v, got:%v", d.Expected, c)
		}
	}
}

func TestCastValue_InvalidFieldType(t *testing.T) {
	f := Field{Type: "invalidType"}
	if _, err := f.CastValue("42"); err == nil {
		t.Errorf("err want:err, got:nil")
	}
}

func TestTestValue(t *testing.T) {
	f := Field{Type: "integer"}
	if !f.TestValue("42") {
		t.Errorf("want:true, got:false")
	}
	if f.TestValue("boo") {
		t.Errorf("want:false, got:true")
	}
}
