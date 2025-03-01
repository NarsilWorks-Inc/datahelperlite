package datahelperlite

// Rows datahelperlite rows interface
type Rows interface {
	Close()                         // Close the rows
	Columns() ([]Column, error)     // Get columns
	Err() error                     // Get last error
	Next() bool                     // Get next number as specified in the upsert configuration
	RawValues() [][]byte            // Raw values in array
	Scan(dest ...interface{}) error // Put the result into destination variables
	Values() ([]interface{}, error) // Return the values
}
