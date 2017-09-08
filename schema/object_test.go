package schema

import "testing"

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
			got, err := encodeObject(d.value)
			if err != nil {
				t.Fatalf("err want:nil got:%q", err)
			}
			if d.want != got {
				t.Fatalf("val want:%s got:%s", d.want, got)
			}
		})
	}
}
