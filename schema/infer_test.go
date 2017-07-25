package schema

import (
	"reflect"
	"sort"
	"testing"
)

func TestInfer_Success(t *testing.T) {
	data := []struct {
		desc    string
		headers []string
		table   [][]string
		want    Schema
	}{
		{"1Cell_Date", []string{"Birthday"}, [][]string{[]string{"1983-10-15"}}, Schema{Fields: []Field{{Name: "Birthday", Type: DateType, Format: defaultFieldFormat}}}},
		{"1Cell_Integer", []string{"Age"}, [][]string{[]string{"10"}}, Schema{Fields: []Field{{Name: "Age", Type: IntegerType, Format: defaultFieldFormat}}}},
		{"1Cell_Number", []string{"Weight"}, [][]string{[]string{"20.2"}}, Schema{Fields: []Field{{Name: "Weight", Type: NumberType, Format: defaultFieldFormat}}}},
		{"ManyCells",
			[]string{"Name", "Age", "Weight"},
			[][]string{
				[]string{"Foo", "10", "20.2"},
				[]string{"Foo", "10", "30"},
				[]string{"Foo", "10", "40.4"},
			},
			Schema{Fields: []Field{
				{Name: "Name", Type: "string", Format: defaultFieldFormat},
				{Name: "Age", Type: "integer", Format: defaultFieldFormat},
				{Name: "Weight", Type: "number", Format: defaultFieldFormat},
			}},
		},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			s, err := Infer(d.headers, d.table)
			if err != nil {
				t.Fatalf("want:nil, got:%q", err)
			}
			sort.Sort(s.Fields)
			sort.Sort(d.want.Fields)
			if !reflect.DeepEqual(s, &d.want) {
				t.Errorf("want:%+v, got:%+v", d.want, s)
			}
		})
	}
}

func TestInfer_Error(t *testing.T) {
	data := []struct {
		desc    string
		headers []string
		table   [][]string
	}{
		{"NotATable", []string{}, [][]string{[]string{"1"}}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			_, err := Infer(d.headers, d.table)
			if err == nil {
				t.Fatalf("want:error, got:nil")
			}
		})
	}
}
