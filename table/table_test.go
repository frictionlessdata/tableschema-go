package table

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/frictionlessdata/tableschema-go/schema"
)

type foo struct {
	Name string
}

func ExampleTable_Iter() {
	tab, _ := CSV(strings.NewReader("\"name\"\nfoo\nbar"), LoadCSVHeaders())
	tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
	iter := tab.Iter()
	var data foo
	for iter.Next(&data) {
		fmt.Println(data.Name)
	}
	// Output:foo
	// bar
}

func TestCSV(t *testing.T) {
	t.Run("LoadCSVHeadersEmptyString", func(t *testing.T) {
		tab, err := CSV(strings.NewReader(""), LoadCSVHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(tab.Headers) != 0 {
			t.Fatalf("len(headers) want:0 got:%v", len(tab.Headers))
		}
	})
	t.Run("LoadCSVHeadersNoComments", func(t *testing.T) {
		in := `"name"
"bar"`
		tab, err := CSV(strings.NewReader(in), LoadCSVHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"name"}
		if !reflect.DeepEqual(want, tab.Headers) {
			t.Fatalf("headers want:%v got:%v", want, tab.Headers)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		var out []foo
		if err := tab.CastAll(&out); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(out) != 1 {
			t.Fatalf("LoadCSVHeaders content must skip first row")
		}
	})
	t.Run("LoadCSVHeadersWithComments", func(t *testing.T) {
		in := `# Foo
# Bar Bez Boo
"name"
"bar"`
		tab, err := CSV(strings.NewReader(in), LoadCSVHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"name"}
		if !reflect.DeepEqual(want, tab.Headers) {
			t.Fatalf("headers want:%v got:%v", want, tab.Headers)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		var out []foo
		if err := tab.CastAll(&out); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(out) != 1 {
			t.Fatalf("LoadCSVHeaders content must skip first row")
		}
	})
	t.Run("CSVHeaders", func(t *testing.T) {
		in := "Foo"
		tab, err := CSV(strings.NewReader(in), CSVHeaders("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"name"}
		if !reflect.DeepEqual(want, tab.Headers) {
			t.Fatalf("val want:%v got:%v", want, tab.Headers)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		var out []foo
		if err := tab.CastAll(&out); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(out) == 0 {
			t.Fatalf("CSVHeaders must not skip first row")
		}
	})
}

func TestTable_CastAll(t *testing.T) {
	data := []struct {
		desc string
		got  []foo
	}{
		{"OutEmpty", []foo{}},
		{"OutNil", nil},
		{"OutInitialized", []foo{{"fooooo"}}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			tab, err := CSV(strings.NewReader("name\nfoo\nbar"))
			if err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
			if err := tab.CastAll(&d.got); err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			want := []foo{{"name"}, {"foo"}, {"bar"}}
			if !reflect.DeepEqual(want, d.got) {
				t.Fatalf("val want:%v got:%v", want, d.got)
			}
		})
	}
	t.Run("EmptyString", func(t *testing.T) {
		tab, err := CSV(strings.NewReader(""))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		var got []foo
		if err := tab.CastAll(&got); err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(got) != 0 {
			t.Fatalf("len(got) want:0 got:%v", len(got))
		}
	})
	t.Run("Error_TableWithNoSchema", func(t *testing.T) {
		tab, err := CSV(strings.NewReader("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if err := tab.CastAll(&[]foo{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("Error_OutNotAPointerToSlice", func(t *testing.T) {
		tab, err := CSV(strings.NewReader("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		tab.Schema = &schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}}
		if err := tab.CastAll([]foo{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
}
