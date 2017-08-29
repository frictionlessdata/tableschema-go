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

func ExampleTable_Iter() {
	tab, _ := New(FromString("\"name\"\nfoo\nbar"), LoadHeaders())
	tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
	iter, _ := tab.Iter()
	for iter.Next() {
		var data csvRow
		iter.CastRow(&data)
		fmt.Println(data.Name)
	}
	// Output:foo
	// bar
}

func ExampleTable_Infer() {
	tab, _ := New(FromString("\"name\"\nfoo\nbar"), LoadHeaders())
	tab.Infer()
	iter, _ := tab.Iter()
	for iter.Next() {
		var data csvRow
		iter.CastRow(&data)
		fmt.Println(data.Name)
	}
	// Output:foo
	// bar
}

func TestLoadHeaders(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		tab, err := New(FromString(""), LoadHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(tab.Headers) != 0 {
			t.Fatalf("len(headers) want:0 got:%v", len(tab.Headers))
		}
	})
	t.Run("SimpleCase", func(t *testing.T) {
		in := `"name"
"bar"`
		tab, err := New(FromString(in), LoadHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"name"}
		if !reflect.DeepEqual(want, tab.Headers) {
			t.Fatalf("headers want:%v got:%v", want, tab.Headers)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		var out []csvRow
		if err := tab.CastAll(&out); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(out) != 1 {
			t.Fatalf("LoadHeaders content must skip first row")
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("ErrorOpts", func(t *testing.T) {
		tab, err := New(FromString(""), errorOpts())
		if tab != nil {
			t.Fatalf("tab want:nil got:%v", tab)
		}
		if err == nil {
			t.Fatalf("err want:error got:nil")
		}
	})
	t.Run("ErrorSource", func(t *testing.T) {
		_, err := New(errorSource(), LoadHeaders())
		if err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

func TestSetHeaders(t *testing.T) {
	in := "Foo"
	tab, err := New(FromString(in), SetHeaders("name"))
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	want := []string{"name"}
	if !reflect.DeepEqual(want, tab.Headers) {
		t.Fatalf("val want:%v got:%v", want, tab.Headers)
	}
	tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
	var out []csvRow
	if err := tab.CastAll(&out); err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	if len(out) == 0 {
		t.Fatalf("CSVHeaders must not skip first row")
	}
}

func TestTable_Infer(t *testing.T) {
	t.Run("SimpleCase", func(t *testing.T) {
		tab, err := New(FromString("\"name\"\nfoo\nbar"), LoadHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if err := tab.Infer(); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		var got []csvRow
		if err := tab.CastAll(&got); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []csvRow{{"foo"}, {"bar"}}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%v got:%v", want, got)
		}
	})
	t.Run("WithErrorSource", func(t *testing.T) {
		tab, err := New(errorSource())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if err := tab.Infer(); err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

func TestTable_Iter(t *testing.T) {
	t.Run("SimpleCase", func(t *testing.T) {
		tab, err := New(FromString("\"name\"\nfoo\nbar"), LoadHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if err := tab.Infer(); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		iter, err := tab.Iter()
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
		tab, err := New(errorSource())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		_, err = tab.Iter()
		if err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

func TestTable_CastAll(t *testing.T) {
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
			tab, err := New(FromString("name\nfoo\nbar"))
			if err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
			if err := tab.CastAll(&d.got); err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			want := []csvRow{{"name"}, {"foo"}, {"bar"}}
			if !reflect.DeepEqual(want, d.got) {
				t.Fatalf("val want:%v got:%v", want, d.got)
			}
		})
	}
	t.Run("MoarData", func(t *testing.T) {
		tab, err := New(FromString(`1,39,Paul
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
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "id", Type: schema.IntegerType}, {Name: "age", Type: schema.IntegerType}, {Name: "name", Type: schema.StringType}}}
		if err := tab.CastAll(&got); err != nil {
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
		tab, err := New(FromString(""))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		var got []csvRow
		if err := tab.CastAll(&got); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(got) != 0 {
			t.Fatalf("len(got) want:0 got:%v", len(got))
		}
	})
	t.Run("Error_TableWithNoSchema", func(t *testing.T) {
		tab, err := New(FromString("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if err := tab.CastAll(&[]csvRow{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("Error_OutNotAPointerToSlice", func(t *testing.T) {
		tab, err := New(FromString("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		if err := tab.CastAll([]csvRow{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
}

func TestTable_WithSchema(t *testing.T) {
	s := &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
	tab, err := New(FromString("name\nfoo\nbar"), WithSchema(s))
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	if !reflect.DeepEqual(s, tab.Schema) {
		t.Fatalf("schema want:%v got:%v", s, tab.Schema)
	}
}

func TestTable_All(t *testing.T) {
	tab, err := New(FromString("name\nfoo\nbar"))
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	want := [][]string{
		[]string{"name"},
		[]string{"foo"},
		[]string{"bar"},
	}
	got, err := tab.All()
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("schema want:%v got:%v", want, got)
	}
}
