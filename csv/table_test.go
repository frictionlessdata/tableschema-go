package csv

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

type csvRow struct {
	Name string
}

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
