// Package table provides the main interfaces used to manipulate tabular data.
// To understand why we need interfaces to process tabular data, it is useful to introduce
// the concepts of the physical and the logical representation of data.
//
// The physical representation of data refers to the representation of data as text on disk,
// for example, in a CSV, JSON or XML file. This representation may have some type information (JSON,
// where the primitive types that JSON supports can be used) or not (CSV, where all data is
// represented in string form). In this project, those are going to be presented as packages that
// provide structs which implement those interfaces. For instance, csv.NewTable creates a Table
// which is backed up by a CSV.
//
// The logical representation of data refers to the "ideal" representation of the data in terms of
// primitive types, data structures, and relations, all as defined by the specification. We could say
// that the specification is about the logical representation of data. That said, functions
// exported for data processing should deal with logic representations. That functionality
// is represented by interfaces in this package.
package table

import (
	"bytes"
	"encoding/csv"
)

// Table provides functionality to iterate and write tabular data. This is the logical
// representation and is meant to be encoding/format agnostic.
type Table interface {
	// Headers returns the headers of the tabular data.
	Headers() []string

	// Iter provides a convenient way to iterate over table's data.
	// The iteration process always start at the beginning of the table and
	// is backed by a new reading.
	Iter() (Iterator, error)

	// ReadAll reads all rows from the table and return it as strings.
	ReadAll() ([][]string, error)
}

// A Writer writes rows to a table file. The idea behind the writer is to
// abstract out the physical representation of the table. Which can have
// many formats, for instance, CSV, XML and JSON
type Writer interface {
	// Write writes a single row to w along with any necessary quoting.
	// A record is a slice of strings with each string being one field.
	Write(record []string) error
	// Flush writes any buffered data to the underlying io.Writer.
	// To check if an error occurred during the Flush, call Error.
	Flush()
	// Error reports any error that has occurred during a previous Write or Flush.
	Error() error
	// WriteAll writes multiple CSV records to w using Write and then calls Flush.
	WriteAll(records [][]string) error
}

// StringWriter is a simple Writer implementation which is backed up by
// an in memory bytes.Buffer.
type StringWriter struct {
	csv.Writer

	content *bytes.Buffer
}

// String returns the content that has been written so far as raw CSV.
func (s *StringWriter) String() string {
	return s.content.String()
}

// NewStringWriter returns a Writer that writes CSV to a string.
// It exports a String() method, which returns its contents.
func NewStringWriter() *StringWriter {
	buf := &bytes.Buffer{}
	return &StringWriter{*csv.NewWriter(buf), buf}
}

// FromSlices creates a new SliceTable using passed-in arguments.
func FromSlices(headers []string, content [][]string) *SliceTable {
	return &SliceTable{headers, content}
}

// SliceTable offers a simple table implementation backed by slices.
type SliceTable struct {
	headers []string
	content [][]string
}

// Headers returns the headers of the tabular data.
func (t *SliceTable) Headers() []string {
	return t.headers
}

// ReadAll reads all rows from the table and return it as strings.
func (t *SliceTable) ReadAll() ([][]string, error) {
	return t.content, nil
}

// Iter provides a convenient way to iterate over table's data.
// The iteration process always start at the beginning of the table and
// is backed by a new reading process.
func (t *SliceTable) Iter() (Iterator, error) {
	return &sliceIterator{content: t.content}, nil
}

type sliceIterator struct {
	content [][]string
	pos     int
}

func (i *sliceIterator) Next() bool {
	i.pos++
	return i.pos <= len(i.content)
}
func (i *sliceIterator) Row() []string { return i.content[i.pos-1] }
func (i *sliceIterator) Err() error    { return nil }
func (i *sliceIterator) Close() error  { return nil }

// Iterator is an interface which provides method to interating over tabular
// data. It is heavly inspired by bufio.Scanner.
// Iterating stops unrecoverably at EOF, the first I/O error, or a token too large to fit in the buffer.
type Iterator interface {
	// Next advances the table interator to the next row, which will be available through the Cast or Row methods.
	// It returns false when the iterator stops, either by reaching the end of the table or an error.
	// After Next returns false, the Err method will return any error that ocurred during the iteration, except if it was io.EOF, Err
	// will return nil.
	// Next could automatically buffer some data, improving reading performance. It could also block, if necessary.
	Next() bool

	// Row returns the most recent row fetched by a call to Next as a newly allocated string slice
	// holding its fields.
	Row() []string

	// Err returns nil if no errors happened during iteration, or the actual error
	// otherwise.
	Err() error

	// Close frees up any resources used during the iteration process.
	Close() error
}
