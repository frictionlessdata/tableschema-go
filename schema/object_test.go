package schema

import "testing"
import "github.com/matryer/is"

type eoStruct struct {
	Name string `json:"name"`
}

func TestEncodeObject(t *testing.T) {
	data := []struct {
		desc  string
		value interface{}
		want  string
	}{
		{"Simple", eoStruct{Name: "Foo"}, `{"name":"Foo"}`},
	}
	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			is := is.New(t)
			got, err := encodeObject(d.value)
			is.NoErr(err)
			is.Equal(d.want, got)
		})
	}
}
