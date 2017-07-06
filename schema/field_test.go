package schema

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDefaultValues(t *testing.T) {
	data := []struct {
		Desc  string
		JSON  string
		Field Field
	}{
		{
			"Default Values",
			`{"name":"n1"}`,
			Field{Name: "n1", Type: defaultFieldType, Format: defaultFieldFormat, TrueValues: defaultTrueValues, FalseValues: defaultFalseValues},
		},
		{
			"Overrinding default values",
			`{"name":"n2","type":"t2","format":"f2","falseValues":["f2"],"trueValues":["t2"]}`,
			Field{Name: "n2", Type: "t2", Format: "f2", TrueValues: []string{"t2"}, FalseValues: []string{"f2"}},
		},
	}
	for _, d := range data {
		var f Field
		if err := json.Unmarshal([]byte(d.JSON), &f); err != nil {
			t.Errorf("err want:nil got:%q", err)
		}
		if !reflect.DeepEqual(f, d.Field) {
			t.Errorf("[%s] want:%+v got:%+v", d.Desc, d.Field, f)
		}
	}
}

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
