package schema

import "time"

func decodeDate(format, value string, c Constraints) (time.Time, error) {
	y, err := decodeDateWithoutChecks(format, value)
	if err != nil {
		return y, err
	}
	var max, min time.Time
	if c.Maximum != "" {
		max, err = decodeDateWithoutChecks(format, c.Maximum)
		if err != nil {
			return max, err
		}
	}
	if c.Minimum != "" {
		min, err = decodeDateWithoutChecks(format, c.Minimum)
		if err != nil {
			return min, err
		}
	}
	return checkConstraints(y, max, min, DateType)
}

func decodeDateWithoutChecks(format, value string) (time.Time, error) {
	return decodeDefaultOrCustomTime("2006-01-02", format, value)
}
