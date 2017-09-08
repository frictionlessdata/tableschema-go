package csv

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
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

func TestLoadHeaders(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		table, err := NewTable(FromString(""), LoadHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		if len(table.Headers()) != 0 {
			t.Fatalf("len(headers) want:0 got:%v", len(table.Headers()))
		}
	})
	t.Run("SimpleCase", func(t *testing.T) {
		in := `"name"
"bar"`
		table, err := NewTable(FromString(in), LoadHeaders())
		if err != nil {
			t.Fatalf("err want:nil got:%q", err)
		}
		want := []string{"name"}
		if !reflect.DeepEqual(want, table.Headers()) {
			t.Fatalf("headers want:%v got:%v", want, table.Headers())
		}

		iter, _ := table.Iter()
		iter.Next()
		want = []string{"bar"}
		if !reflect.DeepEqual(want, iter.Row()) {
			t.Fatalf("headers want:%v got:%v", want, iter.Row())
		}
		if iter.Next() {
			t.Fatalf("want:no more iterations")
		}
	})
}

func TestNewTable(t *testing.T) {
	t.Run("ErrorOpts", func(t *testing.T) {
		table, err := NewTable(FromString(""), errorOpts())
		if table != nil {
			t.Fatalf("reader want:nil got:%v", table)
		}
		if err == nil {
			t.Fatalf("err want:error got:nil")
		}
	})
	t.Run("ErrorSource", func(t *testing.T) {
		_, err := NewTable(errorSource(), LoadHeaders())
		if err == nil {
			t.Fatalf("want:err got:nil")
		}
	})
}

func TestSetHeaders(t *testing.T) {
	in := "Foo"
	table, err := NewTable(FromString(in), SetHeaders("name"))
	if err != nil {
		t.Fatalf("err want:nil got:%q", err)
	}
	want := []string{"name"}
	if !reflect.DeepEqual(want, table.Headers()) {
		t.Fatalf("val want:%v got:%v", want, table.Headers())
	}
	iter, _ := table.Iter()
	iter.Next()
	want = []string{"Foo"}
	if !reflect.DeepEqual(want, iter.Row()) {
		t.Fatalf("headers want:%v got:%v", want, iter.Row())
	}
	if iter.Next() {
		t.Fatalf("want:no more iterations")
	}
}
