package schema

import (
	"encoding/json"
	"io"
)

// Field represents a cell on a table.
type Field struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format"`
}

// Schema describes tabular data.
type Schema struct {
	Fields []Field `json:"fields"`
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
