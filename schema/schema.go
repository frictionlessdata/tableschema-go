package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/frictionlessdata/tableschema-go/table"
)

// InvalidPosition is returned by GetField call when
// it refers to a field that does not exist in the schema.
const InvalidPosition = -1

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

// ReadFromFile reads and parses a schema descrptor from a local file.
func ReadFromFile(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return Read(f)
}

// Fields represents a list of schema fields.
type Fields []Field

func (f Fields) Len() int           { return len(f) }
func (f Fields) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f Fields) Less(i, j int) bool { return strings.Compare(f[i].Name, f[j].Name) == -1 }

// ForeignKeyReference represents the field reference by a foreign key.
type ForeignKeyReference struct {
	Resource          string      `json:"resource,omitempty"`
	Fields            []string    `json:"-"`
	FieldsPlaceholder interface{} `json:"fields,omitempty"`
}

// ForeignKeys defines a schema foreign key
type ForeignKeys struct {
	Fields            []string            `json:"-"`
	FieldsPlaceholder interface{}         `json:"fields,omitempty"`
	Reference         ForeignKeyReference `json:"reference,omitempty"`
}

// Schema describes tabular data.
type Schema struct {
	Fields                Fields      `json:"fields,omitempty"`
	PrimaryKeyPlaceholder interface{} `json:"primaryKey,omitempty"`
	PrimaryKeys           []string    `json:"-"`
	ForeignKeys           ForeignKeys `json:"foreignKeys,omitempty"`
	MissingValues         []string    `json:"missingValues,omitempty"`
}

// Headers returns the headers of the tabular data described
// by the schema.
func (s *Schema) Headers() []string {
	var h []string
	for i := range s.Fields {
		h = append(h, s.Fields[i].Name)
	}
	return h
}

// GetField fetches the index and field referenced by the name argument.
func (s *Schema) GetField(name string) (*Field, int) {
	for i := range s.Fields {
		if name == s.Fields[i].Name {
			return &s.Fields[i], i
		}
	}
	return nil, InvalidPosition
}

// HasField returns checks whether the schema has a field with the passed-in.
func (s *Schema) HasField(name string) bool {
	_, pos := s.GetField(name)
	return pos != InvalidPosition
}

// Validate checks whether the schema is valid. If it is not, returns an error
// describing the problem.
// More at: https://specs.frictionlessdata.io/table-schema/
func (s *Schema) Validate() error {
	// Checking if all fields have a name.
	for _, f := range s.Fields {
		if f.Name == "" {
			return fmt.Errorf("invalid field: attribute name is mandatory")
		}
	}
	// Checking primary keys.
	for _, pk := range s.PrimaryKeys {
		if !s.HasField(pk) {
			return fmt.Errorf("invalid primary key: there is no field %s", pk)
		}
	}
	// Checking foreign keys.
	for _, fk := range s.ForeignKeys.Fields {
		if !s.HasField(fk) {
			return fmt.Errorf("invalid foreign keys: there is no field %s", fk)
		}
	}
	if len(s.ForeignKeys.Reference.Fields) != len(s.ForeignKeys.Fields) {
		return fmt.Errorf("invalid foreign key: foreignKey.fields must contain the same number entries as foreignKey.reference.fields")
	}
	return nil
}

// Write writes the schema descriptor.
func (s *Schema) Write(w io.Writer) error {
	pp, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}
	w.Write(pp)
	return nil
}

// SaveToFile writes the schema descriptor in local file.
func (s *Schema) SaveToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return s.Write(f)
}

