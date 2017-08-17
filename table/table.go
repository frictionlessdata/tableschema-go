package table

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/frictionlessdata/tableschema-go/schema"
)

// CreationOpts defines functional options for creating Tables.
type CreationOpts func(t *Table) error

// CSV creates a Table from the CSV physical representation.
func CSV(source io.Reader, opts ...CreationOpts) (Table, error) {
	t := Table{Source: source}
	for _, opt := range opts {
		if err := opt(&t); err != nil {
			return Table{}, err
		}
	}
	return t, nil
}

// LoadCSVHeaders uses the first line of the CSV as table headers.
func LoadCSVHeaders() CreationOpts {
	return func(t *Table) error {
		r := csv.NewReader(t.Source)
		record, err := r.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		t.Headers = record
		t.skipFirstRow = true
		return nil
	}
}

// CSVFile creates a Table from a CSV local file.
func CSVFile(path string, opts ...CreationOpts) (Table, error) {
	f, err := os.Open(path)
	if err != nil {
		return Table{}, err
	}
	return CSV(f, opts...)
}

// CSVHeaders sets the table headers. It would override headers
// that might exist in the first line of the CSV.
func CSVHeaders(headers ...string) CreationOpts {
	return func(t *Table) error {
		t.Headers = headers
		t.skipFirstRow = true
		return nil
	}
}

// NoCSVHeaders signals the reading library that your CSV has no headers
// defined the first line.
func NoCSVHeaders(headers ...string) CreationOpts {
	return func(t *Table) error {
		t.Headers = headers
		return nil
	}
}

// Table makes it easy to deal with tabular data.
type Table struct {
	Headers []string
	Source  io.Reader
	Schema  *schema.Schema

	skipFirstRow bool
}

// CastAll loads and casts all rows of the table to schema types. The table
// schema must be previously assigned or inferred.
//
// The result argument must necessarily be the address for a slice. The slice
// may be nil or previously allocated.
func (t *Table) CastAll(out interface{}) error {
	if t.Schema == nil {
		return fmt.Errorf("table has no schema")
	}
	r := csv.NewReader(t.Source)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	if t.skipFirstRow {
		records = records[1:]
	}
	outv := reflect.ValueOf(out)
	if outv.Kind() != reflect.Ptr || outv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("out argument must be a slice address")
	}
	slicev := outv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	for i, record := range records {
		if slicev.Len() == i {
			elemp := reflect.New(elemt)
			if err := t.Schema.CastRow(record, elemp.Interface()); err != nil {
				return err
			}
			slicev = reflect.Append(slicev, elemp.Elem())
			slicev = slicev.Slice(0, slicev.Cap())
		} else {
			if err := t.Schema.CastRow(record, slicev.Index(i).Addr().Interface()); err != nil {
				return err
			}
		}
	}
	outv.Elem().Set(slicev.Slice(0, len(records)))
	return nil
}
