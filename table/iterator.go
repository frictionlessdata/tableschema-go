package table

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/frictionlessdata/tableschema-go/schema"
)

// Iterator is an interface which provides method to interating over tabular
// data.
type Iterator interface {
	// Next reads the next row from the input source, blocking if necessary. It automatically buffer some data,
	// improving reading performance.
	// Next returns true if the row was successfully unmarshalled onto result, and false at
	// the end of the table or if an error happened.
	//
	// For example:
	//     iter := myTable.Iter()
	//     for iter.Next(&result) {
	//         log.Printf("Result: %v\n", result.Name)
	//     }
	//     if iter.Err() != nil {
	//         log.Fatal(iter.Err())
	//     }
	Next(out interface{}) bool

	// Err returns nil if no errors happened during iteration, or the actual error
	// otherwise.
	Err() error
}

func newCSVIterator(source io.Reader, s *schema.Schema) *csvIterator {
	if s == nil {
		return &csvIterator{
			reader: nil,
			schema: nil,
			err:    fmt.Errorf("table has no schema"),
		}
	}
	reader := csv.NewReader(source)
	return &csvIterator{
		reader: reader,
		schema: s,
	}
}

type csvIterator struct {
	reader *csv.Reader
	schema *schema.Schema

	err error
}

func (i *csvIterator) Next(out interface{}) bool {
	var next []string
	if i.err == nil {
		next, i.err = i.reader.Read()
	}
	switch i.err {
	case nil:
		if i.err = i.schema.CastRow(next, out); i.err != nil {
			return false
		}
		return true
	case io.EOF:
		i.err = nil
		return false
	default:
		return false
	}
}

func (i *csvIterator) Err() error {
	return i.err
}
