package schema

import (
	"fmt"
	"reflect"
	"time"
)

func decodeTime(format, value string, c Constraints) (time.Time, error) {
	y, err := decodeTimeWithoutCheckConstraints(format, value)
	if err != nil {
		return y, err
	}
	var max, min time.Time
	if c.Maximum != "" {
		max, err = decodeTimeWithoutCheckConstraints(format, c.Maximum)
		if err != nil {
			return y, err
		}
	}
	if c.Minimum != "" {
		min, err = decodeTimeWithoutCheckConstraints(format, c.Minimum)
		if err != nil {
			return y, err
		}
	}
	return checkConstraints(y, max, min, TimeType)
}

func decodeTimeWithoutCheckConstraints(format, value string) (time.Time, error) {
	return decodeDefaultOrCustomTime("03:04:05", format, value)
}

func encodeTime(v interface{}) (string, error) {
	value, ok := v.(time.Time)
	if !ok {
		return "", fmt.Errorf("invalid date - value:%v type:%v", v, reflect.ValueOf(v).Type())
	}
	utc := value.In(time.UTC)
	return utc.Format(time.RFC3339), nil
}
