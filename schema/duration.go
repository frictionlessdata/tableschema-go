package schema

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var durationRegexp = regexp.MustCompile(
	`P(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+\.?\d*S)?`)

const (
	hoursInYear  = time.Duration(24*360) * time.Hour
	hoursInMonth = time.Duration(24*30) * time.Hour
	hoursInDay   = time.Duration(24) * time.Hour
)

func castDuration(value string) (time.Duration, error) {
	matches := durationRegexp.FindStringSubmatch(value)
	if len(matches) == 0 {
		return 0, fmt.Errorf("Invalid duration:\"%s\"", value)
	}
	years := parseIntDuration(matches[1], hoursInYear)
	months := parseIntDuration(matches[2], hoursInMonth)
	days := parseIntDuration(matches[3], hoursInDay)
	hours := parseIntDuration(matches[4], time.Hour)
	minutes := parseIntDuration(matches[5], time.Minute)
	seconds := parseSeconds(matches[6])
	return years + months + days + hours + minutes + seconds, nil
}

func parseIntDuration(v string, multiplier time.Duration) time.Duration {
	if len(v) == 0 {
		return 0
	}
	// Ignoring error here because only digits could come from the regular expression.
	d, _ := strconv.Atoi(v[0 : len(v)-1])
	return time.Duration(d) * multiplier
}

func parseSeconds(v string) time.Duration {
	if len(v) == 0 {
		return 0
	}
	// Ignoring error here because only valid arbitrary precision floats could come from the regular expression.
	d, _ := strconv.ParseFloat(v[0:len(v)-1], 64)
	return time.Duration(d * 10e8)
}

func uncastDuration(in interface{}) (string, error) {
	v, ok := in.(time.Duration)
	if !ok {
		return "", fmt.Errorf("invalid duration - value:%v type:%v", in, reflect.ValueOf(in).Type())
	}
	y := v / hoursInYear
	r := v % hoursInYear
	m := r / hoursInMonth
	r = r % hoursInMonth
	d := r / hoursInDay
	r = r % hoursInDay
	return strings.ToUpper(fmt.Sprintf("P%dY%dM%dDT%s", y, m, d, r.String())), nil
}
