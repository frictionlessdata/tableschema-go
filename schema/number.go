package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func castNumber(decimalChar, groupChar string, bareNumber bool, value string) (float64, error) {
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
	return strconv.ParseFloat(v, 64)
}

var bareNumberRegexp = regexp.MustCompile(`((^[0-9]+\.?[0-9]*)|([0-9]+\.?[0-9]*$))`)

func stripNumberFromString(v string) (string, error) {
	matches := bareNumberRegexp.FindStringSubmatch(v)
	if matches == nil {
		return "", fmt.Errorf("invalid number to strip:%s", v)
	}
	return matches[1], nil
}
