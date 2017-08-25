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
	tab, _ := New(StringSource("\"name\"\nfoo\nbar"), LoadHeaders())
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
	tab, _ := New(StringSource("\"name\"\nfoo\nbar"), LoadHeaders())
	if err := tab.Infer(); err != nil {
		fmt.Println(err)
	}
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
		tab, err := New(StringSource(""), LoadHeaders())
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
		tab, err := New(StringSource(in), LoadHeaders())
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

func TestSetHeaders(t *testing.T) {
	in := "Foo"
	tab, err := New(StringSource(in), SetHeaders("name"))
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
			tab, err := New(StringSource("name\nfoo\nbar"))
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
	t.Run("EmptyString", func(t *testing.T) {
		tab, err := New(StringSource(""))
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
		tab, err := New(StringSource("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if err := tab.CastAll(&[]csvRow{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("Error_OutNotAPointerToSlice", func(t *testing.T) {
		tab, err := New(StringSource("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		if err := tab.CastAll([]csvRow{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
}
