package csv

import (
	"reflect"
	"testing"
)

type iterTestValue struct {
	Name string
}

const (
	dontSkipHeaders = false
	skipHeaders     = true
)

func TestNewIterator(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		iter := newIterator(stringReadCloser(""), dontSkipHeaders)
		if iter.Next() {
			t.Fatalf("more iterations then it should.")
		}
		if iter.Err() != nil {
			t.Fatalf("err want:nil got:%v", iter.Err())
		}
	})
}

func TestIterator_Next(t *testing.T) {
	t.Run("TwoRows", func(t *testing.T) {
		iter := newIterator(stringReadCloser("foo\nbar"), dontSkipHeaders)
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
		iter := newIterator(stringReadCloser("name\nbar"), skipHeaders)
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

func TestIterator_Row(t *testing.T) {
	t.Run("OneRow", func(t *testing.T) {
		iter := newIterator(stringReadCloser("name"), dontSkipHeaders)
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
