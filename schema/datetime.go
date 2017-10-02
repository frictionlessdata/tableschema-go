package schema

import (
	"fmt"
	"strings"
	"time"
)

// Go has a different date formatting style. Converting to the one
// used in https://specs.frictionlessdata.io/table-schema/#date
// https://docs.python.org/2/library/datetime.html#strftime-strptime-behavior
var strftimeToGoConversionTable = map[string]string{
	"%d":  "02",
	"%-d": "2",
	"%B":  "January",
	"%b":  "Jan",
	"%h":  "Jan",
	"%m":  "01",
	"%_m": " 1",
	"%-m": "1",
	"%Y":  "2006",
	"%y":  "06",
	"%H":  "15",
	"%I":  "03",
	"%M":  "04",
	"%S":  "05",
	"%f":  "999999",
	"%z":  "Z0700",
	"%:z": "Z07:00",
	"%Z":  "MST",
	"%p":  "PM",
}

func castDate(format, value string) (time.Time, error) {
	return decodeDefaultOrCustomTime("2006-01-02", format, value)
}

func decodeYearMonth(value string, c Constraints) (time.Time, error) {
	y, err := decodeYearMonthWithoutChecks(value)
	if err != nil {
		return time.Now(), err
	}
	var max, min time.Time
	if c.Maximum != "" {
		max, err = decodeYearMonthWithoutChecks(c.Maximum)
		if err != nil {
			return time.Now(), err
		}
	}
	if c.Minimum != "" {
		min, err = decodeYearMonthWithoutChecks(c.Minimum)
		if err != nil {
			return time.Now(), err
		}
	}
	return checkConstraints(y, max, min, YearMonthType)
}

func decodeYearMonthWithoutChecks(value string) (time.Time, error) {
	return time.Parse("2006-01", value)
}

func decodeYearWithoutChecks(value string) (time.Time, error) {
	return time.Parse("2006", value)
}

func decodeYear(value string, c Constraints) (time.Time, error) {
	y, err := decodeYearWithoutChecks(value)
	if err != nil {
		return time.Now(), err
	}
	var max, min time.Time
	if c.Maximum != "" {
		max, err = decodeYearWithoutChecks(c.Maximum)
		if err != nil {
			return time.Now(), err
		}
	}
	if c.Minimum != "" {
		min, err = decodeYearWithoutChecks(c.Minimum)
		if err != nil {
			return time.Now(), err
		}
	}
	return checkConstraints(y, max, min, YearType)
}

func checkConstraints(v, max, min time.Time, t string) (time.Time, error) {
	if !max.IsZero() && v.After(max) {
		return time.Now(), fmt.Errorf("constraint check error: %s:%v > maximum:%v", t, v, max)
	}
	if !min.IsZero() && v.Before(min) {
		return time.Now(), fmt.Errorf("constraint check error: %s:%v < minimum:%v", t, v, min)
	}
	return v, nil
}

func castDateTime(format, value string) (time.Time, error) {
	return decodeDefaultOrCustomTime(time.RFC3339, format, value)
}

func decodeDefaultOrCustomTime(defaultFormat, format, value string) (time.Time, error) {
	switch format {
	case "", defaultFieldFormat:
		t, err := time.Parse(defaultFormat, value)
		if err != nil {
			return time.Now(), err
		}
		return t.In(time.UTC), nil
	case AnyDateFormat:
		return time.Unix(0, 0), fmt.Errorf("any date format not yet supported. Please file an issue at github.com/frictionlessdata/tableschema-go")
	}
	goFormat := format
	for f, s := range strftimeToGoConversionTable {
		goFormat = strings.Replace(goFormat, f, s, -1)
	}
	t, err := time.Parse(goFormat, value)
	if err != nil {
		return time.Now(), err
	}
	return t.In(time.UTC), nil
}
