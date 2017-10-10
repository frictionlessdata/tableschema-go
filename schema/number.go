package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func castNumber(decimalChar, groupChar string, bareNumber bool, value string, c Constraints) (float64, error) {
	dc := decimalChar
	if groupChar != "" {
		dc = decimalChar
	}
	v := strings.Replace(value, dc, ".", 1)
	gc := defaultGroupChar
	if groupChar != "" {
		gc = groupChar
	}
	v = strings.Replace(v, gc, "", -1)
	if !bareNumber {
		var err error
		v, err = stripNumberFromString(v)
		if err != nil {
			return 0, err
		}
	}
	returned, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}
	if c.Maximum != "" {
		max, err := strconv.ParseFloat(c.Maximum, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid maximum number: %v", c.Maximum)
		}
		if returned > max {
			return 0, fmt.Errorf("constraint check error: integer:%f > maximum:%f", returned, max)
		}
	}
	if c.Minimum != "" {
		min, err := strconv.ParseFloat(c.Minimum, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid minimum integer: %v", c.Minimum)
		}
		if returned < min {
			return 0, fmt.Errorf("constraint check error: integer:%f < minimum:%f", returned, min)
		}
	}
	return returned, nil
}

var bareNumberRegexp = regexp.MustCompile(`((^[0-9]+\.?[0-9]*)|([0-9]+\.?[0-9]*$))`)

func stripNumberFromString(v string) (string, error) {
	matches := bareNumberRegexp.FindStringSubmatch(v)
	if matches == nil {
		return "", fmt.Errorf("invalid number to strip:%s", v)
	}
	return matches[1], nil
}
