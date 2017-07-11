package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Schema describes tabular data.
type Schema struct {
	Fields []Field `json:"fields"`
}

// Headers returns the headers of the tabular data described
// by the schema.
func (s Schema) Headers() []string {
	var h []string
	for i := range s.Fields {
		h = append(h, s.Fields[i].Name)
	}
	return h
}

// Read reads and parses a descriptor to create a schema.
//
// Example - Reading a schema from a file:
//
//  f, err := os.Open("foo/bar/schema.json")
//  if err != nil {
//    panic(err)
//  }
//  s, err := Read(f)
//  if err != nil {
//    panic(err)
//  }
//  fmt.Println(s)
func Read(r io.Reader) (*Schema, error) {
	var s Schema
	dec := json.NewDecoder(r)
	if err := dec.Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// CastRow casts a row to schema types.  The out value must be pointer to a
// struct. Only exported fields will be cast. The lowercased field name is used
// as the key for each exported field.
//
// If a value in the row cannot be cast to its respective schema field
// (Field.CastValue), this call will return an error. Furthermore, this call
// is also going to return an error if the schema field value can not be cast
// to the struct field type.
func (s Schema) CastRow(row []string, out interface{}) error {
	if reflect.ValueOf(out).Kind() != reflect.Ptr || reflect.Indirect(reflect.ValueOf(out)).Kind() != reflect.Struct {
		return fmt.Errorf("CastRow only accepts a pointer to a struct.")
	}
	outv := reflect.Indirect(reflect.ValueOf(out))
	outt := outv.Type()
	for i := 0; i < outt.NumField(); i++ {
		fieldValue := outv.Field(i)
		if fieldValue.CanSet() { // Only consider exported fields.
			field := outt.Field(i)
			fieldName := strings.ToLower(field.Name)
			for j := range s.Fields {
				if fieldName == s.Fields[j].Name {
					cell := row[j]
					v, err := s.Fields[j].CastValue(cell)
					if err != nil {
						return err
					}
					toSetValue := reflect.ValueOf(v)
					toSetType := toSetValue.Type()
					if !toSetType.ConvertibleTo(field.Type) {
						return fmt.Errorf("value:%s field:%s - can not convert from %v to %v", fieldName, cell, toSetType, field.Type)
					}
					fieldValue.Set(toSetValue.Convert(field.Type))
					break
				}
			}
		}
	}
	return nil
}
