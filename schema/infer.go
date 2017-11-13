package schema

import (
	"fmt"

	"github.com/frictionlessdata/tableschema-go/table"
)

var (
	// https://specs.frictionlessdata.io/table-schema/#boolean
	booleanValues = map[string]struct{}{
		"true":  struct{}{},
		"True":  struct{}{},
		"TRUE":  struct{}{},
		"1":     struct{}{},
		"false": struct{}{},
		"False": struct{}{},
		"FALSE": struct{}{},
		"0":     struct{}{},
	}
	// This structure is optmized for querying.
	// It should point a type to what is allowed to be implicitly cast.
	// The inner set must be sorted by the narrower first.
	implicitCast = map[string][]string{
		IntegerType:   []string{IntegerType, NumberType, StringType},
		NumberType:    []string{NumberType, StringType},
		BooleanType:   []string{BooleanType, IntegerType, NumberType, StringType},
		YearMonthType: []string{YearMonthType, DateType, StringType},
		YearType:      []string{YearType, IntegerType, NumberType, StringType},
		DateType:      []string{DateType, DateTimeType, StringType},
		DateTimeType:  []string{DateTimeType, StringType},
		TimeType:      []string{TimeType, StringType},
		DurationType:  []string{DurationType, StringType},
		ObjectType:    []string{ObjectType, StringType},
		ArrayType:     []string{ArrayType, StringType},
		GeoPointType:  []string{GeoPointType, ArrayType, StringType},
		StringType:    []string{},
	}

	// Types ordered from narrower to wider.
	orderedTypes = []string{BooleanType, YearType, IntegerType, GeoPointType, NumberType, YearMonthType, DateType, DateTimeType, TimeType, DurationType, ArrayType, ObjectType}

	noConstraints = Constraints{}
)

// Maximum number of rows used to infer schema.
const maxNumRowsInfer = 100

// Infer infers a schema from a slice of the tabular data. For columns that contain
// cells that can inferred as different types, the most popular type is set as the field
// type. For instance, a column with values 10.1, 10, 10 will inferred as being of type
// "integer".
func Infer(tab table.Table) (*Schema, error) {
	s, err := sample(tab)
	if err != nil {
		return nil, err
	}
	return infer(tab.Headers(), s)
}

func sample(tab table.Table) ([][]string, error) {
	iter, err := tab.Iter()
	if err != nil {
		return nil, err
	}
	var t [][]string
	for count := 0; count < maxNumRowsInfer && iter.Next(); count++ {
		t = append(t, iter.Row())
	}
	if iter.Err() != nil {
		return nil, iter.Err()
	}
	return t, nil
}

func infer(headers []string, table [][]string) (*Schema, error) {
	inferredTypes := make([]map[string]int, len(headers))
	for rowID := range table {
		row := table[rowID]
		// TODO(danielfireman): the python version does some normalization on
		// the number of columns and headers. Need to look closer at this.
		if len(headers) != len(row) {
			return nil, fmt.Errorf("data is not tabular. headers:%v row[%d]:%v", headers, rowID, row)
		}
		for cellIndex, cell := range row {
			if inferredTypes[cellIndex] == nil {
				inferredTypes[cellIndex] = make(map[string]int)
			}
			// The list bellow must be ordered by the narrower field type.
			t := findType(cell, orderedTypes)
			inferredTypes[cellIndex][t]++
		}
	}
	schema := Schema{}
	for index := range headers {
		schema.Fields = append(schema.Fields,
			Field{
				Name:   headers[index],
				Type:   defaultFieldType,
				Format: defaultFieldFormat,
			})
		count := 0
		for t, c := range inferredTypes[index] {
			if c > count {
				f := &schema.Fields[index]
				f.Type = t
				count = c
			}
		}
	}
	return &schema, nil
}

// InferImplicitCasting uses a implicit casting for infering the type of columns
// that have cells of diference types. For instance, a column with values 10.1, 10, 10
// will inferred as being of type "number" ("integer" can be implicitly cast to "number").
//
// For medium to big tables, this method is faster than the Infer.
func InferImplicitCasting(tab table.Table) (*Schema, error) {
	s, err := sample(tab)
	if err != nil {
		return nil, err
	}
	return inferImplicitCasting(tab.Headers(), s)
}

func inferImplicitCasting(headers []string, table [][]string) (*Schema, error) {
	inferredTypes := make([]string, len(headers))
	for rowID := range table {
		row := table[rowID]
		// TODO(danielfireman): the python version does some normalization on
		// the number of columns and headers. Need to look closer at this.
		if len(headers) != len(row) {
			return nil, fmt.Errorf("data is not tabular. headers:%v row[%d]:%v", headers, rowID, row)
		}
		for cellIndex, cell := range row {
			if inferredTypes[cellIndex] == "" {
				t := findType(cell, orderedTypes)
				inferredTypes[cellIndex] = t
			} else {
				inferredTypes[cellIndex] = findType(cell, implicitCast[inferredTypes[cellIndex]])
			}
		}
	}
	schema := Schema{}
	for index := range headers {
		schema.Fields = append(schema.Fields,
			Field{
				Name:   headers[index],
				Type:   inferredTypes[index],
				Format: defaultFieldFormat,
			})
	}
	return &schema, nil
}

func findType(value string, checkOrder []string) string {
	for _, t := range checkOrder {
		switch t {
		case BooleanType:
			if _, ok := booleanValues[value]; ok {
				return BooleanType
			}
		case IntegerType:
			if _, err := castInt(defaultBareNumber, value, noConstraints); err == nil {
				return IntegerType
			}
		case NumberType:
			if _, err := castNumber(defaultDecimalChar, defaultGroupChar, defaultBareNumber, value, noConstraints); err == nil {
				return NumberType
			}
		case DateType:
			if _, err := castDate(defaultFieldFormat, value, noConstraints); err == nil {
				return DateType
			}
		case ArrayType:
			if _, err := castArray(value); err == nil {
				return ArrayType
			}
		case ObjectType:
			if _, err := castObject(value); err == nil {
				return ObjectType
			}
		case TimeType:
			if _, err := castTime(defaultFieldFormat, value, noConstraints); err == nil {
				return TimeType
			}
		case YearMonthType:
			if _, err := castYearMonth(value, noConstraints); err == nil {
				return YearMonthType
			}
		case YearType:
			if _, err := castYear(value, noConstraints); err == nil {
				return YearType
			}
		case DateTimeType:
			if _, err := castDateTime(value, noConstraints); err == nil {
				return DateTimeType
			}
		case DurationType:
			if _, err := castDuration(value); err == nil {
				return DurationType
			}
		case GeoPointType:
			if _, err := castGeoPoint(defaultFieldFormat, value); err == nil {
				return GeoPointType
			}
		}
	}
	return StringType
}
