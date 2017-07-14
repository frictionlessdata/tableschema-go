package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// PrimaryKeys defines a schema primary key.
// More: https://specs.frictionlessdata.io/table-schema/#primary-key
type PrimaryKeys []string

type ForeignKeyReference struct {
	Resource string   `json:"resource"`
	Fields   []string `json:"fields"`
}

// ForeignKeys defines a schema foreign key
// More: https://specs.frictionlessdata.io/table-schema/#foreign-keys
type ForeignKeys struct {
	Fields    []string            `json:"fields"`
	Reference ForeignKeyReference `json:"reference"`
}

// Schema describes tabular data.
type Schema struct {
	Fields      []Field     `json:"fields"`
	PrimaryKeys PrimaryKeys `json:"primaryKey"`
	ForeignKeys ForeignKeys `json:"foreignKeys"`
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

// GetField fetches the index and field referenced by the name argument. The third
// return value is true if there is a field with the passed-in name in
// the schema and false otherwise.
func (s Schema) GetField(name string) (int, *Field, bool) {
	for i := range s.Fields {
		if name == s.Fields[i].Name {
			return i, &s.Fields[i], true
		}
	}
	return 0, nil, false
}

// Validate checks whether the schema is valid. If it is not, returns an error
// describing the problem.
// More at: https://specs.frictionlessdata.io/table-schema/
func (s Schema) Validate() error {
	// Checking if all fields have a name.
	for _, f := range s.Fields {
		if f.Name == "" {
			return fmt.Errorf("invalid field: attribute name is mandatory")
		}
	}
	// Checking primary keys.
	for _, pk := range s.PrimaryKeys {
		if _, _, ok := s.GetField(pk); !ok {
			return fmt.Errorf("invalid primary key: there is no field %s", pk)
		}
	}
	// Checking foreign keys.
	for _, fk := range s.ForeignKeys.Fields {
		if _, _, ok := s.GetField(fk); !ok {
			return fmt.Errorf("invalid foreign keys: there is no field %s", fk)
		}
	}
	if len(s.ForeignKeys.Reference.Fields) != len(s.ForeignKeys.Fields) {
		return fmt.Errorf("invalid foreign key: foreignKey.fields must contain the same number entries as foreignKey.reference.fields")
	}
	return nil
}

// Read reads, parses and validates a descriptor to create a schema.
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
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return &s, nil
}

// CastRow casts a row to schema types. The out value must be pointer to a
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
			fieldIndex, f, ok := s.GetField(fieldName)
			if ok {
				cell := row[fieldIndex]
				v, err := f.CastValue(cell)
				if err != nil {
					return err
				}
				toSetValue := reflect.ValueOf(v)
				toSetType := toSetValue.Type()
				if !toSetType.ConvertibleTo(field.Type) {
					return fmt.Errorf("value:%s field:%s - can not convert from %v to %v", fieldName, cell, toSetType, field.Type)
				}
				fieldValue.Set(toSetValue.Convert(field.Type))
			}
		}
	}
	return nil
}
