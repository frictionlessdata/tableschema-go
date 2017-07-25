package schema

import "fmt"

type inferredType struct {
	t string
	f string
}

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
)

// Infer infers a schema from tabular data.
func Infer(headers []string, table [][]string) (*Schema, error) {
	inferredTypes := make(map[int]map[string]int, len(headers))
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
			// The list bellow must ordered by the narrower fieldt type.
			for _, t := range []string{BooleanType, IntegerType, NumberType, DateType} {
				found := false
				switch t {
				case BooleanType:
					// https://specs.frictionlessdata.io/table-schema/#boolean
					if _, ok := booleanValues[cell]; ok {
						inferredTypes[cellIndex][BooleanType]++
						found = true
					}
				case IntegerType:
					if _, err := castInt(cell); err == nil {
						inferredTypes[cellIndex][IntegerType]++
						found = true
					}
				case NumberType:
					if _, err := castNumber(cell); err == nil {
						inferredTypes[cellIndex][NumberType]++
						found = true
					}
				case DateType:
					if _, err := castDate(defaultFieldFormat, cell); err == nil {
						inferredTypes[cellIndex][DateType]++
						found = true
					}
				}
				if found {
					break
				}
			}
		}
	}
	schema := Schema{}
	for index := range headers {
		schema.Fields = append(schema.Fields, Field{Name: headers[index], Type: defaultFieldType, Format: defaultFieldFormat})
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
