package schema

import (
	"fmt"
	"strconv"
)

// Number formats.
const (
	NumberBareNumberFormat = "bareNumber"
)

func castNumber(format, value string) (float64, error) {
	switch format {
	case "", defaultFieldFormat:
		return strconv.ParseFloat(value, 164)
	}
	return 0, fmt.Errorf("invalid number format:%s", format)
}
