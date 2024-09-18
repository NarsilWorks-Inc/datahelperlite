package datahelperlite

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	cfg "github.com/eaglebush/config"
	pgr "github.com/eaglebush/pager"
)

// DataHelperLite interface for usage
type DataHelperLite interface {
	NewHelper() DataHelperLite
	Begin() error                                                                       // Begin a transaction. If there is an existing transaction, begin is ignored
	BeginDR() (string, error)                                                           // Begin a transaction with transaction id to use when rollback is deferred
	Commit(...string) error                                                             // Commit the transaction
	Close() error                                                                       // Close connection
	Discard(name string) error                                                          // Discard a savepoint
	Escape(fv string) string                                                            // Escape a field value (fv) from disruption by single quote
	Exec(sql string, args ...interface{}) (int64, error)                                // Exec executes a non-returning query
	Exists(sqlWithParams string, args ...interface{}) (bool, error)                     // Checks existence of a record
	Mark(name string) error                                                             // Mark a savepoint
	Next(serial string, next *int64) error                                              // Get next value of a serial
	Open(ctx context.Context, di *cfg.DatabaseInfo) error                               // Open a new connection
	Query(sql string, args ...interface{}) (Rows, error)                                // Query to a database to return one or more records
	QueryArray(sql string, out interface{}, args ...interface{}) error                  // Query to a database to return one or more records and store to an array
	QueryRow(sql string, args ...interface{}) Row                                       // QueryRow to a database and return one record
	QueryPaged(pinfo pgr.Parameter, sql string, args ...interface{}) (Rows, error)      // Query with paged results
	Rollback(...string) error                                                           // Rollback a transaction
	Save(name string) error                                                             // Save a transaction
	VerifyWithin(tableName string, values []VerifyExpression) (Valid bool, Error error) // VerifyWithin a set of validation expression against the underlying database table
	DatabaseVersion() string                                                            // Get database version
	Now() *time.Time                                                                    // Get time now
	NowUTC() *time.Time                                                                 // Get the time in UTC
}

// ReadType - read types in data retrieval
type ReadType string

// ReturnKind - kinds of return in data retrieval
type ReturnKind string

// ReadTypes for data access
const (
	READALL           ReadType = `all`   // Read all
	READBYKEY         ReadType = `key`   // Read by key
	READBYLATERALKEYS ReadType = `lkeys` // Read by lateral keys
	READBYCODE        ReadType = `code`  // Read by code
	READFORFORM       ReadType = `form`  // Read for form
	READELSE          ReadType = `else`  // Read else
)

// ReturnKind returns the data depends on kind
const (
	RETURNALL       ReturnKind = `all`
	RETURNFORFORM   ReturnKind = `form`
	RETURNESSENTIAL ReturnKind = `essential`
)

// Parameter Types
type (
	VarChar     string
	VarCharMax  string
	NVarCharMax string
)

type ParameterType interface {
	VarChar | VarCharMax | NVarCharMax
}

// Helper for datahelperlite
var Helper map[string]DataHelperLite

// Errors
var (
	ErrNoRows                error // ErrNoRows for no rows returned
	ErrArrayTypeNotSupported error = errors.New(`array type not supported`)
	ErrNoConn                error = errors.New(`no connection of the object was initialized`)
	ErrNoTx                  error = errors.New(`no transaction was initialized`)
	ErrNoPagerSet            error = errors.New(`no pager was set or initialized`)
	ErrVarMustBeInit         error = errors.New(`variable in next parameter must be initialized`)
)

// New creates new datahelper lite
func New(dhl *DataHelperLite, helperId string) (DataHelperLite, error) {
	var (
		ndh DataHelperLite
	)
	// copy existing postgresql helper
	if dhl != nil {
		ndh = *dhl
	}
	if ndh == nil {
		ndhi, present := Helper[helperId]
		if !present {
			return nil, errors.New(`no helper name of such`)
		}
		// create new helper instance
		ndh = ndhi.NewHelper()
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

// InterpolateTable interpolates table name that has been enclosed with curly braces
func InterpolateTable(sql string, schema string) string {
	if schema != "" {
		schema = schema + `.`
	}
	re := regexp.MustCompile(`\{([a-zA-Z0-9\[\]\"\_\-]*)\}`)
	return re.ReplaceAllString(sql, schema+`$1`)
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

// ToDBType converts string or string types to desired DBType
func ToDBType[T ParameterType](value any) T {
	if value == nil {
		return GetZero[T]()
	}
	switch t := value.(type) {
	case string:
		return T(t)
	case *string:
		return T(*t)
	default:
		v := reflect.ValueOf(t)
		if reflect.TypeOf(t).Kind() == reflect.Ptr {
			v = v.Elem()
		}
		x := v.String()
		return T(x)
	}
}

// GetZero gets the zero value of the generic type
func GetZero[T any]() T {
	var result T
	return result
}
