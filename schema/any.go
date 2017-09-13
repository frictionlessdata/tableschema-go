package schema

import "fmt"

func castAny(value interface{}) (interface{}, error) {
	return value, nil
}

func encodeAny(value interface{}) (string, error) {
	return fmt.Sprintf("%v", value), nil
}
