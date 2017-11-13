package schema

import "time"

func castDate(format, value string, c Constraints) (time.Time, error) {
	y, err := castDateWithoutChecks(format, value)
	if err != nil {
		return y, err
	}
	var max, min time.Time
	if c.Maximum != "" {
		max, err = castDateWithoutChecks(format, c.Maximum)
		if err != nil {
			return max, err
		}
	}
	if c.Minimum != "" {
		min, err = castDateWithoutChecks(format, c.Minimum)
		if err != nil {
			return min, err
		}
	}
	return checkConstraints(y, max, min, DateType)
}

func castDateWithoutChecks(format, value string) (time.Time, error) {
	return castDefaultOrCustomTime("2006-01-02", format, value)
}
