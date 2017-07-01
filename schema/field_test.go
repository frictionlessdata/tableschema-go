package schema

import (
	"reflect"
	"testing"
)

func TestCastValue_Integer(t *testing.T) {
	f := Field{Type: "integer"}
	c, err := f.CastValue("42")
	if err != nil {
		t.Errorf("[Field.CastValue(integer)] err want:nil, got:%q", err)
	}
	intValue, ok := c.(int64)
	if !ok {
		t.Errorf("[Field.CastValue(integer)] cast want:int64, got:%s", reflect.TypeOf(c))
	}

	if intValue != 42 {
		t.Errorf("[Field.CastValue(integer)] val want:42, got:%d", intValue)
	}
}
