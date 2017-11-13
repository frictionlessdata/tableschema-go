package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/frictionlessdata/tableschema-go/table"
)

// InvalidPosition is returned by GetField call when
// it refers to a field that does not exist in the schema.
const InvalidPosition = -1

// Unexported tagname for the tableheader
const tableheaderTag = "tableheader"

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
	if len(s.MissingValues) == 0 {
		return &s, nil
	}
	// Transforming the list in a set.
	valueSet := make(map[string]struct{}, len(s.MissingValues))
	for _, v := range s.MissingValues {
		valueSet[v] = struct{}{}
	}
	// Updating fields.
	for i := range s.Fields {
		s.Fields[i].MissingValues = make(map[string]struct{}, len(valueSet))
		for k, v := range valueSet {
			s.Fields[i].MissingValues[k] = v
		}
	}
	return &s, nil
}

// LoadFromFile loads and parses a schema descriptor from a local file.
func LoadFromFile(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return Read(f)
}

var (
	httpClient *http.Client
	once       sync.Once
)

const remoteFetchTimeoutSecs = 15

// LoadRemote downloads and parses a schema descriptor from the specified URL.
func LoadRemote(url string) (*Schema, error) {
	once.Do(func() {
		httpClient = &http.Client{
			Timeout: remoteFetchTimeoutSecs * time.Second,
		}
	})
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return Read(resp.Body)
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

// CastRow casts the passed-in row to schema types and stores it in the value pointed
// by out. The out value must be pointer to a struct. Only exported fields will be unmarshalled.
// The lowercased field name is used as the key for each exported field.
//
// If a value in the row cannot be marshalled to its respective schema field (Field.Unmarshal),
// this call will return an error. Furthermore, this call is also going to return an error if
// the schema field value can not be unmarshalled to the struct field type.
func (s *Schema) CastRow(row []string, out interface{}) error {
	if reflect.ValueOf(out).Kind() != reflect.Ptr || reflect.Indirect(reflect.ValueOf(out)).Kind() != reflect.Struct {
		return fmt.Errorf("can only cast pointer to structs")
	}
	outv := reflect.Indirect(reflect.ValueOf(out))
	outt := outv.Type()
	for i := 0; i < outt.NumField(); i++ {
		fieldValue := outv.Field(i)
		if fieldValue.CanSet() { // Only consider exported fields.
			field := outt.Field(i)
			fieldName, ok := field.Tag.Lookup(tableheaderTag)
			if !ok { // if no tag is set use own name
				fieldName = field.Name
			}
			f, fieldIndex := s.GetField(fieldName)
			if fieldIndex != InvalidPosition {
				cell := row[fieldIndex]
				if s.isMissingValue(cell) {
					continue
				}
				v, err := f.Cast(cell)
				if err != nil {
					return err
				}
				toSetValue := reflect.ValueOf(v)
				toSetType := toSetValue.Type()
				if !toSetType.ConvertibleTo(field.Type) {
					return fmt.Errorf("value:%s field:%s - can not convert from %v to %v", field.Name, cell, toSetType, field.Type)
				}
				fieldValue.Set(toSetValue.Convert(field.Type))
			}
		}
	}
	return nil
}

type rawCell struct {
	pos int
	val string
}

type rawRow []rawCell

func (r rawRow) Len() int           { return len(r) }
func (r rawRow) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r rawRow) Less(i, j int) bool { return r[i].pos < r[j].pos }

// UncastRow uncasts struct into a row. This method can only uncast structs (or pointer to structs) and
// will error out if nil is passed.
// The order of the cells in the returned row is the schema declaration order.
func (s *Schema) UncastRow(in interface{}) ([]string, error) {
	inValue := reflect.Indirect(reflect.ValueOf(in))
	if inValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can only uncast structs and does not support nil pointers")
	}
	inType := inValue.Type()
	var row rawRow
	for i := 0; i < inType.NumField(); i++ {
		structFieldValue := inValue.Field(i)
		fieldName, ok := inType.Field(i).Tag.Lookup(tableheaderTag)
		if !ok {
			fieldName = inType.Field(i).Name
		}
		f, fieldIndex := s.GetField(fieldName)
		if fieldIndex != InvalidPosition {
			cell, err := f.Uncast(structFieldValue.Interface())
			if err != nil {
				return nil, err
			}
			row = append(row, rawCell{fieldIndex, cell})
		}
	}
	sort.Sort(row)
	ret := make([]string, len(row))
	for i := range row {
		ret[i] = row[i].val
	}
	return ret, nil
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

// uniqueKey represents field ID and field value which then can be used for equality tests (e.g. in a map key)
type uniqueKey struct {
	KeyIndex int
	KeyValue interface{}
}

// CastTable loads and casts all table rows.
//
// The result argument must necessarily be the address for a slice. The slice
// may be nil or previously allocated.
func (s *Schema) CastTable(tab table.Table, out interface{}) error {
	outv := reflect.ValueOf(out)
	if outv.Kind() != reflect.Ptr || outv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("out argument must be a slice address")
	}
	iter, err := tab.Iter()
	if err != nil {
		return err
	}
	defer iter.Close()

	uniqueFieldIndexes := extractUniqueFieldIndexes(s)
	uniqueCache := make(map[uniqueKey]struct{})

	slicev := outv.Elem()
	slicev = slicev.Slice(0, 0) // Trucantes the passed-in slice.
	elemt := slicev.Type().Elem()
	i := 0
	for iter.Next() {
		i++
		elemp := reflect.New(elemt)
		if err := s.CastRow(iter.Row(), elemp.Interface()); err != nil {
			return err
		}
		for _, k := range uniqueFieldIndexes {
			field := elemp.Elem().Field(k)
			if _, ok := uniqueCache[uniqueKey{k, field.Interface()}]; ok {
				return fmt.Errorf("field(s) '%s' duplicates in row %v", elemp.Elem().Type().Field(k).Name, i)
			}
			uniqueCache[uniqueKey{k, field.Interface()}] = struct{}{}
		}
		slicev = reflect.Append(slicev, elemp.Elem())
		slicev = slicev.Slice(0, slicev.Len())
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	outv.Elem().Set(slicev.Slice(0, i))
	return nil
}

func extractUniqueFieldIndexes(s *Schema) []int {
	uniqueIndexes := make(map[int]struct{})
	for _, pk := range s.PrimaryKeys {
		_, index := s.GetField(pk)
		uniqueIndexes[index] = struct{}{}
	}
	for i := range s.Fields {
		if _, ok := uniqueIndexes[i]; !ok && s.Fields[i].Constraints.Unique {
			uniqueIndexes[i] = struct{}{}
		}
	}
	keys := make([]int, 0, len(uniqueIndexes))
	for k := range uniqueIndexes {
		keys = append(keys, k)
	}
	return keys
}

// UncastTable uncasts each element (struct) of the passed-in slice and
func (s *Schema) UncastTable(in interface{}) ([][]string, error) {
	inVal := reflect.Indirect(reflect.ValueOf(in))
	if inVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("tables must be slice of structs")
	}
	var t [][]string
	for i := 0; i < inVal.Len(); i++ {
		r, err := s.UncastRow(inVal.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		t = append(t, r)
	}
	return t, nil
}

// String returns an human readable version of the schema.
func (s *Schema) String() string {
	var buf bytes.Buffer
	pp, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	buf.Write(pp)
	return buf.String()
}
