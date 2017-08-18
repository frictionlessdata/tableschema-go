package table

import (
	"reflect"
	"strings"
	"testing"

	"github.com/frictionlessdata/tableschema-go/schema"
)

type foo struct {
	Name string
}

func TestCSV(t *testing.T) {
	t.Run("LoadCSVHeaders", func(t *testing.T) {
		in := `name
		Foo`
		tab, err := CSV(strings.NewReader(in), LoadCSVHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"name"}
		if !reflect.DeepEqual(want, tab.Headers) {
			t.Fatalf("val want:%v got:%v", want, tab.Headers)
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
		if len(out) != 0 {
			t.Fatalf("CSVHeaders must skip first row")
		}
	})
	t.Run("NoCSVHeaders", func(t *testing.T) {
		in := "Foo"
		tab, err := CSV(strings.NewReader(in), NoCSVHeaders("name"))
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
		if !reflect.DeepEqual([]foo{{"Foo"}}, out) {
			t.Fatalf("val want:%v got:%v", []foo{{"Foo"}}, out)
		}
	})
}

func TestCastAll(t *testing.T) {
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
}

func TestCastAll_Error(t *testing.T) {
	t.Run("TableWithNoSchema", func(t *testing.T) {
		tab, err := CSV(strings.NewReader("name"))
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if err := tab.CastAll(&[]foo{}); err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("OutNotAPointerToSlice", func(t *testing.T) {
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
