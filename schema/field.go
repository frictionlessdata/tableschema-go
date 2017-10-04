package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
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
	defaultDecimalChar = "."
	defaultGroupChar   = ","
	defaultBareNumber  = true
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
	AnyType       = "any"
)

// Formats.
const (
	AnyDateFormat = "any"
)

// Constraints can be used by consumers to list constraints for validating
// field values.
type Constraints struct {
	// Required indicates whether this field is allowed to be null.
	// Schema.MissingValues define how the string representation can
	// represent null values.
	Required bool `json:"required,omitempty"`

	Maximum         string `json:"maximum,omitempty"`
	Minimum         string `json:"minimum,omitempty"`
	MinLength       int    `json:"minLength,omitempty"`
	MaxLength       int    `json:"maxLength,omitempty"`
	Pattern         string `json:"pattern,omitempty"`
	compiledPattern *regexp.Regexp
}

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

	// Number/Integer properties.

	// A string whose value is used to represent a decimal point within the number. The default value is ".".
	DecimalChar string `json:"decimalChar,omitempty"`
	// A string whose value is used to group digits within the number. The default value is null. A common value is "," e.g. "100,000".
	GroupChar string `json:"groupChar,omitempty"`
	// If true the physical contents of this field must follow the formatting constraints already set out.
	// If false the contents of this field may contain leading and/or trailing non-numeric characters which
	// are going to be stripped. Default value is true:
	BareNumber bool `json:"bareNumber,omitempty"`

	// MissingValues is a map which dictates which string values should be treated as null
	// values.
	MissingValues map[string]struct{} `json:"-"`

	// Constraints can be used by consumers to list constraints for validating
	// field values.
	Constraints Constraints
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
		DecimalChar: defaultDecimalChar,
		GroupChar:   defaultGroupChar,
		BareNumber:  defaultBareNumber,
	}
	if err := json.Unmarshal(data, u); err != nil {
		return err
	}
	*f = Field(*u)

	if f.Constraints.Pattern != "" {
		p, err := regexp.Compile(f.Constraints.Pattern)
		if err != nil {
			return err
		}
		f.Constraints.compiledPattern = p
	}
	return nil
}

// Decode decodes the passed-in string against field type. Returns an error
// if the value can not be cast or any field constraint can not be satisfied.
func (f *Field) Decode(value string) (interface{}, error) {
	if f.Constraints.Required {
		_, ok := f.MissingValues[value]
		if ok {
			return nil, fmt.Errorf("%s is required", f.Name)
		}
	}
	switch f.Type {
	case IntegerType:
		return castInt(f.BareNumber, value, f.Constraints)
	case StringType:
		return decodeString(f.Format, value, f.Constraints)
	case BooleanType:
		return castBoolean(value, f.TrueValues, f.FalseValues)
	case NumberType:
		return castNumber(f.DecimalChar, f.GroupChar, f.BareNumber, value, f.Constraints)
	case DateType:
		return decodeDate(f.Format, value, f.Constraints)
	case ObjectType:
		return castObject(value)
	case ArrayType:
		return castArray(value)
	case TimeType:
		return decodeTime(f.Format, value, f.Constraints)
	case YearMonthType:
		return decodeYearMonth(value, f.Constraints)
	case YearType:
		return decodeYear(value, f.Constraints)
	case DateTimeType:
		return decodeDateTime(value, f.Constraints)
	case DurationType:
		return castDuration(value)
	case GeoPointType:
		return castGeoPoint(f.Format, value)
	case AnyType:
		return castAny(value)
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
	case AnyType:
		return encodeAny(in)
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
