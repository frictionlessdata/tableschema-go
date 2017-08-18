package table

import (
	"fmt"
	"io"
	"reflect"

	"github.com/frictionlessdata/tableschema-go/schema"
)

// CreationOpts defines functional options for creating Tables.
type CreationOpts func(t *Table) error

// Table makes it easy to deal with tabular data.
type Table struct {
	Headers []string
	Source  io.Reader
	Schema  *schema.Schema
}

// CastAll loads and casts all rows of the table to schema types. The table
// schema must be previously assigned or inferred.
//
// The result argument must necessarily be the address for a slice. The slice
// may be nil or previously allocated.
func (t *Table) CastAll(out interface{}) error {
	outv := reflect.ValueOf(out)
	if outv.Kind() != reflect.Ptr || outv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("out argument must be a slice address")
	}
	slicev := outv.Elem()
	slicev = slicev.Slice(0, 0) // Trucantes the passed-in slice.
	elemt := slicev.Type().Elem()
	iter := t.Iter()
	i := 0
	for elemp := reflect.New(elemt); iter.Next(elemp.Interface()); {
		slicev = reflect.Append(slicev, elemp.Elem())
		slicev = slicev.Slice(0, slicev.Cap())
		i++
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	outv.Elem().Set(slicev.Slice(0, i))
	return nil
}

// Iter returns an Iterator to read the table.
func (t *Table) Iter() Iterator {
	return newCSVIterator(t.Source, t.Schema)
}
