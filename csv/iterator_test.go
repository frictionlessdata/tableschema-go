package csv

import (
	"testing"

	"github.com/matryer/is"
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
		is := is.New(t)
		iter := newIterator(stringReadCloser(""), defaultDialect, dontSkipHeaders)
		is.True(!iter.Next()) // more iterations than it should
		is.NoErr(iter.Err())
	})
}

func TestIterator_Next(t *testing.T) {
	t.Run("TwoRows", func(t *testing.T) {
		is := is.New(t)
		iter := newIterator(stringReadCloser("foo\nbar"), defaultDialect, dontSkipHeaders)
		is.True(iter.Next())  // want two more iterations
		is.True(iter.Next())  // want one more interation
		is.True(!iter.Next()) // more iterations than it should
		is.NoErr(iter.Err())
	})
	t.Run("TwoRowsSkipHeaders", func(t *testing.T) {
		is := is.New(t)
		iter := newIterator(stringReadCloser("name\nbar"), defaultDialect, skipHeaders)
		is.True(iter.Next())  // want one interation
		is.True(!iter.Next()) // more iterations than it should
		is.NoErr(iter.Err())
	})
}

func TestIterator_Row(t *testing.T) {
	t.Run("OneRow", func(t *testing.T) {
		is := is.New(t)
		iter := newIterator(stringReadCloser("name"), defaultDialect, dontSkipHeaders)
		is.True(iter.Next()) // want one iteration

		got := iter.Row()
		want := []string{"name"}
		is.Equal(want, got)
	})
}
