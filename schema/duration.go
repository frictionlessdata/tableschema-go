package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var durationRegexp = regexp.MustCompile(
	`P(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)

const (
	hoursInYear  = time.Duration(24*36) * time.Hour
	hoursInMonth = time.Duration(24*30) * time.Hour
	hoursInDay   = time.Duration(24) * time.Hour
)

func castDuration(value string) (time.Duration, error) {
	matches := durationRegexp.FindStringSubmatch(value)
	if len(matches) == 0 {
		return 0, invalidDurationError(value)
	}
	years, err := parseDuration(matches[1], hoursInYear)
	if err != nil {
		return 0, invalidDurationError(value)
	}
	months, err := parseDuration(matches[2], hoursInMonth)
	if err != nil {
		return 0, invalidDurationError(value)
	}
	days, err := parseDuration(matches[3], hoursInDay)
	if err != nil {
		return 0, invalidDurationError(value)
	}
	hours, err := parseDuration(matches[4], time.Hour)
	if err != nil {
		return 0, invalidDurationError(value)
	}
	minutes, err := parseDuration(matches[5], time.Minute)
	if err != nil {
		return 0, invalidDurationError(value)
	}
	seconds, err := parseDuration(matches[6], time.Second)
	if err != nil {
		return 0, invalidDurationError(value)
	}
	return years + months + days + hours + minutes + seconds, nil
}

func parseDuration(v string, multiplier time.Duration) (time.Duration, error) {
	if len(v) == 0 {
		return 0, nil
	}
	d, err := strconv.ParseFloat(v[0:len(v)-1], 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(d) * multiplier, nil
}

func invalidDurationError(v string) error {
	return fmt.Errorf("Invalid duration:\"%s\"", v)
}
