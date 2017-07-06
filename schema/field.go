package schema

import (
	"encoding/json"
	"fmt"
)

// Default for schema fields.
const (
	defaultFieldType   = "string"
	defaultFieldFormat = "default"
)

// Default schema variables.
var (
	defaultTrueValues  = []string{"yes", "y", "true", "t", "1"}
	defaultFalseValues = []string{"no", "n", "false", "f", "0"}
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

	// Boolean properties. Define set of the values that represent true and false, respectively.
	// https://specs.frictionlessdata.io/table-schema/#boolean
	TrueValues  []string `json:"trueValues"`
	FalseValues []string `json:"falseValues"`
}

// UnmarshalJSON sets *f to a copy of data. It will respect the default values
// described at: https://specs.frictionlessdata.io/table-schema/
func (f *Field) UnmarshalJSON(data []byte) error {
	// This is neded so it does not call UnmarshalJSON from recursively.
	type fieldAlias Field
	u := &fieldAlias{
		Type:        defaultFieldType,
		Format:      defaultFieldFormat,
		TrueValues:  defaultTrueValues,
		FalseValues: defaultFalseValues,
	}
	if err := json.Unmarshal(data, u); err != nil {
		return err
	}
	*f = Field(*u)
	return nil
}

// CastValue casts a value against field. Returns an error if the value can
// not be cast or any field constraint can no be satisfied.
func (f *Field) CastValue(value string) (interface{}, error) {
	switch f.Type {
	case IntegerType:
		return castInt(value)
	case StringType:
		return castString(f.Type, value)
	}
	return nil, fmt.Errorf("invalid field type: %s", f.Type)
}

// TestValue checks whether the value can be casted against the field.
func (f *Field) TestValue(value string) bool {
	_, err := f.CastValue(value)
	return err == nil
}
