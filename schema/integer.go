package schema

import (
	"fmt"
	"regexp"
	"strconv"
)

// CastInt casts an integer value (passed-in as unicode string) against a field. Returns an
// error if the value can not be converted to integer.
func castInt(bareNumber bool, value string) (int64, error) {
	v := value
	if !bareNumber {
		var err error
		v, err = stripIntegerFromString(v)
		if err != nil {
			return 0, err
		}
	}
	return strconv.ParseInt(v, 10, 64)
}

var bareIntegerRegexp = regexp.MustCompile(`((^[0-9]+)|([0-9]+$))`)

func stripIntegerFromString(v string) (string, error) {
	matches := bareIntegerRegexp.FindStringSubmatch(v)
	if matches == nil {
		return "", fmt.Errorf("invalid integer to strip:%s", v)
	}
	return matches[1], nil
}
