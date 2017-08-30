package schema

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestRead_Sucess(t *testing.T) {
	data := []struct {
		Desc   string
		JSON   string
		Schema Schema
	}{
		{
			"OneField",
			`{
                "fields":[{"name":"n","title":"ti","type":"integer","description":"desc","format":"f","trueValues":["ntrue"],"falseValues":["nfalse"]}]
            }`,
			Schema{
				Fields: []Field{{Name: "n", Title: "ti", Type: "integer", Description: "desc", Format: "f", TrueValues: []string{"ntrue"}, FalseValues: []string{"nfalse"}}},
			},
		},
		{
			"MultipleFields",
			`{
                "fields":[{"name":"n1","type":"t1","format":"f1","falseValues":[]}, {"name":"n2","type":"t2","format":"f2","trueValues":[]}]
            }`,
			Schema{
				Fields: []Field{
					{Name: "n1", Type: "t1", Format: "f1", TrueValues: defaultTrueValues, FalseValues: []string{}},
					{Name: "n2", Type: "t2", Format: "f2", TrueValues: []string{}, FalseValues: defaultFalseValues},
				},
			},
		},
		{
			"PKString",
			`{"fields":[{"name":"n1"}], "primaryKey":"n1"}`,
			Schema{Fields: []Field{asJSONField(Field{Name: "n1"})}, PrimaryKeys: []string{"n1"}},
		},
		{
			"PKSlice",
			`{"fields":[{"name":"n1"}], "primaryKey":["n1"]}`,
			Schema{Fields: []Field{asJSONField(Field{Name: "n1"})}, PrimaryKeys: []string{"n1"}},
		},
		{
			"FKFieldsString",
			`{"fields":[{"name":"n1"}], "foreignKeys":{"fields":"n1"}}`,
			Schema{Fields: []Field{asJSONField(Field{Name: "n1"})}, ForeignKeys: ForeignKeys{Fields: []string{"n1"}}},
		},
		{
			"FKFieldsSlice",
			`{"fields":[{"name":"n1"}], "foreignKeys":{"fields":["n1"]}}`,
			Schema{Fields: []Field{asJSONField(Field{Name: "n1"})}, ForeignKeys: ForeignKeys{Fields: []string{"n1"}}},
		},
		{
			"FKReferenceFieldsString",
			`{"fields":[{"name":"n1"}], "foreignKeys":{"reference":{"fields":"n1"}}}`,
			Schema{Fields: []Field{asJSONField(Field{Name: "n1"})}, ForeignKeys: ForeignKeys{Reference: ForeignKeyReference{Fields: []string{"n1"}}}},
		},
		{
			"FKReferenceFieldsSlice",
			`{"fields":[{"name":"n1"}], "foreignKeys":{"reference":{"fields":["n1"]}}}`,
			Schema{Fields: []Field{asJSONField(Field{Name: "n1"})}, ForeignKeys: ForeignKeys{Reference: ForeignKeyReference{Fields: []string{"n1"}}}},
		},
	}
	for _, d := range data {
		t.Run(d.Desc, func(t *testing.T) {
			s, err := Read(strings.NewReader(d.JSON))
			if err != nil {
				t.Fatalf("want:nil, got:%q", err)
			}
			if !reflect.DeepEqual(s, &d.Schema) {
				t.Errorf("want:%+v, got:%+v", d.Schema, s)
			}
		})
	}
}

