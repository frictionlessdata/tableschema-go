package schema

import (
	"encoding/json"
	"io"
)

// Schema describes tabular data.
type Schema struct {
	Fields []Field `json:"fields"`
}

// Headers returns the headers of the tabular data described
// by the schema.
func (s Schema) Headers() []string {
	var h []string
	for i := range s.Fields {
		h = append(h, s.Fields[i].Name)
	}
	return h
}

// Read reads and parses a descriptor to create a schema.
//
// Example - Reading a schema from a file:
//
//  f, err := os.Open("foo/bar/schema.json")
//  if err != nil {
//    panic(err)
//  }
//  s, err := Read(f)
//  if err != nil {
//    panic(err)
//  }
//  fmt.Println(s)
func Read(r io.Reader) (*Schema, error) {
	var s Schema
	dec := json.NewDecoder(r)
	if err := dec.Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
