package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/frictionlessdata/tableschema-go/table"
)

// Table represents a Table backed by a CSV physical representation.
type Table struct {
	headers     []string
	source      Source
	skipHeaders bool
}

// NewTable creates a table.Table from the CSV table physical representation.
// CreationOpts are executed in the order they are declared.
func NewTable(source Source, opts ...CreationOpts) (*Table, error) {
	t := Table{source: source}
	for _, opt := range opts {
		if err := opt(&t); err != nil {
			return nil, err
		}
	}
	return &t, nil
}

// Iter returns an Iterator to read the table. Iter returns an error
// if the table physical source can not be iterated.
// The iteration process always start at the beginning of the CSV and
// is backed by a new reading.
func (table *Table) Iter() (table.Iterator, error) {
	src, err := table.source()
	if err != nil {
		return nil, err
	}
	return newIterator(src, table.skipHeaders), nil
}

// ReadAll reads all rows from the table and return it as strings.
func (table *Table) ReadAll() ([][]string, error) {
	var r [][]string
	iter, err := table.Iter()
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	for iter.Next() {
		r = append(r, iter.Row())
	}
	return r, nil
}

// Headers returns the headers of the tabular data.
func (table *Table) Headers() []string {
	return table.headers
}

// String returns a string version of the table.
func (table *Table) String() string {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	rows, err := table.ReadAll()
	if err != nil {
		return ""
	}
	w.WriteAll(rows)
	return buf.String()
}

func newIterator(source io.ReadCloser, skipHeaders bool) *csvIterator {
	return &csvIterator{
		source:      source,
		reader:      csv.NewReader(source),
		skipHeaders: skipHeaders,
	}
}

type csvIterator struct {
	reader *csv.Reader
	source io.ReadCloser

	current     []string
	err         error
	skipHeaders bool
}

func (i *csvIterator) Next() bool {
	if i.err != nil {
		return false
	}
	var err error
	i.current, err = i.reader.Read()
	if err != io.EOF {
		i.err = err
	}
	if i.skipHeaders {
		i.skipHeaders = false
		i.Next()
	}
	return err == nil
}

func (i *csvIterator) Row() []string {
	return i.current
}

func (i *csvIterator) Err() error {
	return i.err
}

func (i *csvIterator) Close() error {
	return i.source.Close()
}

// CreationOpts defines functional options for creating Tables.
type CreationOpts func(t *Table) error

// Source defines a table physical data source.
type Source func() (io.ReadCloser, error)

// FromFile defines a file-based Source.
func FromFile(path string) Source {
	return func() (io.ReadCloser, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

var (
	httpClient *http.Client
	once       sync.Once
)

const remoteFetchTimeoutSecs = 15

// Remote fetches the source schema from a remote URL.
func Remote(url string) Source {
	return func() (io.ReadCloser, error) {
		once.Do(func() {
			httpClient = &http.Client{
				Timeout: remoteFetchTimeoutSecs * time.Second,
			}
		})
		resp, err := httpClient.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		return stringReadCloser(string(body)), nil
	}
}

// FromString defines a string-based source.
func FromString(str string) Source {
	return func() (io.ReadCloser, error) {
		return stringReadCloser(str), nil
	}
}

func stringReadCloser(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func errorSource() Source {
	return func() (io.ReadCloser, error) {
		return nil, fmt.Errorf("error source")
	}
}

// LoadHeaders uses the first line of the CSV as table headers.
// The header line will be skipped during iteration
func LoadHeaders() CreationOpts {
	return func(reader *Table) error {
		iter, err := reader.Iter()
		if err != nil {
			return err
		}
		if iter.Next() {
			reader.headers = iter.Row()
		}
		reader.skipHeaders = true
		return nil
	}
}

// SetHeaders sets the table headers.
func SetHeaders(headers ...string) CreationOpts {
	return func(reader *Table) error {
		reader.headers = headers
		return nil
	}
}

func errorOpts(headers ...string) CreationOpts {
	return func(_ *Table) error {
		return fmt.Errorf("error opts")
	}
}

// NewWriter creates a writer which appends records to a CSV raw file.
//
// As returned by NewWriter, a csv.Writer writes records terminated by a
// newline and uses ',' as the field delimiter. The exported fields can be
// changed to customize the details before the first call to Write or WriteAll.
//
// Comma is the field delimiter.
//
// If UseCRLF is true, the csv.Writer ends each record with \r\n instead of \n.
func NewWriter(w io.Writer) *csv.Writer {
	return csv.NewWriter(w)
}
