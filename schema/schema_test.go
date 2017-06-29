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
                "fields":[{"name":"n","type":"t","format":"f"}]
            }`,
			&Schema{
				[]Field{{Name: "n", Type: "t", Format: "f"}},
			},
		},
		{
			"Multiple fields",
			`{
                "fields":[{"name":"n1","type":"t1","format":"f1"}, {"name":"n2","type":"t2","format":"f2"}]
            }`,
			&Schema{
				[]Field{{Name: "n1", Type: "t1", Format: "f1"}, {Name: "n2", Type: "t2", Format: "f2"}},
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
