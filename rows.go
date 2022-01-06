package datahelperlite

import "database/sql"

// Rows datahelperlite rows interface
type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
	Columns() ([]sql.ColumnType, error)
	Values() ([]interface{}, error)
	RawValues() [][]byte
}