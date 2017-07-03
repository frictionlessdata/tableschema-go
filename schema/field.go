package schema

import "fmt"

// Default for schema fields.
const (
	defaultFieldType   = "string"
	defaultFieldFormat = "default"
)

// Field Types.
const (
	IntegerType = "integer"
	StringType  = "string"
)

// Field represents a cell on a table.
type Field struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format"`
}

// CastValue casts a value against field. Returns an error if the value can
// not be cast or any field constraint can no be satisfied.
func (f Field) CastValue(value string) (interface{}, error) {
	switch f.Type {
	case IntegerType:
		return castInt(value)
	case StringType:
		return castString(f.Type, value)
	}
	return nil, fmt.Errorf("invalid field type: %s", f.Type)
}

// TestValue checks whether the value can be casted against the field.
func (f Field) TestValue(value string) bool {
	if _, err := f.CastValue(value); err != nil {
		return false
	}
	return true
}

func setDefaultValues(f *Field) {
	if f.Type == "" {
		f.Type = defaultFieldType
	}
	if f.Format == "" {
		f.Format = defaultFieldFormat
	}
}
