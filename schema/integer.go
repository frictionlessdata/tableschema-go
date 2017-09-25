package schema

import (
	"fmt"
	"regexp"
	"strconv"
)

// CastInt casts an integer value (passed-in as unicode string) against a field. Returns an
// error if the value can not be converted to integer.
func castInt(bareNumber bool, value string, c Constraints) (int64, error) {
	v := value
	if !bareNumber {
		var err error
		v, err = stripIntegerFromString(v)
		if err != nil {
			return 0, err
		}
	}
	returned, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	if c.Maximum != "" {
		max, err := strconv.ParseInt(c.Maximum, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid maximum integer: %v", c.Maximum)
		}
		if returned > max {
			return 0, fmt.Errorf("constraint check error: integer:%d > maximum:%d", returned, max)
		}
	}
	if c.Minimum != "" {
		min, err := strconv.ParseInt(c.Minimum, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid minimum integer: %v", c.Minimum)
		}
		if returned < min {
			return 0, fmt.Errorf("constraint check error: integer:%d < minimum:%d", returned, min)
		}
	}
	return returned, nil
}

var bareIntegerRegexp = regexp.MustCompile(`((^[0-9]+)|([0-9]+$))`)

func stripIntegerFromString(v string) (string, error) {
	matches := bareIntegerRegexp.FindStringSubmatch(v)
	if matches == nil {
		return "", fmt.Errorf("invalid integer to strip:%s", v)
	}
	return matches[1], nil
}
