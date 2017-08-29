package csv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/frictionlessdata/tableschema-go/schema"
	"github.com/frictionlessdata/tableschema-go/table"
)

// Maximum number of rows used to infer schema.
const maxNumRowsInfer = 100

// Table implements a Table which is backed by a CSV source.
type Table struct {
	Headers []string
	Source  Source
	Schema  *schema.Schema

	skipHeaders bool
}

// Iter returns an Iterator to read the table. Iter returns an error
// if the table physical source can not be iterated.
func (t *Table) Iter() (table.Iterator, error) {
	reader, err := t.Source()
	if err != nil {
		return nil, err
	}
	return newIterator(reader, t.Schema, t.skipHeaders), nil
}

// Infer tries to infer a suitable schema for the table.
func (t *Table) Infer() error {
	iter, err := t.Iter()
	if err != nil {
		return err
	}
	var table [][]string
	for i := 0; i < maxNumRowsInfer; i++ {
		if !iter.Next() {
			break
		}

		table = append(table, iter.Row())
	}
	s, err := schema.Infer(t.Headers, table)
	if err != nil {
		return err
	}
	t.Schema = s
	return nil
}

// CastAll loads and casts all rows of the table to schema types. The table
// schema must be previously assigned or inferred.
//
// The result argument must necessarily be the address for a slice. The slice
// may be nil or previously allocated.
func (t *Table) CastAll(out interface{}) error {
	iter, err := t.Iter()
	if err != nil {
		return err
	}
	return table.CastAll(iter, out)
}

// All returns all rows of the table.
func (t *Table) All() ([][]string, error) {
	iter, err := t.Iter()
	if err != nil {
		return nil, err
	}
	var all [][]string
	for iter.Next() {
		all = append(all, iter.Row())
	}
	return all, nil
}

// CreationOpts defines functional options for creating Tables.
type CreationOpts func(t *Table) error

// Source defines a table physical data source.
type Source func() (io.Reader, error)

// FromFile defines a file-based Source.
func FromFile(path string) Source {
	return func() (io.Reader, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return bufio.NewReader(f), nil
	}
}

// FromString defines a string-based source.
func FromString(str string) Source {
	return func() (io.Reader, error) {
		return strings.NewReader(str), nil
	}
}

func errorSource() Source {
	return func() (io.Reader, error) {
		return nil, fmt.Errorf("error source")
	}
}

// New creates a Table from the CSV physical representation.
func New(source Source, opts ...CreationOpts) (*Table, error) {
	t := Table{Source: source}
	for _, opt := range opts {
		if err := opt(&t); err != nil {
			return nil, err
		}
	}
	return &t, nil
}

// LoadHeaders uses the first line of the CSV as table headers.
// The header line will be skipped during iteration
func LoadHeaders() CreationOpts {
	return func(t *Table) error {
		iter, err := t.Iter()
		if err != nil {
			return err
		}
		if iter.Next() {
			t.Headers = iter.Row()
		}
		t.skipHeaders = true
		return nil
	}
}

// SetHeaders sets the table headers.
func SetHeaders(headers ...string) CreationOpts {
	return func(t *Table) error {
		t.Headers = headers
		return nil
	}
}

// WithSchema associates an schema to the CSV table being created.
func WithSchema(s *schema.Schema) CreationOpts {
	return func(t *Table) error {
		t.Schema = s
		return nil
	}
}

func errorOpts(headers ...string) CreationOpts {
	return func(t *Table) error {
		return fmt.Errorf("error opts")
	}
}
