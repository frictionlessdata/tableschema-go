package schema

import (
	"testing"

	"github.com/matryer/is"
)

func TestCastAny(t *testing.T) {
	is := is.New(t)
	got, err := castAny("foo")
	is.NoErr(err)
	is.Equal("foo", got)
}

func TestEncodeAny(t *testing.T) {
	is := is.New(t)
	got, err := encodeAny(10)
	is.NoErr(err)
	is.Equal("10", got)
}
