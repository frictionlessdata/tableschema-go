package schema

import (
	"encoding/json"
	"fmt"
)

func castArray(value string) (interface{}, error) {
	var obj interface{}
	if err := json.Unmarshal([]byte(value), &obj); err != nil {
		return nil, err
	}
	arr, ok := obj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s is not an JSON array", value)
	}
	return arr, nil
}
