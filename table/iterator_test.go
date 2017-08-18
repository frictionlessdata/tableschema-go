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

func TestIterator_Next(t *testing.T) {
	data := []struct {
		desc string
		want []iterTestValue
	}{
		{"AllTable", []iterTestValue{{"name"}, {"foo"}, {"bar"}}},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			iter := newCSVIterator(
				strings.NewReader("name\nfoo\nbar"),
				&schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}},
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
	t.Run("EmptyString", func(t *testing.T) {
		iter := newCSVIterator(
			strings.NewReader(""),
			&schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}},
		)
		var nothing iterTestValue
		if iter.Next(&nothing) {
			t.Fatalf("more iterations then it should.")
		}
		if iter.Err() != nil {
			t.Fatalf("err want:nil got:%v", iter.Err())
		}
	})
	t.Run("Error_NilSchema", func(t *testing.T) {
		iter := newCSVIterator(strings.NewReader("name"), nil)
		if iter.Err() == nil {
			t.Fatalf("want:err got:nil")
		}
	})
	t.Run("Error_Casting", func(t *testing.T) {
		iter := newCSVIterator(strings.NewReader("name"), nil)
		if iter.Next(nil) {
			t.Fatalf("want:err got:nil")
		}
	})
}
