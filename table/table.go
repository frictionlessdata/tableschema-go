package table

import (
	"fmt"
	"reflect"
)

// Table makes it easy to deal with physical tabular data.
type Table interface {
	// Iter returns an Iterator to read the table. Iter returns an error
	// if the table physical source can not be iterated.
	Iter() (Iterator, error)

	// Infer tries to infer a suitable schema for the table.
	Infer() error

	// CastAll loads and casts all rows of the table to schema types. The table
	// schema must be previously assigned or inferred.
	//
	// The result argument must necessarily be the address for a slice. The slice
	// may be nil or previously allocated.
	CastAll(out interface{}) error

	// All returns all rows of the table.
	All() ([][]string, error)
}

// CastAll loads and casts all rows returned by the iterator to schema types.
//
// The result argument must necessarily be the address for a slice. The slice
// may be nil or previously allocated.
func CastAll(iter Iterator, out interface{}) error {
	outv := reflect.ValueOf(out)
	if outv.Kind() != reflect.Ptr || outv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("out argument must be a slice address")
	}
	slicev := outv.Elem()
	slicev = slicev.Slice(0, 0) // Trucantes the passed-in slice.
	elemt := slicev.Type().Elem()
	i := 0
	for iter.Next() {
		elemp := reflect.New(elemt)
		if err := iter.CastRow(elemp.Interface()); err != nil {
			return err
		}
		slicev = reflect.Append(slicev, elemp.Elem())
		slicev = slicev.Slice(0, slicev.Len())
		i++
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	outv.Elem().Set(slicev.Slice(0, i))
	return nil
}
