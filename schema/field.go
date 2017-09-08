package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
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

// Field types.
const (
	IntegerType   = "integer"
	StringType    = "string"
	BooleanType   = "boolean"
	NumberType    = "number"
	DateType      = "date"
	ObjectType    = "object"
	ArrayType     = "array"
	DateTimeType  = "datetime"
	TimeType      = "time"
	YearMonthType = "yearmonth"
	YearType      = "year"
	DurationType  = "duration"
	GeoPointType  = "geopoint"
)

// Formats.
const (
	AnyDateFormat = "any"
)

// Field describes a single field in the table schema.
// More: https://specs.frictionlessdata.io/table-schema/#field-descriptors
type Field struct {
	// Name of the field. It is mandatory and shuold correspond to the name of field/column in the data file (if it has a name).
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
	// A human readable label or title for the field.
	Title string `json:"title,omitempty"`
	// A description for this field e.g. "The recipient of the funds"
	Description string `json:"description,omitempty"`

	// Boolean properties. Define set of the values that represent true and false, respectively.
	// https://specs.frictionlessdata.io/table-schema/#boolean
	TrueValues  []string `json:"trueValues,omitempty"`
	FalseValues []string `json:"falseValues,omitempty"`
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

// Decode decodes the passed-in string against field type. Returns an error
// if the value can not be cast or any field constraint can not be satisfied.
func (f *Field) Decode(value string) (interface{}, error) {
	switch f.Type {
	case IntegerType:
		return castInt(value)
	case StringType:
		return castString(f.Format, value)
	case BooleanType:
		return castBoolean(value, f.TrueValues, f.FalseValues)
	case NumberType:
		return castNumber(value)
	case DateType:
		return castDate(f.Format, value)
	case ObjectType:
		return castObject(value)
	case ArrayType:
		return castArray(value)
	case TimeType:
		return castTime(f.Format, value)
	case YearMonthType:
		return castYearMonth(value)
	case YearType:
		return castYear(value)
	case DateTimeType:
		return castDateTime(f.Format, value)
	case DurationType:
		return castDuration(value)
	case GeoPointType:
		return castGeoPoint(f.Format, value)
	}
	return nil, fmt.Errorf("invalid field type: %s", f.Type)
}

// Encode encodes the passed-in value into a string. It returns an error if the
// the type of the passed-in value can not be converted to field type.
func (f *Field) Encode(in interface{}) (string, error) {
	// This indirect avoids the need to custom-case pointer types.
	inValue := reflect.Indirect(reflect.ValueOf(in))
	inInterface := inValue.Interface()
	ok := false
	switch f.Type {
	case IntegerType:
		var a int64
		ok = reflect.TypeOf(inInterface).ConvertibleTo(reflect.ValueOf(a).Type())
		if ok {
			inInterface = inValue.Convert(reflect.ValueOf(a).Type()).Interface()
		}
	case NumberType:
		var a float64
		ok = reflect.TypeOf(inInterface).ConvertibleTo(reflect.ValueOf(a).Type())
		if ok {
			inInterface = inValue.Convert(reflect.ValueOf(a).Type()).Interface()
		}
	case BooleanType:
		return encodeBoolean(in, f.TrueValues, f.FalseValues)
	case DurationType:
		return encodeDuration(inInterface)
	case GeoPointType:
		return encodeGeoPoint(f.Format, in)
	case DateType, DateTimeType, TimeType, YearMonthType, YearType:
		return encodeTime(inInterface)
	case ObjectType:
		return encodeObject(inInterface)
	case StringType:
		_, ok = inInterface.(string)
	case ArrayType:
		ok = reflect.TypeOf(inInterface).Kind() == reflect.Slice
	}
	if !ok {
		return "", fmt.Errorf("can not convert \"%d\" which type is %s to type %s", in, reflect.TypeOf(in), f.Type)
	}
	return fmt.Sprintf("%v", inInterface), nil
}

// TestString checks whether the value can be unmarshalled to the field type.
func (f *Field) TestString(value string) bool {
	_, err := f.Decode(value)
	return err == nil
}

// asReadField returns the field passed-in as parameter like it's been read as JSON.
// That include setting default values.
// Created for being used in tests.
// IMPORTANT: Not ready for being used in production due to possibly bad performance.
func asJSONField(f Field) Field {
	var out Field
	data, _ := json.Marshal(&f)
	json.Unmarshal(data, &out)
	return out
}
