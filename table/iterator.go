package table

// Iterator is an interface which provides method to interating over tabular
// data. It is heavly inspired by bufio.Scanner.
type Iterator interface {
	// Next advances the table interator to the next row, which will be available through the Cast or Row methods.
	// It returns false when the iterator stops, either by reaching the end of the table or an error.
	// After Next returns false, the Err method will return any error that ocurred during the iteration, except if it was io.EOF, Err
	// will return nil.
	// Next could automatically buffer some data, improving reading performance. It could also block, if necessary.
	Next() bool

	// CastRow casts the most recent row fetched by a call to Next. Cast will error if the table has no schema
	// associateed to it of if the row can not be cast to its respective schema. More at Schema.CastRow.
	CastRow(out interface{}) error

	// Row returns the most recent row fetched by a call to Next as a newly allocated string slice
	// holding its fields.
	Row() []string

	// Err returns nil if no errors happened during iteration, or the actual error
	// otherwise.
	Err() error
}