// Decode decodes the passed-in row to schema types and stores it in the value pointed
// by out. The out value must be pointer to a struct. Only exported fields will be unmarshalled.
// The lowercased field name is used as the key for each exported field.
//
// If a value in the row cannot be marshalled to its respective schema field (Field.Unmarshal),
// this call will return an error. Furthermore, this call is also going to return an error if
// the schema field value can not be unmarshalled to the struct field type.
func (s *Schema) Decode(row []string, out interface{}) error {
	if reflect.ValueOf(out).Kind() != reflect.Ptr || reflect.Indirect(reflect.ValueOf(out)).Kind() != reflect.Struct {
		return fmt.Errorf("UnmarshalRow only accepts a pointer to a struct.")
	}
	outv := reflect.Indirect(reflect.ValueOf(out))
	outt := outv.Type()
	for i := 0; i < outt.NumField(); i++ {
		fieldValue := outv.Field(i)
		if fieldValue.CanSet() { // Only consider exported fields.
			field := outt.Field(i)
			fieldName := strings.ToLower(field.Name)
			f, fieldIndex := s.GetField(fieldName)
			if fieldIndex != InvalidPosition {
				cell := row[fieldIndex]
				if s.isMissingValue(cell) {
					continue
				}
				v, err := f.UnmarshalString(cell)
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

func (s *Schema) isMissingValue(value string) bool {
	for _, mv := range s.MissingValues {
		if mv == value {
			return true
		}
	}
	return false
}

// UnmarshalJSON sets *f to a copy of data. It will respect the default values
// described at: https://specs.frictionlessdata.io/table-schema/
func (s *Schema) UnmarshalJSON(data []byte) error {
	// This is neded so it does not call UnmarshalJSON from recursively.
	type schemaAlias Schema
	var a schemaAlias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	if err := processPlaceholder(a.PrimaryKeyPlaceholder, &a.PrimaryKeys); err != nil {
		return fmt.Errorf("primaryKey must be either a string or list")
	}
	a.PrimaryKeyPlaceholder = nil
	if err := processPlaceholder(a.ForeignKeys.FieldsPlaceholder, &a.ForeignKeys.Fields); err != nil {
		return fmt.Errorf("foreignKeys.fields must be either a string or list")
	}
	a.ForeignKeys.FieldsPlaceholder = nil
	if err := processPlaceholder(a.ForeignKeys.Reference.FieldsPlaceholder, &a.ForeignKeys.Reference.Fields); err != nil {
		return fmt.Errorf("foreignKeys.reference.fields must be either a string or list")
	}
	a.ForeignKeys.Reference.FieldsPlaceholder = nil
	*s = Schema(a)
	return nil
}

// MarshalJSON returns the JSON encoding of s.
func (s *Schema) MarshalJSON() ([]byte, error) {
	type schemaAlias Schema
	a := schemaAlias(*s)
	a.PrimaryKeyPlaceholder = a.PrimaryKeys
	a.ForeignKeys.Reference.FieldsPlaceholder = a.ForeignKeys.Reference.Fields
	return json.Marshal(a)
}

func processPlaceholder(ph interface{}, v *[]string) error {
	if ph == nil {
		return nil
	}
	if vStr, ok := ph.(string); ok {
		*v = append(*v, vStr)
		return nil
	}
	if vSlice, ok := ph.([]interface{}); ok {
		for i := range vSlice {
			*v = append(*v, vSlice[i].(string))
		}
		return nil
	}
	// Only for signalling that an error happened. The caller knows the best
	// error message.
	return fmt.Errorf("")
}

// DecodeTable loads and decodes all table rows.
//
// The result argument must necessarily be the address for a slice. The slice
// may be nil or previously allocated.
func (s *Schema) DecodeTable(tab table.Table, out interface{}) error {
	outv := reflect.ValueOf(out)
	if outv.Kind() != reflect.Ptr || outv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("out argument must be a slice address")
	}
	iter, err := tab.Iter()
	if err != nil {
		return err
	}
	defer iter.Close()
	slicev := outv.Elem()
	slicev = slicev.Slice(0, 0) // Trucantes the passed-in slice.
	elemt := slicev.Type().Elem()
	i := 0
	for iter.Next() {
		elemp := reflect.New(elemt)
		if err := s.Decode(iter.Row(), elemp.Interface()); err != nil {
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
