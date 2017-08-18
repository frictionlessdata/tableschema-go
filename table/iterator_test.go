package table

import (
	"reflect"
	"strings"
	"testing"

	"github.com/frictionlessdata/tableschema-go/schema"
)

type iterTestValue struct {
	Name string
}

func TestNext(t *testing.T) {
	data := []struct {
		desc         string
		want         []iterTestValue
		skipFirstRow bool
	}{
		{"AllTable", []iterTestValue{{"name"}, {"foo"}, {"bar"}}, false},
		{"SkipFirstRow", []iterTestValue{{"foo"}, {"bar"}}, true},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			iter := newCSVIterator(
				strings.NewReader("name\nfoo\nbar"),
				&schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}},
				d.skipFirstRow,
			)
			for _, want := range d.want {
				var got iterTestValue
				if iter.Next(&got) {
					if !reflect.DeepEqual(want, got) {
						t.Fatalf("val want:%v got:%v", want, got)
					}
				}
				if iter.Err() != nil {
					t.Fatalf("err want:nil got:%v", iter.Err())
				}
			}
			var nothing iterTestValue
			if iter.Next(&nothing) {
				t.Fatalf("more iterations then it should.")
			}
			if iter.Err() != nil {
				t.Fatalf("err want:nil got:%v", iter.Err())
			}
		})
	}

}

func TestNext_Error(t *testing.T) {
	t.Run("NilSchema", func(t *testing.T) {
		iter := newCSVIterator(strings.NewReader("name"), nil, true)
		if iter.Err() == nil {
			t.Fatalf("want:err got:nil")
		}
	})
	t.Run("ErrorCasting", func(t *testing.T) {
		iter := newCSVIterator(strings.NewReader("name"), nil, true)
		if iter.Next(nil) {
			t.Fatalf("want:err got:nil")
		}
	})
}