func TestRead_Error(t *testing.T) {
	data := []struct {
		Desc string
		JSON string
	}{
		{"InvalidSchema", `{"fields":"f1"}`},
		{"EmptyDescriptor", ""},
		{"InvalidPKType", `{"fields":[{"name":"n1"}], "primaryKey":1}`},
		{"InvalidFKFieldsType", `{"fields":[{"name":"n1"}], "foreignKeys":{"fields":1}}`},
		{"InvalidFKReferenceFieldsType", `{"fields":[{"name":"n1"}], "foreignKeys":{"reference":{"fields":1}}}`},
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

func TestSchema_CastRow(t *testing.T) {
	t.Run("NoImplicitCast", func(t *testing.T) {
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
	})
	t.Run("ImplicitCastToInt", func(t *testing.T) {
		t1 := struct{ Age int }{}
		s := Schema{Fields: []Field{{Name: "name", Type: StringType}, {Name: "age", Type: IntegerType}}}
		if err := s.CastRow([]string{"Foo", "42"}, &t1); err != nil {
			t.Fatalf("err want:nil, got:%q", err)
		}
		if t1.Age != 42 {
			t.Errorf("value:Name want:42, got:%d", t1.Age)
		}
	})
	t.Run("Error_SchemaFieldAndStructFieldDifferentTypes", func(t *testing.T) {
		// Field is string and struct is int.
		t1 := struct{ Age int }{}
		s := Schema{Fields: []Field{{Name: "age", Type: StringType}}}
		if err := s.CastRow([]string{"42"}, &t1); err == nil {
			t.Fatalf("want:error, got:nil")
		}
	})
	t.Run("Error_NotAPointerToStruct", func(t *testing.T) {
		t1 := struct{ Age int }{}
		s := Schema{Fields: []Field{{Name: "name", Type: StringType}}}
		if err := s.CastRow([]string{"Foo", "42"}, t1); err == nil {
			t.Fatalf("want:error, got:nil")
		}
	})
	t.Run("Error_CellCanNotBeCast", func(t *testing.T) {
		// Field is string and struct is int.
		t1 := struct{ Age int }{}
		s := Schema{Fields: []Field{{Name: "age", Type: IntegerType}}}
		if err := s.CastRow([]string{"foo"}, &t1); err == nil {
			t.Fatalf("want:error, got:nil")
		}
	})
	t.Run("Error_CastToNil", func(t *testing.T) {
		t1 := &struct{ Age int }{}
		t1 = nil
		s := Schema{Fields: []Field{{Name: "age", Type: IntegerType}}}
		if err := s.CastRow([]string{"foo"}, &t1); err == nil {
			t.Fatalf("want:error, got:nil")
		}
	})
}

func TestValidate_SimpleValid(t *testing.T) {
	data := []struct {
		Desc   string
		Schema Schema
	}{
		{"PrimaryKey", Schema{Fields: []Field{{Name: "p"}, {Name: "i"}},
			PrimaryKeys: []string{"p"},
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
		{"PKNonexistingField", Schema{Fields: []Field{{Name: "n1"}}, PrimaryKeys: []string{"n2"}}},
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

func TestWrite(t *testing.T) {
	s := Schema{
		Fields:      []Field{{Name: "Foo"}, {Name: "Bar"}},
		PrimaryKeys: []string{"Foo"},
		ForeignKeys: ForeignKeys{Reference: ForeignKeyReference{Fields: []string{"Foo"}}},
	}
	buf := bytes.NewBufferString("")
	if err := s.Write(buf); err != nil {
		t.Errorf("want:nil got:err")
	}
	want := `{
    "fields": [
    {
        "name": "Foo"
    },
    {
        "name": "Bar"
    }
    ],
    "primaryKey": [
        "Foo"
    ],
    "foreignKeys": {
        "reference": {
            "fields": [
                "Foo"
            ]
        }
    }
}`
	if reflect.DeepEqual(buf.String(), want) {
		t.Errorf("val want:%s got:%s", want, buf.String())
	}
}

func TestGetField(t *testing.T) {
	t.Run("HasField", func(t *testing.T) {
		s := Schema{Fields: []Field{{Name: "Foo"}, {Name: "Bar"}}}
		field, pos := s.GetField("Foo")
		if pos != 0 {
			t.Fatalf("pos want:0 got:%d", pos)
		}
		if field == nil {
			t.Fatalf("field want:field got nil")
		}
		field, pos = s.GetField("Bar")
		if pos != 1 {
			t.Fatalf("pos want:1 got:%d", pos)
		}
		if field == nil {
			t.Fatalf("field value want:field got nil")
		}
	})
	t.Run("DoesNotHaveField", func(t *testing.T) {
		s := Schema{Fields: []Field{{Name: "Bez"}}}
		field, pos := s.GetField("Foo")
		if pos != InvalidPosition {
			t.Fatalf("pos want:InvalidPosition got:%d", pos)
		}
		if field != nil {
			t.Fatalf("field value want:nil got:%v", field)
		}
	})
}

func TestHasField(t *testing.T) {
	t.Run("HasField", func(t *testing.T) {
		s := Schema{Fields: []Field{{Name: "Foo"}, {Name: "Bar"}}}
		if !s.HasField("Foo") {
			t.Fatalf("field exist want:true got:false")
		}
		if !s.HasField("Bar") {
			t.Fatalf("field existence want:true got:false")
		}
	})
	t.Run("DoesNotHaveField", func(t *testing.T) {
		s := Schema{Fields: []Field{{Name: "Bez"}}}
		if s.HasField("Bar") {
			t.Fatalf("field existence want:false got:true")
		}
	})
}

func TestMissingValues(t *testing.T) {
	s := Schema{
		Fields:        []Field{{Name: "Foo"}},
		MissingValues: []string{"f"},
	}
	row := struct {
		Foo string
	}{}
	s.CastRow([]string{"f"}, &row)
	if row.Foo != "" {
		t.Fatalf("want:\"\" got:%s", row.Foo)
	}
}
