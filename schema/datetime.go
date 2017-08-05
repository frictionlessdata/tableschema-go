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
	switch format {
	case "", defaultFieldFormat:
		return time.Parse("2006-01-02", value)
	case AnyDateFormat:
		return time.Unix(0, 0), fmt.Errorf("any date format not yet supported. Please file an issue at github.com/frictionlessdata/tableschema-go")
	}
	goFormat := format
	for f, s := range strftimeToGoConversionTable {
		goFormat = strings.Replace(goFormat, f, s, -1)
	}
	return time.Parse(goFormat, value)
}

func castTime(format, value string) (time.Time, error) {
	switch format {
	case "", defaultFieldFormat:
		return time.Parse("03:04:05", value)
	case AnyDateFormat:
		return time.Unix(0, 0), fmt.Errorf("any date format not yet supported. Please file an issue at github.com/frictionlessdata/tableschema-go")
	}
	goFormat := format
	for f, s := range strftimeToGoConversionTable {
		goFormat = strings.Replace(goFormat, f, s, -1)
	}
	return time.Parse(goFormat, value)
}
