// Package DataHelperLite
//
// v3.0
// 2025.10.21
package datahelperlite

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DataHelperLite interface
type DataHelperLite interface {
	NewHelper() DataHelperLite                                       // Create a new helper
	Acquire(ctx context.Context, h *DataHelperHandle) error          // Acquire sets all queries to a new context to isolate from pool context.
	Begin() error                                                    // Begin a transaction that supports deferred rollback.
	BeginManually() error                                            // Begin a transaction that should be committed or rolled back manually.
	Commit() error                                                   // Commit the transaction
	DatabaseVersion() string                                         // Get database version
	Discard(name string) error                                       // Discard a savepoint
	Escape(fv string) string                                         // Escape a field value (fv) from disruption by single quote
	Exec(sql string, args ...any) (int64, error)                     // Exec executes a non-returning query
	Exists(sqlWithParams string, args ...any) (bool, error)          // Checks existence of a record
	ExistsExt(tableName string, values []ColumnFilter) (bool, error) // Checks existence of a record by using a set of column filters against the underlying database table
	Mark(name string) error                                          // Mark a savepoint
	Next(serial string, next *int64) error                           // Get next value of a serial
	Now() *time.Time                                                 // Get time now
	NowUTC() *time.Time                                              // Get the time in UTC
	Ping() error                                                     // Ping the connection of the helper
	Query(sql string, args ...any) (Rows, error)                     // Query to a database to return one or more records
	QueryArray(sql string, out any, args ...any) error               // Query to a database to return one or more records and store to an array
	QueryRow(sql string, args ...any) Row                            // QueryRow to a database and return one record
	Rollback() error                                                 // Rollback a transaction
	Save(name string) error                                          // Save a transaction
}

// ReadType - read types in data retrieval
type ReadType string

// ReturnKind - kinds of return in data retrieval
type ReturnKind string

// ReadTypes for data access
const (
	ReadAll           ReadType = `all`   // Read all
	ReadByKey         ReadType = `key`   // Read by key
	ReadByLateralKeys ReadType = `lkeys` // Read by lateral keys
	ReadByCode        ReadType = `code`  // Read by code
	ReadElse          ReadType = `else`  // Read else
)

// ReturnKind returns the data depends on kind
const (
	ReturnAll       ReturnKind = `all`
	ReturnForForm   ReturnKind = `form`
	ReturnEssential ReturnKind = `essential`
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

var (
	Helper map[string]DataHelperLite
)

// Errors
var (
	ErrNoRows                error // ErrNoRows for no rows returned
	ErrArrayTypeNotSupported error = errors.New(`array type not supported`)
	ErrNoTx                  error = errors.New(`no transaction was initialized`)
	ErrVarMustBeInit         error = errors.New(`variable in next parameter must be initialized`)
	ErrHandleNotSet          error = errors.New(`handle not set`)
)

// New creates new datahelper lite if the dhl parameter is null.
func New(dhl DataHelperLite, helperId string) (DataHelperLite, error) {
	if dhl != nil {
		return dhl, nil
	}
	ndh, present := Helper[helperId]
	if !present {
		return nil, fmt.Errorf("'%s' helper name is invalid", helperId)
	}
	return ndh.NewHelper(), nil
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
	var result T
	if isReallyNil(value) {
		return result
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

func isReallyNil(i any) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	if !v.IsValid() {
		return true
	}
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return true
		}
		elem := v.Elem()
		if !elem.IsValid() {
			return true
		}
		switch elem.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
			return elem.IsNil()
		}
	}
	switch v.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return v.IsNil()
	}
	return false
}
