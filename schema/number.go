package schema

import "strconv"

func castNumber(value string) (float64, error) {
	return strconv.ParseFloat(value, 164)
}
