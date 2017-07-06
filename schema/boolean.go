package schema

import "fmt"

func castBoolean(value string, trueValues, falseValues []string) (bool, error) {
	for _, v := range trueValues {
		if value == v {
			return true, nil
		}
	}
	for _, v := range falseValues {
		if value == v {
			return false, nil
		}
	}
	return false, fmt.Errorf("invalid boolean value:%s", value)
}
