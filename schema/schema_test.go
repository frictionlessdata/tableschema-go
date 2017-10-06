package schema

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"reflect"
	"strings"
	"testing"

	"github.com/frictionlessdata/tableschema-go/table"
)

func ExampleSchema_Decode() {
	// Lets assume we have a schema ...
	s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}

	// And a Table.
	t := table.FromSlices([]string{"Name", "Age"}, [][]string{
		{"Foo", "42"},
		{"Bar", "43"}})

	// And we would like to process them using Go types. First we need to create a struct to
	// hold the content of each row.
	// The tag tableheader maps the field to the schema. If no tag is set the name of the field
	// has to be the same like inside the schema.
	type person struct {
		MyName string `tableheader:"Name"`
		Age    int
	}

	// Now it is a matter of iterate over the table and Decode each row.
	iter, _ := t.Iter()
	for iter.Next() {
		var p person
		s.Decode(iter.Row(), &p)
		fmt.Printf("%+v\n", p)
	}
	// Output: {MyName:Foo Age:42}
	// {MyName:Bar Age:43}
}

func ExampleSchema_DecodeTable() {
	// Lets assume we have a schema ...
	s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}

	// And a Table.
	t := table.FromSlices([]string{"Name", "Age"}, [][]string{
		{"Foo", "42"},
		{"Bar", "43"}})

	// And we would like to process them using Go types. First we need to create a struct to
	// hold the content of each row.
	// The tag tableheader maps the field to the schema. If no tag is set the name of the field
	// has to be the same like inside the schema.
	type person struct {
		MyName string `tableheader:"Name"`
		Age    int
	}
	var people []person
	s.DecodeTable(t, &people)
	fmt.Print(people)
	// Output: [{Foo 42} {Bar 43}]
}

func ExampleSchema_Encode() {
	// Lets assume we have a schema.
	s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}

	// And would like to create a CSV out of this list. The tag tableheader maps
	// the field to the schema name. If no tag is set the name of the field
	// has to be the same like inside the schema.
	people := []struct {
		MyName string `tableheader:"Name"`
		Age    int
	}{{"Foo", 42}, {"Bar", 43}}

	// First create the writer and write the header.
	w := table.NewStringWriter()
	w.Write([]string{"Name", "Age"})

	// Then write the list
	for _, person := range people {
		row, _ := s.Encode(person)
		w.Write(row)
	}
	w.Flush()
	fmt.Print(w.String())
	// Output: Name,Age
	// Foo,42
	// Bar,43
}

func ExampleSchema_EncodeTable() {
	// Lets assume we have a schema.
	s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}

	// And would like to create a CSV out of this list. The tag tableheader maps
	// the field to the schema name. If no tag is set the name of the field
	// has to be the same like inside the schema.
	people := []struct {
		MyName string `tableheader:"Name"`
		Age    int
	}{{"Foo", 42}, {"Bar", 43}}

	// Then encode the people slice into a slice of rows.
	rows, _ := s.EncodeTable(people)

	// Now, simply write it down.
	w := table.NewStringWriter()
	w.Write([]string{"Name", "Age"})
	w.WriteAll(rows)
	w.Flush()
	fmt.Print(w.String())
	// Output: Name,Age
	// Foo,42
	// Bar,43
}

