package table

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/frictionlessdata/tableschema-go/schema"
)

// Iterator is an interface which provides method to interating over tabular
// data.
type Iterator interface {
	Next(out interface{}) bool
	Err() error
}

func newCSVIterator(source io.Reader, s *schema.Schema, skiFirstRow bool) *csvIterator {
	reader := csv.NewReader(bufio.NewReader(source))
	var err error
	if skiFirstRow {
		_, err = reader.Read()
	}
	if s == nil {
		err = fmt.Errorf("table has no schema")
	}
	return &csvIterator{
		reader: reader,
		schema: s,
		err:    err,
	}
}

type csvIterator struct {
	reader *csv.Reader
	schema *schema.Schema

	err error
}

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
func (i *csvIterator) Next(out interface{}) bool {
	// If there is an error skipping the first line, this will be catch
	// at the first call to Next(), thus for loops are not going to be
	// executed.
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

// Err returns nil if no errors happened during iteration, or the actual error
// otherwise.
func (i *csvIterator) Err() error {
	return i.err
}
