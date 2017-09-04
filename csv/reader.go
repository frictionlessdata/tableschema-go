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

type tableDef struct {
	Headers []string
	Source  Source
	Schema  *schema.Schema
}

// Reader provides funcionality to read a table which is backed by a CSV source.
type Reader struct {
	tableDef

	skipHeaders bool
}

// Iter returns an Iterator to read the table. Iter returns an error
// if the table physical source can not be iterated.
func (reader *Reader) Iter() (table.Iterator, error) {
	src, err := reader.Source()
	if err != nil {
		return nil, err
	}
	return newIterator(src, reader.Schema, reader.skipHeaders), nil
}

// CastAll loads and casts all rows of the table to schema types. The table
// schema must be previously assigned or inferred.
//
// The result argument must necessarily be the address for a slice. The slice
// may be nil or previously allocated.
func (reader *Reader) CastAll(out interface{}) error {
	iter, err := reader.Iter()
	if err != nil {
		return err
	}
	return table.CastAll(iter, out)
}

// All returns all rows of the table.
func (reader *Reader) All() ([][]string, error) {
	iter, err := reader.Iter()
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
type CreationOpts func(t *Reader) error

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
// CreationOpts are executed in the order they are declared.
func New(source Source, opts ...CreationOpts) (*Reader, error) {
	t := Reader{tableDef: tableDef{Source: source}, skipHeaders: false}
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
	return func(reader *Reader) error {
		iter, err := reader.Iter()
		if err != nil {
			return err
		}
		if iter.Next() {
			reader.Headers = iter.Row()
		}
		reader.skipHeaders = true
		return nil
	}
}

// SetHeaders sets the table headers.
func SetHeaders(headers ...string) CreationOpts {
	return func(reader *Reader) error {
		reader.Headers = headers
		return nil
	}
}

// WithSchema associates an schema to the CSV table being created.
func WithSchema(s *schema.Schema) CreationOpts {
	return func(reader *Reader) error {
		reader.Schema = s
		return nil
	}
}

// InferSchema tries to infer a suitable schema for the table data being read.
func InferSchema() CreationOpts {
	return func(reader *Reader) error {
		iter, err := reader.Iter()
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
		s, err := schema.Infer(reader.Headers, table)
		if err != nil {
			return err
		}
		reader.Schema = s
		return nil
	}
}

func errorOpts(headers ...string) CreationOpts {
	return func(_ *Reader) error {
		return fmt.Errorf("error opts")
	}
}
