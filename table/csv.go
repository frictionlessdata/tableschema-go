package table

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// CSV creates a Table from the CSV physical representation.
func CSV(source io.Reader, opts ...CreationOpts) (Table, error) {
	t := Table{Source: source}
	for _, opt := range opts {
		if err := opt(&t); err != nil {
			return Table{}, err
		}
	}
	return t, nil
}

// LoadCSVHeaders uses the first line of the CSV as table headers.
func LoadCSVHeaders() CreationOpts {
	return func(t *Table) error {
		r := bufio.NewReader(t.Source)
		t.Source = r
		var line string
		var err error
		for {
			line, err = r.ReadString('\n')
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			if strings.HasPrefix(line, "#") {
				continue
			}
			break
		}
		t.Headers = strings.Split(line[:len(line)-1], ",")
		for i, h := range t.Headers {
			t.Headers[i], err = strconv.Unquote(h)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// CSVFile creates a Table from a CSV local file.
func CSVFile(path string, opts ...CreationOpts) (Table, error) {
	f, err := os.Open(path)
	if err != nil {
		return Table{}, err
	}
	return CSV(f, opts...)
}

// CSVHeaders sets the table headers.
func CSVHeaders(headers ...string) CreationOpts {
	return func(t *Table) error {
		t.Headers = headers
		return nil
	}
}
