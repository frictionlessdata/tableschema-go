package schema

import "encoding/json"

func castObject(value string) (interface{}, error) {
	var obj interface{}
	if err := json.Unmarshal([]byte(value), &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func uncastObject(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	return string(b), err
}
