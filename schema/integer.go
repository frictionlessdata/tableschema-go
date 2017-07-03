package schema

import "strconv"

// CastInt casts an integer value (passed-in as unicode string) against a field. Returns an
// error if the value can not be converted to integer.
func castInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