func TestLoadRemote(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"fields": [{"name": "ID", "type": "integer"}]}`)
	}
	ts := httptest.NewServer(http.HandlerFunc(h))
	defer ts.Close()
	got, err := LoadRemote(ts.URL)
	if err != nil {
		t.Fatalf("want:nil, got:%q", err)
	}
	want := &Schema{Fields: []Field{asJSONField(Field{Name: "ID", Type: "integer"})}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want:%+v, got:%+v", want, got)
	}
	t.Run("Error", func(t *testing.T) {
		_, err := LoadRemote("invalidURL")
		if err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

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
				Fields: []Field{{Name: "n", Title: "ti", Type: "integer", Description: "desc", Format: "f", TrueValues: []string{"ntrue"}, FalseValues: []string{"nfalse"},
					DecimalChar: defaultDecimalChar, GroupChar: defaultGroupChar, BareNumber: defaultBareNumber}},
			},
		},
		{
			"MultipleFields",
			`{
                "fields":[{"name":"n1","type":"t1","format":"f1","falseValues":[]}, {"name":"n2","type":"t2","format":"f2","trueValues":[]}]
            }`,
			Schema{
				Fields: []Field{
					{Name: "n1", Type: "t1", Format: "f1", TrueValues: defaultTrueValues, FalseValues: []string{}, DecimalChar: defaultDecimalChar, GroupChar: defaultGroupChar, BareNumber: defaultBareNumber},
					{Name: "n2", Type: "t2", Format: "f2", TrueValues: []string{}, FalseValues: defaultFalseValues, DecimalChar: defaultDecimalChar, GroupChar: defaultGroupChar, BareNumber: defaultBareNumber},
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
				t.Errorf("want:%+v, got:%+v", &d.Schema, s)
			}
		})
	}
	t.Run("MissingValues", func(t *testing.T) {
		reader := strings.NewReader(`{"fields":[{"name":"n","type":"integer"}],"missingValues":["na"]}`)
		s, err := Read(reader)
		if err != nil {
			t.Fatalf("want:nil, got:%q", err)
		}
		f := s.Fields[0]
		if _, ok := f.MissingValues["na"]; !ok {
			t.Fatalf("want:ok got:!ok")
		}
	})
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

func TestSchema_Decode(t *testing.T) {
	t.Run("NoImplicitCast", func(t *testing.T) {
		t1 := struct {
			Name string
			Age  int64
		}{}
		s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}
		if err := s.Decode([]string{"Foo", "42"}, &t1); err != nil {
			t.Fatalf("err want:nil, got:%q", err)
		}
		if t1.Name != "Foo" {
			t.Errorf("value:Name want:Foo got:%s", t1.Name)
		}
		if t1.Age != 42 {
			t.Errorf("value:Age want:42 got:%d", t1.Age)
		}
	})
	t.Run("StructWithTags", func(t *testing.T) {
		t1 := struct {
			MyName string `tableheader:"Name"`
			MyAge  int64  `tableheader:"Age"`
		}{}
		s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}
		if err := s.Decode([]string{"Foo", "42"}, &t1); err != nil {
			t.Fatalf("err want:nil, got:%q", err)
		}
		if t1.MyName != "Foo" {
			t.Errorf("value:Name want:Foo got:%s", t1.MyName)
		}
		if t1.MyAge != 42 {
			t.Errorf("value:Age want:42 got:%d", t1.MyAge)
		}
	})
	t.Run("ImplicitCastToInt", func(t *testing.T) {
		t1 := struct{ Age int }{}
		s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}
		if err := s.Decode([]string{"Foo", "42"}, &t1); err != nil {
			t.Fatalf("err want:nil, got:%q", err)
		}
		if t1.Age != 42 {
			t.Errorf("value:Name want:42, got:%d", t1.Age)
		}
	})
	t.Run("Error_SchemaFieldAndStructFieldDifferentTypes", func(t *testing.T) {
		// Field is string and struct is int.
		t1 := struct{ Age int }{}
		s := Schema{Fields: []Field{{Name: "Age", Type: StringType}}}
		if err := s.Decode([]string{"42"}, &t1); err == nil {
			t.Fatalf("want:error, got:nil")
		}
	})
	t.Run("Error_NotAPointerToStruct", func(t *testing.T) {
		t1 := struct{ Age int }{}
		s := Schema{Fields: []Field{{Name: "Name", Type: StringType}}}
		if err := s.Decode([]string{"Foo", "42"}, t1); err == nil {
			t.Fatalf("want:error, got:nil")
		}
	})
	t.Run("Error_CellCanNotBeCast", func(t *testing.T) {
		// Field is string and struct is int.
		t1 := struct{ Age int }{}
		s := Schema{Fields: []Field{{Name: "Age", Type: IntegerType}}}
		if err := s.Decode([]string{"foo"}, &t1); err == nil {
			t.Fatalf("want:error, got:nil")
		}
	})
	t.Run("Error_CastToNil", func(t *testing.T) {
		t1 := &struct{ Age int }{}
		t1 = nil
		s := Schema{Fields: []Field{{Name: "Age", Type: IntegerType}}}
		if err := s.Decode([]string{"foo"}, &t1); err == nil {
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
	s.Decode([]string{"f"}, &row)
	if row.Foo != "" {
		t.Fatalf("want:\"\" got:%s", row.Foo)
	}
}

type csvRow struct {
	Name string
}

func TestDecodeTable(t *testing.T) {
	data := []struct {
		desc string
		got  []csvRow
	}{
		{"OutEmpty", []csvRow{}},
		{"OutNil", nil},
		{"OutInitialized", []csvRow{{"fooooo"}}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			tab := table.FromSlices(
				[]string{"Name"},
				[][]string{{"foo"}, {"bar"}})
			s := &Schema{Fields: []Field{{Name: "Name", Type: StringType}}}
			if err := s.DecodeTable(tab, &d.got); err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			want := []csvRow{{"foo"}, {"bar"}}
			if !reflect.DeepEqual(want, d.got) {
				t.Fatalf("val want:%v got:%v", want, d.got)
			}
		})
	}
	t.Run("MoarData", func(t *testing.T) {
		tab := table.FromSlices(
			[]string{"ID", "Age", "Name"},
			[][]string{{"1", "39", "Paul"}, {"2", "23", "Jimmy"}, {"3", "36", "Jane"}, {"4", "28", "Judy"}, {"5", "37", "Iñtërnâtiônàlizætiøn"}})

		type data struct {
			ID   int
			Age  int
			Name string
		}
		s := &Schema{Fields: []Field{{Name: "ID", Type: IntegerType}, {Name: "Age", Type: IntegerType}, {Name: "Name", Type: StringType}}}
		got := []data{}
		if err := s.DecodeTable(tab, &got); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []data{{1, 39, "Paul"}, {2, 23, "Jimmy"}, {3, 36, "Jane"}, {4, 28, "Judy"}, {5, 37, "Iñtërnâtiônàlizætiøn"}}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%v got:%v", want, got)
		}
	})
	t.Run("EmptyTable", func(t *testing.T) {
		tab := table.FromSlices([]string{}, [][]string{})
		s := &Schema{Fields: []Field{{Name: "name", Type: StringType}}}
		var got []csvRow
		if err := s.DecodeTable(tab, &got); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(got) != 0 {
			t.Fatalf("len(got) want:0 got:%v", len(got))
		}
	})
	t.Run("Error_OutNotAPointerToSlice", func(t *testing.T) {
		tab := table.FromSlices([]string{"name"}, [][]string{{""}})
		s := &Schema{Fields: []Field{{Name: "name", Type: StringType}}}
		if err := s.DecodeTable(tab, []csvRow{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
}

func TestSchema_Encode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		type rowType struct {
			Name string
			Age  int
		}
		s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}
		got, err := s.Encode(rowType{Name: "Foo", Age: 42})
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"Foo", "42"}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%v got:%v", want, got)
		}
	})
	t.Run("SuccessWithTags", func(t *testing.T) {
		type rowType struct {
			MyName string `tableheader:"Name"`
			MyAge  int    `tableheader:"Age"`
		}
		s := Schema{Fields: []Field{{Name: "Name", Type: StringType}, {Name: "Age", Type: IntegerType}}}
		got, err := s.Encode(rowType{MyName: "Foo", MyAge: 42})
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"Foo", "42"}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%v got:%v", want, got)
		}
	})
	t.Run("Error_Encoding", func(t *testing.T) {
		type rowType struct {
			Age string
		}
		s := Schema{Fields: []Field{{Name: "Age", Type: IntegerType}}}
		_, err := s.Encode(rowType{Age: "10"})
		if err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("Error_NotStruct", func(t *testing.T) {
		s := Schema{Fields: []Field{{Name: "name", Type: StringType}}}
		in := "string"
		_, err := s.Encode(in)
		if err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("Error_StructIsNil", func(t *testing.T) {
		s := Schema{Fields: []Field{{Name: "name", Type: StringType}}}
		var in *csvRow
		_, err := s.Encode(in)
		if err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
}

func TestEncodeTable(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		people := []struct {
			Name string
		}{{"Foo"}, {"Bar"}, {"Bez"}}
		s := Schema{Fields: []Field{{Name: "Name", Type: StringType}}}
		got, err := s.EncodeTable(people)
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := [][]string{{"Foo"}, {"Bar"}, {"Bez"}}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%s got:%q", want, got)
		}
	})

	t.Run("Error_InputIsNotASlice", func(t *testing.T) {
		s := Schema{Fields: []Field{{Name: "Name", Type: StringType}}}
		_, err := s.EncodeTable(10)
		if err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("Error_ErrorEncoding", func(t *testing.T) {
		people := []struct {
			Name string
		}{{"Foo"}}
		s := Schema{Fields: []Field{{Name: "Name", Type: IntegerType}}}
		_, err := s.EncodeTable(people)
		if err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
}
