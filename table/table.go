package table

import (
	"io"

	"github.com/frictionlessdata/tableschema-go/schema"
)

// RawIter iterates over a set of rows in its raw form.
type RawIter struct {
}

// Next retrieves the next row from the table, blocking if necessary.
// The row is returned in its raw form, as a slice of strings.
func (i *RawIter) Next() ([]string, error) {
	return nil, nil
}

// Done returns true only if a follow up Next call is guaranteed to return false.
func (i *RawIter) Done() bool {
	return false
}

// RowIter iterates over a set of rows in a table, casting its rows to a
// specified schema.
type RowIter struct {
	Schema schema.Schema
}

// Next retrieves the next row from the table, blocking if necessary.
// The row is returned as a struct, cast to using the Schema.
// See: http://godoc.org/github.com/frictionlessdata/tableschema-go/schema#Schema.CastRow
func (i *RowIter) Next(out interface{}) error {
	return nil
}

// Table represents a tabular data structure. Tabular data consists of a set of rows.
// Each row has a set of fields (columns). We usually expect that each row has
// the same set of fields and thus we can talk about the fields for the table as a whole.
// More at: https://specs.frictionlessdata.io/table-schema/#concepts
type Table struct {
}

// Raw allows to iterate over the table in its raw form (values as strings).
// No validation or cast is performed.
func (t *Table) Raw() (*RawIter, error) {
	return nil, nil
}

// Row allows to iterate over the table casting its rows to the passed-in
// schema.
func (t *Table) Row(schema schema.Schema) (*RowIter, error) {
	return nil, nil
}

// CSV creates a Table from the CSV physical representation.
func CSV(source io.Reader) (*Table, error) {
	return nil, nil
}

// JSON creates a Table from the JSON physical representation.
func JSON(source io.Reader) (*Table, error) {
	return nil, nil
}

// ([][]string, error) {
// 	return csv.NewReader(bufio.NewReader(source)).ReadAll()
// }
