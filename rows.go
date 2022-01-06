package datahelperlite

// Rows datahelperlite rows interface
type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
	Columns() ([]Column, error)
	Values() ([]interface{}, error)
	RawValues() [][]byte
}
