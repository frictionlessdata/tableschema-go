package csv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func ExampleTable_Iter() {
	table, _ := NewTable(FromString("\"name\"\nfoo\nbar"), LoadHeaders())
	iter, _ := table.Iter()
	defer iter.Close()
	for iter.Next() {
		fmt.Println(iter.Row())
	}
	// Output:[foo]
	// [bar]
}

func ExampleTable_ReadAll() {
	table, _ := NewTable(FromString("\"name\"\nfoo\nbar"), LoadHeaders())
	rows, _ := table.ReadAll()
	fmt.Print(rows)
	// Output:[[foo] [bar]]
}

func ExampleTable_ReadColumn() {
	table, _ := NewTable(FromString("name,age\nfoo,25\nbar,48"), LoadHeaders())
	cols, _ := table.ReadColumn("name")
	fmt.Print(cols)
	// Output:[foo bar]
}

func ExampleNewWriter() {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Write([]string{"foo", "bar"})
	w.Flush()
	fmt.Println(buf.String())
	// Output:foo,bar
}

func TestRemote(t *testing.T) {
	is := is.New(t)
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "\"name\"\nfoo\nbar")
	}
	ts := httptest.NewServer(http.HandlerFunc(h))
	defer ts.Close()
	table, _ := NewTable(Remote(ts.URL), LoadHeaders())
	got, _ := table.ReadAll()
	want := [][]string{{"foo"}, {"bar"}}
	is.Equal(want, got)

	t.Run("Error", func(t *testing.T) {
		is := is.New(t)
		_, err := NewTable(Remote("invalidURL"), LoadHeaders())
		is.True(err != nil)
	})
}

func TestLoadHeaders(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		is := is.New(t)
		table, err := NewTable(FromString(""), LoadHeaders())
		is.NoErr(err)
		is.Equal(len(table.Headers()), 0)
	})
	t.Run("SimpleCase", func(t *testing.T) {
		is := is.New(t)
		in := `"name"
"bar"`
		table, err := NewTable(FromString(in), LoadHeaders())
		is.NoErr(err)

		want := []string{"name"}
		is.Equal(want, table.Headers())

		iter, _ := table.Iter()
		iter.Next()
		want = []string{"bar"}
		is.Equal(want, iter.Row())
		is.True(!iter.Next())
	})
}

func TestNewTable(t *testing.T) {
	t.Run("ErrorOpts", func(t *testing.T) {
		is := is.New(t)
		table, err := NewTable(FromString(""), errorOpts())
		is.True(table == nil)
		is.True(err != nil)
	})
	t.Run("ErrorSource", func(t *testing.T) {
		is := is.New(t)
		_, err := NewTable(errorSource(), LoadHeaders())
		is.True(err != nil)
	})
}

func TestSetHeaders(t *testing.T) {
	is := is.New(t)
	in := "Foo"
	table, err := NewTable(FromString(in), SetHeaders("name"))
	is.NoErr(err)
	want := []string{"name"}
	is.Equal(want, table.Headers())

	iter, _ := table.Iter()
	iter.Next()
	want = []string{"Foo"}
	is.Equal(want, iter.Row())
	is.True(!iter.Next())
}

func TestDelimiter(t *testing.T) {
	is := is.New(t)
	in := "Foo;Bar"
	table, err := NewTable(FromString(in), Delimiter(';'))
	is.NoErr(err)
	contents, err := table.ReadAll()
	is.NoErr(err)
	is.Equal(contents, [][]string{{"Foo", "Bar"}})
}

func TestConsiderInitialSpace(t *testing.T) {
	is := is.New(t)
	in := " Foo"
	table, err := NewTable(FromString(in), ConsiderInitialSpace())
	is.NoErr(err)
	contents, err := table.ReadAll()
	is.NoErr(err)
	is.Equal(contents, [][]string{{" Foo"}})
}

func TestReadAll(t *testing.T) {
	is := is.New(t)
	in := "name\nfoo\nbar"
	want := [][]string{[]string{"name"}, []string{"foo"}, []string{"bar"}}

	table, err := NewTable(FromString(in))
	is.NoErr(err)
	rows, err := table.ReadAll()
	is.NoErr(err)
	is.Equal(want, rows)

	table, err = NewTable(
		func() (io.ReadCloser, error) {
			return nil, errors.New("this is a source test error")
		})
	is.NoErr(err)
	_, err = table.ReadAll()
	is.True(err != nil)
}

func TestString(t *testing.T) {
	is := is.New(t)
	in := "name\nfoo\nbar"
	want := "name\nfoo\nbar\n"
	table, err := NewTable(FromString(in))
	is.NoErr(err)
	is.Equal(want, table.String())
}

func TestReadColumn(t *testing.T) {
	t.Run("HeaderNotFound", func(t *testing.T) {
		is := is.New(t)
		tab, err := NewTable(FromString("name\nfoo"), LoadHeaders())
		is.NoErr(err)
		_, err = tab.ReadColumn("age")
		is.True(err != nil) // Must err as there is no column called age.
	})
	t.Run("ErrorCreatingIter", func(t *testing.T) {
		is := is.New(t)
		tab, err := NewTable(errorSource())
		is.NoErr(err)
		tab.headers = []string{"age"}
		_, err = tab.ReadColumn("age")
		is.True(err != nil) // Must err as the source will error.
	})
}
