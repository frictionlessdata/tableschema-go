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

func castYearMonth(value string, c Constraints) (time.Time, error) {
	y, err := castYearMonthWithoutChecks(value)
	if err != nil {
		return y, err
	}
	var max, min time.Time
	if c.Maximum != "" {
		max, err = castYearMonthWithoutChecks(c.Maximum)
		if err != nil {
			return y, err
		}
	}
	if c.Minimum != "" {
		min, err = castYearMonthWithoutChecks(c.Minimum)
		if err != nil {
			return y, err
		}
	}
	return checkConstraints(y, max, min, YearMonthType)
}

func castYearMonthWithoutChecks(value string) (time.Time, error) {
	return time.Parse("2006-01", value)
}

func castYearWithoutChecks(value string) (time.Time, error) {
	return time.Parse("2006", value)
}

func castYear(value string, c Constraints) (time.Time, error) {
	y, err := castYearWithoutChecks(value)
	if err != nil {
		return y, err
	}
	var max, min time.Time
	if c.Maximum != "" {
		max, err = castYearWithoutChecks(c.Maximum)
		if err != nil {
			return y, err
		}
	}
	if c.Minimum != "" {
		min, err = castYearWithoutChecks(c.Minimum)
		if err != nil {
			return y, err
		}
	}
	return checkConstraints(y, max, min, YearType)
}

func castDateTime(value string, c Constraints) (time.Time, error) {
	dt, err := castDateTimeWithoutChecks(value)
	if err != nil {
		return dt, err
	}
	var max, min time.Time
	if c.Maximum != "" {
		max, err = castDateTimeWithoutChecks(c.Maximum)
		if err != nil {
			return dt, err
		}
	}
	if c.Minimum != "" {
		min, err = castDateTimeWithoutChecks(c.Minimum)
		if err != nil {
			return dt, err
		}
	}
	return checkConstraints(dt, max, min, DateTimeType)
}

func castDateTimeWithoutChecks(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func checkConstraints(v, max, min time.Time, t string) (time.Time, error) {
	if !max.IsZero() && v.After(max) {
		return v, fmt.Errorf("constraint check error: %s:%v > maximum:%v", t, v, max)
	}
	if !min.IsZero() && v.Before(min) {
		return v, fmt.Errorf("constraint check error: %s:%v < minimum:%v", t, v, min)
	}
	return v, nil
}

func castDefaultOrCustomTime(defaultFormat, format, value string) (time.Time, error) {
	switch format {
	case "", defaultFieldFormat:
		t, err := time.Parse(defaultFormat, value)
		if err != nil {
			return t, err
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
		return t, err
	}
	return t.In(time.UTC), nil
}
