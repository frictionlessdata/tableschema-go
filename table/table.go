package table

import (
	"encoding/csv"
	"io"
)

// Table tabular data representation.
type Table struct {
	Source io.Reader
}

// ReadAll reads all the remaining records from table's Source. Each record is
// a slice of fields.
func (t *Table) Read() ([][]string, error) {
	return csv.NewReader(t.Source).ReadAll()
}
