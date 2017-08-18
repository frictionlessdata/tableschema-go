package table

import (
	"reflect"
	"strings"
	"testing"

	"github.com/frictionlessdata/tableschema-go/schema"
)

func TestLoadCSVHeaders(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		tab, err := CSV(strings.NewReader(""), LoadCSVHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(tab.Headers) != 0 {
			t.Fatalf("len(headers) want:0 got:%v", len(tab.Headers))
		}
	})
	t.Run("NoComments", func(t *testing.T) {
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
	t.Run("Comments", func(t *testing.T) {
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
}

func TestCSVHeaders(t *testing.T) {
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
}
