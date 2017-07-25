package schema

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
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
		{"1", Field{Type: BooleanType, TrueValues: []string{"1"}}, true},
		{"0", Field{Type: BooleanType, FalseValues: []string{"0"}}, false},
		{"42.5", Field{Type: NumberType}, 42.5},
		{"2015-10-15", Field{Type: DateType}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{"2015-10-15", Field{Type: DateType, Format: defaultFieldFormat}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{"15/10/2015", Field{Type: DateType, Format: "%d/%m/%Y"}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
	}
	for _, d := range data {
		c, err := d.Field.CastValue(d.Value)
		if err != nil {
			t.Fatalf("err want:nil got:%s", err)
		}
		if c != d.Expected {
			t.Errorf("val want:%v, got:%v", d.Expected, c)
		}
	}
}

func TestCastValue_InvalidDate(t *testing.T) {
	data := []struct {
		desc  string
		field Field
		value string
	}{
		{"InvalidFormat_Any", Field{Type: DateType, Format: "any"}, "2015-10-15"},
		{"InvalidFormat_Strftime", Field{Type: DateType, Format: "Fooo"}, "2015-10-15"},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			if _, err := d.field.CastValue(d.value); err == nil {
				t.Errorf("want:err got:nil")
			}
		})
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
