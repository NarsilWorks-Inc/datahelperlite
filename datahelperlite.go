package datahelperlite

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"

	cfg "github.com/eaglebush/config"
	std "github.com/eaglebush/stdutil"
)

// DataHelperLite interface for usage
type DataHelperLite interface {
	Open(ctx context.Context, di *cfg.DatabaseInfo) error                                                    // Open a new connection
	Close() error                                                                                            // Close connection
	Begin() error                                                                                            // Begin a transaction. If there is an existing transaction, begin is ignored
	Commit() error                                                                                           // Commit a transaction
	Rollback() error                                                                                         // Rollback a transaction
	Mark(name string) error                                                                                  // Mark a savepoint
	Discard(name string) error                                                                               // Discard a savepoint
	Save(name string) error                                                                                  // Save a savepoint
	Query(sql string, args ...interface{}) (Rows, error)                                                     // Query to a database and return one or more records
	QueryRow(sql string, args ...interface{}) Row                                                            // QueryRow to a database and return one record
	Exec(sql string, args ...interface{}) (int64, error)                                                     // Exec executes a non-returning query
	VerifyWithin(tablename string, values []std.VerifyExpression) (Valid bool, QueryOK bool, Message string) // VerifyWithin a set of validation expression against the underlying database table
	Escape(fv string) string                                                                                 // Escape a field value (fv) from disruption by single quote
}

// ReadType - read types in data retrieval
type ReadType string

// ReturnKind - kinds of return in data retrieval
type ReturnKind string

// ReadTypes for data access
const (
	READALL           ReadType = `all`
	READBYKEY         ReadType = `key`
	READBYLATERALKEYS ReadType = `lkeys`
	READBYCODE        ReadType = `code`
	READFORFORM       ReadType = `form`
)

// ReturnKind returns the data depends on kind
const (
	RETURNALL       ReturnKind = `all`
	RETURNFORFORM   ReturnKind = `form`
	RETURNESSENTIAL ReturnKind = `essential`
)

// Helper for datahelperlite
var Helper map[string]DataHelperLite

// New creates new datahelper lite
func New(dhl *DataHelperLite, helperid string) (DataHelperLite, error) {

	var (
		ndh     DataHelperLite
		present bool
	)
	// copy existing postgresql helper
	ndh = *dhl

	if ndh == nil {
		ndh, present = Helper[helperid]
		if !present {
			return nil, errors.New(`No helper name of such`)
		}
	}

	return ndh, nil
}

// SetHelper - set helper object
func SetHelper(name string, dhl DataHelperLite) {
	if Helper == nil {
		Helper = make(map[string]DataHelperLite)
	}

	Helper[name] = dhl
}

// Row datahelperlite row interface
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows datahelperlite rows interface
type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
	Values() ([]interface{}, error)
	RawValues() [][]byte
}

// InterpolateTable interpolates table name that has been enclosed with curly braces
func InterpolateTable(sql string, schema string) string {
	if schema != "" {
		schema = schema + `.`
	}

	re := regexp.MustCompile(`\{([a-zA-Z0-9\[\]\"\_\-]*)\}`)
	sql = re.ReplaceAllString(sql, schema+`$1`)

	return sql
}

// ReplaceQueryParamMarker replaces SQL string with parameters set as ?
func ReplaceQueryParamMarker(preparedQuery string, paramInSeq bool, paramPlaceHolder string) string {

	retstr := preparedQuery
	defph := `?`
	pattern := `\` + defph //search for ?

	// if the paramPlaceHolder was set
	// by the configuration the same as default place holder, we exit
	if paramPlaceHolder == defph {
		return retstr
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(preparedQuery, -1)

	for i, match := range matches {
		if paramInSeq {
			retstr = strings.Replace(retstr, match, paramPlaceHolder+strconv.Itoa((i+1)), 1)
		} else {
			retstr = strings.Replace(retstr, match, paramPlaceHolder, 1) // replace one at a time
		}
	}

	return retstr
}
