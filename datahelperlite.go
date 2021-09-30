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
	Begin() error                                                                                            // Begin a transaction. If there is an existing transaction, begin is ignored
	Commit() error                                                                                           // Commit the transaction
	Close() error                                                                                            // Close connection
	Discard(name string) error                                                                               // Discard a savepoint
	Escape(fv string) string                                                                                 // Escape a field value (fv) from disruption by single quote
	Exec(sql string, args ...interface{}) (int64, error)                                                     // Exec executes a non-returning query
	Exists(sqlwparams string, args ...interface{}) (bool, error)                                             // Checks existence of a record
	Mark(name string) error                                                                                  // Mark a savepoint
	Next(serial string, next *int64) error                                                                   // Get next value of a serial
	Open(ctx context.Context, di *cfg.DatabaseInfo) error                                                    // Open a new connection
	Query(sql string, args ...interface{}) (Rows, error)                                                     // Query to a database to return one or more records
	QueryArray(sql string, out interface{}, args ...interface{}) error                                       // Query to a database to return one or more records and store to an array
	QueryRow(sql string, args ...interface{}) Row                                                            // QueryRow to a database and return one record
	Rollback() error                                                                                         // Rollback a transaction
	Save(name string) error                                                                                  // Save a transaction
	VerifyWithin(tablename string, values []std.VerifyExpression) (Valid bool, QueryOK bool, Message string) // VerifyWithin a set of validation expression against the underlying database table
	DatabaseVersion() string                                                                                 // Get database version
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

// Errors
var (
	ErrNoRows                error // ErrNoRows for no rows returned
	ErrNoConn                error = errors.New(`no connection of the object was initialized`)
	ErrNoTx                  error = errors.New(`no transaction was initialized`)
	ErrVarMustBeInit         error = errors.New(`variable in next parameter must be initialized`)
	ErrArrayTypeNotSupported error = errors.New(`array type not supported`)
)

// New creates new datahelper lite
func New(dhl *DataHelperLite, helperid string) (DataHelperLite, error) {

	var (
		ndh     DataHelperLite
		present bool
	)

	// copy existing postgresql helper
	if dhl != nil {
		ndh = *dhl
	}

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

// SetErrNoRows sets the error when there are no rows
func SetErrNoRows(err error) {
	ErrNoRows = err
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
