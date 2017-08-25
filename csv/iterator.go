package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/frictionlessdata/tableschema-go/schema"
)

func newIterator(source io.Reader, s *schema.Schema, skipHeaders bool) *csvIterator {
	return &csvIterator{
		reader:      csv.NewReader(source),
		schema:      s,
		skipHeaders: skipHeaders,
	}
}

type csvIterator struct {
	reader *csv.Reader
	schema *schema.Schema

	current     []string
	err         error
	skipHeaders bool
}

func (i *csvIterator) Next() bool {
	if i.err != nil {
		return false
	}
	var err error
	i.current, err = i.reader.Read()
	if err != io.EOF {
		i.err = err
	}
	if i.skipHeaders {
		i.skipHeaders = false
		i.Next()
	}
	return err == nil
}

func (i *csvIterator) CastRow(out interface{}) error {
	if i.schema == nil {
		return fmt.Errorf("table has no schema")
	}
	return i.schema.CastRow(i.current, out)
}

func (i *csvIterator) Row() []string {
	return i.current
}

func (i *csvIterator) Err() error {
	return i.err
}
