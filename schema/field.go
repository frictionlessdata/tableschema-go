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

	// Unique indicates whether this field is allowed to have duplicates.
	// This constrain is only relevant for Schema.CastTable
	Unique bool `json:"unique,omitempty"`

	Maximum         string `json:"maximum,omitempty"`
	Minimum         string `json:"minimum,omitempty"`
	MinLength       int    `json:"minLength,omitempty"`
	MaxLength       int    `json:"maxLength,omitempty"`
	Pattern         string `json:"pattern,omitempty"`
	compiledPattern *regexp.Regexp

	// Enum indicates that the value of the field must exactly match a value in the enum array.
	// The values of the fields could need encoding, depending on the type.
	// It applies to all field types.
	Enum []interface{} `json:"enum,omitempty"`
	// rawEnum keeps the raw version of the enum objects, to make validation faster and easier.
	rawEnum map[string]struct{}
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
	// Transformation/Validation that should be done at creation time.
	if f.Constraints.Pattern != "" {
		p, err := regexp.Compile(f.Constraints.Pattern)
		if err != nil {
			return err
		}
		f.Constraints.compiledPattern = p
	}
	if f.Constraints.Enum != nil {
		f.Constraints.rawEnum = make(map[string]struct{})
		for i := range f.Constraints.Enum {
			e, err := f.Uncast(f.Constraints.Enum[i])
			if err != nil {
				return err
			}
			f.Constraints.rawEnum[e] = struct{}{}
		}
	}
	return nil
}

// Cast casts the passed-in string against field type. Returns an error
// if the value can not be cast or any field constraint can not be satisfied.
func (f *Field) Cast(value string) (interface{}, error) {
	if f.Constraints.Required {
		_, ok := f.MissingValues[value]
		if ok {
			return nil, fmt.Errorf("%s is required", f.Name)
		}
	}
	var castd interface{}
	var err error
	switch f.Type {
	case IntegerType:
		castd, err = castInt(f.BareNumber, value, f.Constraints)
	case StringType:
		castd, err = castString(f.Format, value, f.Constraints)
	case BooleanType:
		castd, err = castBoolean(value, f.TrueValues, f.FalseValues)
	case NumberType:
		castd, err = castNumber(f.DecimalChar, f.GroupChar, f.BareNumber, value, f.Constraints)
	case DateType:
		castd, err = castDate(f.Format, value, f.Constraints)
	case ObjectType:
		castd, err = castObject(value)
	case ArrayType:
		castd, err = castArray(value)
	case TimeType:
		castd, err = castTime(f.Format, value, f.Constraints)
	case YearMonthType:
		castd, err = castYearMonth(value, f.Constraints)
	case YearType:
		castd, err = castYear(value, f.Constraints)
	case DateTimeType:
		castd, err = castDateTime(value, f.Constraints)
	case DurationType:
		castd, err = castDuration(value)
	case GeoPointType:
		castd, err = castGeoPoint(f.Format, value)
	case AnyType:
		castd, err = castAny(value)
	}
	if err != nil {
		return nil, err
	}
	if castd == nil {
		return nil, fmt.Errorf("invalid field type: %s", f.Type)
	}
	if len(f.Constraints.rawEnum) > 0 {
		rawValue, err := f.Uncast(castd)
		if err != nil {
			return nil, err
		}
		if _, ok := f.Constraints.rawEnum[rawValue]; !ok {
			return nil, fmt.Errorf("castd value:%s does not match enum constraints:%v", rawValue, f.Constraints.rawEnum)
		}
	}
	return castd, nil
}

// Uncast uncasts the passed-in value into a string. It returns an error if the
// the type of the passed-in value can not be converted to field type.
func (f *Field) Uncast(in interface{}) (string, error) {
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
		return uncastBoolean(in, f.TrueValues, f.FalseValues)
	case DurationType:
		return uncastDuration(inInterface)
	case GeoPointType:
		return uncastGeoPoint(f.Format, in)
	case DateType, DateTimeType, TimeType, YearMonthType, YearType:
		return uncastTime(inInterface)
	case ObjectType:
		return uncastObject(inInterface)
	case StringType:
		_, ok = inInterface.(string)
	case ArrayType:
		ok = reflect.TypeOf(inInterface).Kind() == reflect.Slice
	case AnyType:
		return uncastAny(in)
	}
	if !ok {
		return "", fmt.Errorf("can not convert \"%d\" which type is %s to type %s", in, reflect.TypeOf(in), f.Type)
	}
	return fmt.Sprintf("%v", inInterface), nil
}

// TestString checks whether the value can be unmarshalled to the field type.
func (f *Field) TestString(value string) bool {
	_, err := f.Cast(value)
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
