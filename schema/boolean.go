package schema

import (
	"fmt"
	"reflect"
)

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

func uncastBoolean(value interface{}, trueValues, falseValues []string) (string, error) {
	switch value.(type) {
	case bool:
		return fmt.Sprintf("%v", value), nil
	case string:
		for _, v := range trueValues {
			if value == v {
				return value.(string), nil
			}
		}
		for _, v := range falseValues {
			if value == v {
				return value.(string), nil
			}
		}
	}
	return "", fmt.Errorf("invalid boolean - value:\"%v\" type:%v", value, reflect.ValueOf(value).Type())
}
