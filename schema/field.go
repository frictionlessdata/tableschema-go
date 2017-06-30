package schema

// Default for schema fields.
const (
	defaultFieldType   = "string"
	defaultFieldFormat = "default"
)

// Field represents a cell on a table.
type Field struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format"`
}

func setDefaultValues(f *Field) {
	if f.Type == "" {
		f.Type = defaultFieldType
	}
	if f.Format == "" {
		f.Format = defaultFieldFormat
	}
}
