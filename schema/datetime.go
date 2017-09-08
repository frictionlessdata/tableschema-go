package schema

import (
	"fmt"
	"reflect"
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
	return castDefaultOrCustomTime("2006-01-02", format, value)
}

func castTime(format, value string) (time.Time, error) {
	return castDefaultOrCustomTime("03:04:05", format, value)
}

func castYearMonth(value string) (time.Time, error) {
	return time.Parse("2006-01", value)
}

func castYear(value string) (time.Time, error) {
	return time.Parse("2006", value)
}

func castDateTime(format, value string) (time.Time, error) {
	return castDefaultOrCustomTime(time.RFC3339, format, value)
}

func castDefaultOrCustomTime(defaultFormat, format, value string) (time.Time, error) {
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

func encodeTime(v interface{}) (string, error) {
	value, ok := v.(time.Time)
	if !ok {
		return "", fmt.Errorf("invalid date - value:%v type:%v", v, reflect.ValueOf(v).Type())
	}
	utc := value.In(time.UTC)
	return utc.Format(time.RFC3339), nil
}
