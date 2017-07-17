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
			"OneField",
			`{
                "fields":[{"name":"n","title":"ti","type":"integer","description":"desc","format":"f","trueValues":["ntrue"],"falseValues":["nfalse"]}]
            }`,
			&Schema{
				Fields: []Field{{Name: "n", Title: "ti", Type: "integer", Description: "desc", Format: "f", TrueValues: []string{"ntrue"}, FalseValues: []string{"nfalse"}}},
			},
		},
		{
			"MultipleFields",
			`{
                "fields":[{"name":"n1","type":"t1","format":"f1","falseValues":[]}, {"name":"n2","type":"t2","format":"f2","trueValues":[]}]
            }`,
			&Schema{
				Fields: []Field{
					{Name: "n1", Type: "t1", Format: "f1", TrueValues: defaultTrueValues, FalseValues: []string{}},
					{Name: "n2", Type: "t2", Format: "f2", TrueValues: []string{}, FalseValues: defaultFalseValues},
				},
			},
		},
	}
	for _, d := range data {
		t.Run(d.Desc, func(t *testing.T) {
			s, err := Read(strings.NewReader(d.JSON))
			if err != nil {
				t.Fatalf("want:nil, got:%q", err)
			}
			if !reflect.DeepEqual(s, d.Schema) {
				t.Errorf("want:%+v, got:%+v", d.Schema, s)
			}
		})
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
		t.Run(d.Desc, func(t *testing.T) {
			_, err := Read(strings.NewReader(d.JSON))
			if err == nil {
				t.Fatalf("want:error, got:nil")
			}
		})
	}
}

func TestCastRow_NoImplicitCast(t *testing.T) {
	t1 := struct {
		Name string
		Age  int64
	}{}
	s := Schema{Fields: []Field{{Name: "name", Type: StringType}, {Name: "age", Type: IntegerType}}}
	if err := s.CastRow([]string{"Foo", "42"}, &t1); err != nil {
		t.Fatalf("err want:nil, got:%q", err)
	}
	if t1.Name != "Foo" {
		t.Errorf("value:Name want:Foo got:%s", t1.Name)
	}
	if t1.Age != 42 {
		t.Errorf("value:Age want:42 got:%d", t1.Age)
	}
}

func TestCastRow_ImplicitCastToInt(t *testing.T) {
	t1 := struct{ Age int }{}
	s := Schema{Fields: []Field{{Name: "name", Type: StringType}, {Name: "age", Type: IntegerType}}}
	if err := s.CastRow([]string{"Foo", "42"}, &t1); err != nil {
		t.Fatalf("err want:nil, got:%q", err)
	}
	if t1.Age != 42 {
		t.Errorf("value:Name want:42, got:%d", t1.Age)
	}
}

func TestCastRow_SchemaFieldAndStructFieldDifferentTypes(t *testing.T) {
	// Field is string and struct is int.
	t1 := struct{ Age int }{}
	s := Schema{Fields: []Field{{Name: "age", Type: StringType}}}
	if err := s.CastRow([]string{"42"}, &t1); err == nil {
		t.Fatalf("want:error, got:nil")
	}
}

func TestCastRow_NotAPointerToStruct(t *testing.T) {
	t1 := struct{ Age int }{}
	s := Schema{Fields: []Field{{Name: "name", Type: StringType}}}
	if err := s.CastRow([]string{"Foo", "42"}, t1); err == nil {
		t.Fatalf("want:error, got:nil")
	}
}

func TestCastRow_CellCanNotBeCast(t *testing.T) {
	// Field is string and struct is int.
	t1 := struct{ Age int }{}
	s := Schema{Fields: []Field{{Name: "age", Type: IntegerType}}}
	if err := s.CastRow([]string{"foo"}, &t1); err == nil {
		t.Fatalf("want:error, got:nil")
	}
}

func TestValidate_SimpleValid(t *testing.T) {
	data := []struct {
		Desc   string
		Schema Schema
	}{
		{"PrimaryKey", Schema{Fields: []Field{{Name: "p"}, {Name: "i"}},
			PrimaryKeys: PrimaryKeys{"p"},
			ForeignKeys: ForeignKeys{
				Fields:    []string{"p"},
				Reference: ForeignKeyReference{Resource: "", Fields: []string{"i"}},
			}},
		},
	}
	for _, d := range data {
		t.Run(d.Desc, func(t *testing.T) {
			if err := d.Schema.Validate(); err != nil {
				t.Errorf("want:nil got:%q", err)
			}
		})
	}
}

func TestValidate_Invalid(t *testing.T) {
	data := []struct {
		Desc   string
		Schema Schema
	}{
		{"MissingName", Schema{Fields: []Field{{Type: IntegerType}}}},
		{"PKNonexistingField", Schema{Fields: []Field{{Name: "n1"}}, PrimaryKeys: PrimaryKeys{"n2"}}},
		{"FKNonexistingField", Schema{Fields: []Field{{Name: "n1"}},
			ForeignKeys: ForeignKeys{Fields: []string{"n2"}},
		}},
		{"InvalidReferences", Schema{Fields: []Field{{Name: "n1"}},
			ForeignKeys: ForeignKeys{
				Fields:    []string{"n1"},
				Reference: ForeignKeyReference{Resource: "", Fields: []string{"n1", "n2"}},
			}},
		},
	}
	for _, d := range data {
		t.Run(d.Desc, func(t *testing.T) {
			if err := d.Schema.Validate(); err == nil {
				t.Errorf("want:err got:nil")
			}
		})
	}
}
