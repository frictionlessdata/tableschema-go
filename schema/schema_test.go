package schema

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadSucess(t *testing.T) {
	data := []struct {
		Desc   string
		JSON   string
		Schema *Schema
	}{
		{
			"One field",
			`{
                "fields":[{"name":"n","type":"t","format":"f","trueValues":["ntrue"],"falseValues":["nfalse"]}]
            }`,
			&Schema{
				[]Field{{Name: "n", Type: "t", Format: "f", TrueValues: []string{"ntrue"}, FalseValues: []string{"nfalse"}}},
			},
		},
		{
			"Multiple fields",
			`{
                "fields":[{"name":"n1","type":"t1","format":"f1","falseValues":[]}, {"name":"n2","type":"t2","format":"f2","trueValues":[]}]
            }`,
			&Schema{
				[]Field{
					{Name: "n1", Type: "t1", Format: "f1", TrueValues: defaultTrueValues, FalseValues: []string{}},
					{Name: "n2", Type: "t2", Format: "f2", TrueValues: []string{}, FalseValues: defaultFalseValues},
				},
			},
		},
	}
	for _, d := range data {
		s, err := Read(strings.NewReader(d.JSON))
		if err != nil {
			t.Fatalf("[%s] want:nil, got:%q", d.Desc, err)
		}
		if !reflect.DeepEqual(s, d.Schema) {
			t.Errorf("[%s] want:%+v, got:%+v", d.Desc, d.Schema, s)
		}
	}
}

func TestHeaders(t *testing.T) {
	// Empty schema, empty headers.
	s := Schema{}
	if len(s.Headers()) > 0 {
		t.Errorf("want:0 got:%d", len(s.Headers()))
	}

	s1 := Schema{Fields: []Field{{Name: "f1"}, {Name: "f2"}}}
	expected := []string{"f1", "f2"}
	if !reflect.DeepEqual(s1.Headers(), expected) {
		t.Errorf("want:%v got:%v", expected, s1.Headers())
	}
}

func TestReadError(t *testing.T) {
	data := []struct {
		Desc string
		JSON string
	}{
		{"empty descriptor", ""},
	}
	for _, d := range data {
		_, err := Read(strings.NewReader(d.JSON))
		if err == nil {
			t.Fatalf("[%s] want:error, got:nil", d.Desc)
		}
	}
}
