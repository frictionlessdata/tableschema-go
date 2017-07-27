package schema

import "fmt"

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
	// This structure is optmized for querying. It a set inside a map.
	// It should point a type to what is allowed to be implicitly cast.
	// The inner set must be sorted by the narrower first.
	implicitCast = map[string][]string{
		IntegerType: []string{IntegerType, NumberType, StringType},
		NumberType:  []string{NumberType, StringType},
		BooleanType: []string{BooleanType, IntegerType, NumberType, StringType},
		DateType:    []string{DateType, DateTimeType, StringType},
		ObjectType:  []string{ObjectType, StringType},
		ArrayType:   []string{ArrayType, StringType},
		StringType:  []string{},
	}
)

// Infer infers a schema from tabular data.
func Infer(headers []string, table [][]string) (*Schema, error) {
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
			t := findType(cell, []string{BooleanType, IntegerType, NumberType, DateType, ArrayType, ObjectType})
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
// that have cells of diference types. For instance, if a certain column has a cell
// with value 1 and another cell with value 1.2, the type inferred type of this column
// would end up with the type "number" ("integer" can be implicitly cast to "number").
//
// For medium to big tables, this method is faster than the Infer.
func InferImplicitCasting(headers []string, table [][]string) (*Schema, error) {
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
				// The list bellow must be ordered by the narrower field type.
				t := findType(cell, []string{BooleanType, IntegerType, NumberType, DateType, ArrayType, ObjectType})
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
			if _, err := castInt(value); err == nil {
				return IntegerType
			}
		case NumberType:
			if _, err := castNumber(value); err == nil {
				return NumberType
			}
		case DateType:
			if _, err := castDate(defaultFieldFormat, value); err == nil {
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
		}
	}
	return StringType
}
