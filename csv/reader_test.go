package csv

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/frictionlessdata/tableschema-go/schema"
)

type csvRow struct {
	Name string
}

func ExampleReader_Iter() {
	reader, _ := NewReader(FromString("\"name\"\nfoo\nbar"), LoadHeaders())
	reader.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
	iter, _ := reader.Iter()
	for iter.Next() {
		var data csvRow
		iter.UnmarshalRow(&data)
		fmt.Println(data.Name)
	}
	// Output:foo
	// bar
}

func ExampleInferSchema() {
	reader, _ := NewReader(FromString("\"name\"\nfoo\nbar"), LoadHeaders(), InferSchema())
	iter, _ := reader.Iter()
	for iter.Next() {
		var data csvRow
		iter.UnmarshalRow(&data)
		fmt.Println(data.Name)
	}
	// Output:foo
	// bar
}

func TestLoadHeaders(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		reader, err := NewReader(FromString(""), LoadHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(reader.Headers) != 0 {
			t.Fatalf("len(headers) want:0 got:%v", len(reader.Headers))
		}
	})
	t.Run("SimpleCase", func(t *testing.T) {
		in := `"name"
"bar"`
		reader, err := NewReader(FromString(in), LoadHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"name"}
		if !reflect.DeepEqual(want, reader.Headers) {
			t.Fatalf("headers want:%v got:%v", want, reader.Headers)
		}
		reader.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		var out []csvRow
		if err := reader.UnmarshalAll(&out); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(out) != 1 {
			t.Fatalf("LoadHeaders content must skip first row")
		}
	})
}

func TestNewReader(t *testing.T) {
	t.Run("ErrorOpts", func(t *testing.T) {
		reader, err := NewReader(FromString(""), errorOpts())
		if reader != nil {
			t.Fatalf("reader want:nil got:%v", reader)
		}
		if err == nil {
			t.Fatalf("err want:error got:nil")
		}
	})
	t.Run("ErrorSource", func(t *testing.T) {
		_, err := NewReader(errorSource(), LoadHeaders())
		if err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

func TestSetHeaders(t *testing.T) {
	in := "Foo"
	reader, err := NewReader(FromString(in), SetHeaders("name"))
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	want := []string{"name"}
	if !reflect.DeepEqual(want, reader.Headers) {
		t.Fatalf("val want:%v got:%v", want, reader.Headers)
	}
	reader.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
	var out []csvRow
	if err := reader.UnmarshalAll(&out); err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	if len(out) == 0 {
		t.Fatalf("CSVHeaders must not skip first row")
	}
}

func TestInferSchema(t *testing.T) {
	t.Run("SimpleCase", func(t *testing.T) {
		reader, err := NewReader(FromString("\"name\"\nfoo\nbar"), LoadHeaders(), InferSchema())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		var got []csvRow
		if err := reader.UnmarshalAll(&got); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []csvRow{{"foo"}, {"bar"}}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%v got:%v", want, got)
		}
	})
	t.Run("WithErrorSource", func(t *testing.T) {
		_, err := NewReader(errorSource(), InferSchema())
		if err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

func TestReader_Iter(t *testing.T) {
	t.Run("SimpleCase", func(t *testing.T) {
		reader, err := NewReader(FromString("\"name\"\nfoo\nbar"), LoadHeaders(), InferSchema())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		iter, err := reader.Iter()
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := [][]string{{"foo"}, {"bar"}}
		for i := range want {
			if !iter.Next() {
				t.Fatalf("want more values")
			}
			if !reflect.DeepEqual(want[i], iter.Row()) {
				t.Fatalf("val want:%v got:%v", want[i], iter.Row())
			}
			if iter.Err() != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
		}
	})
	t.Run("WithErrorSource", func(t *testing.T) {
		reader, err := NewReader(errorSource())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		_, err = reader.Iter()
		if err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

func TestReader_UnmarshalAll(t *testing.T) {
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
			reader, err := NewReader(FromString("name\nfoo\nbar"))
			if err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			reader.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
			if err := reader.UnmarshalAll(&d.got); err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			want := []csvRow{{"name"}, {"foo"}, {"bar"}}
			if !reflect.DeepEqual(want, d.got) {
				t.Fatalf("val want:%v got:%v", want, d.got)
			}
		})
	}
	t.Run("MoarData", func(t *testing.T) {
		reader, err := NewReader(FromString(`1,39,Paul
2,23,Jimmy
3,36,Jane
4,28,Judy
5,37,Iñtërnâtiônàlizætiøn`))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		type data struct {
			ID   int
			Age  int
			Name string
		}
		got := []data{}
		reader.Schema = &schema.Schema{Fields: []schema.Field{{Name: "id", Type: schema.IntegerType}, {Name: "age", Type: schema.IntegerType}, {Name: "name", Type: schema.StringType}}}
		if err := reader.UnmarshalAll(&got); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []data{
			{1, 39, "Paul"},
			{2, 23, "Jimmy"},
			{3, 36, "Jane"},
			{4, 28, "Judy"},
			{5, 37, "Iñtërnâtiônàlizætiøn"},
		}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%v got:%v", want, got)
		}
	})
	t.Run("EmptyString", func(t *testing.T) {
		reader, err := NewReader(FromString(""))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		reader.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		var got []csvRow
		if err := reader.UnmarshalAll(&got); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(got) != 0 {
			t.Fatalf("len(got) want:0 got:%v", len(got))
		}
	})
	t.Run("Error_ReaderWithNoSchema", func(t *testing.T) {
		reader, err := NewReader(FromString("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if err := reader.UnmarshalAll(&[]csvRow{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("Error_OutNotAPointerToSlice", func(t *testing.T) {
		reader, err := NewReader(FromString("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		reader.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		if err := reader.UnmarshalAll([]csvRow{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
}

func TestWithSchema(t *testing.T) {
	s := &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
	reader, err := NewReader(FromString("name\nfoo\nbar"), WithSchema(s))
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	if !reflect.DeepEqual(s, reader.Schema) {
		t.Fatalf("schema want:%v got:%v", s, reader.Schema)
	}
}

func TestReader_All(t *testing.T) {
	reader, err := NewReader(FromString("name\nfoo\nbar"))
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	want := [][]string{
		[]string{"name"},
		[]string{"foo"},
		[]string{"bar"},
	}
	got, err := reader.All()
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("schema want:%v got:%v", want, got)
	}
}
