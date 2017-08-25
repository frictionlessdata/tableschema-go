package csv

import (
	"reflect"
	"strings"
	"testing"

	"github.com/frictionlessdata/tableschema-go/schema"
)

type iterTestValue struct {
	Name string
}

const (
	dontSkipHeaders = false
	skipHeaders     = true
)

func TestNewIterator_EmptyString(t *testing.T) {
	iter := newIterator(
		strings.NewReader(""),
		&schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}},
		dontSkipHeaders,
	)
	if iter.Next() {
		t.Fatalf("more iterations then it should.")
	}
	if iter.Err() != nil {
		t.Fatalf("err want:nil got:%v", iter.Err())
	}
}

func TestIterator_Next(t *testing.T) {
	t.Run("TwoRows", func(t *testing.T) {
		iter := newIterator(strings.NewReader("foo\nbar"), nil, dontSkipHeaders)
		if !iter.Next() {
			t.Fatalf("want two more iterations.")
		}
		if !iter.Next() {
			t.Fatalf("want one more iteration")
		}
		if iter.Next() {
			t.Fatalf("more iterations then it should.")
		}
		if iter.Err() != nil {
			t.Fatalf("err want:nil got:%v", iter.Err())
		}
	})
	t.Run("TwoRowsSkipHeaders", func(t *testing.T) {
		iter := newIterator(strings.NewReader("name\nbar"), nil, skipHeaders)
		if !iter.Next() {
			t.Fatalf("want one iteration")
		}
		if iter.Next() {
			t.Fatalf("more iterations then it should.")
		}
		if iter.Err() != nil {
			t.Fatalf("err want:nil got:%v", iter.Err())
		}
	})
}

func TestIterator_CastRow(t *testing.T) {
	t.Run("OneRow", func(t *testing.T) {
		iter := newIterator(
			strings.NewReader("name"),
			&schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}},
			dontSkipHeaders,
		)
		if !iter.Next() {
			t.Fatalf("want:one iteration got:zero")
		}
		var got iterTestValue
		if err := iter.CastRow(&got); err != nil {
			t.Fatalf("err want:nil got:%v", iter.Err())
		}
		want := iterTestValue{"name"}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%v got:%v", want, got)
		}
	})
	t.Run("Error_NilSchema", func(t *testing.T) {
		iter := newIterator(strings.NewReader("name"), nil, dontSkipHeaders)
		if !iter.Next() {
			t.Fatalf("next want:true got:false")
		}
		var got iterTestValue
		if err := iter.CastRow(&got); err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
	t.Run("Error_NilOutput", func(t *testing.T) {
		iter := newIterator(strings.NewReader("name"), nil, dontSkipHeaders)
		if !iter.Next() {
			t.Fatalf("next want:true got:false")
		}
		if err := iter.CastRow(nil); err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

func TestIterator_Row(t *testing.T) {
	t.Run("OneRow", func(t *testing.T) {
		iter := newIterator(
			strings.NewReader("name"),
			&schema.Schema{Fields: []schema.Field{{Name: "name", Type: schema.StringType}}},
			dontSkipHeaders,
		)
		if !iter.Next() {
			t.Fatalf("want:one iteration got:zero")
		}
		got := iter.Row()
		want := []string{"name"}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("val want:%v got:%v", want, got)
		}
	})
}
