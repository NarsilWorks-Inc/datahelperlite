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
	Open(ctx context.Context, di *cfg.DatabaseInfo) error
	Close() error
	Begin() error
	Commit() error
	Rollback() error
	Mark(name string) error
	Discard(name string) error
	Save(name string) error
	Query(sql string, args ...interface{}) (Rows, error)
	QueryRow(sql string, args ...interface{}) Row
	Exec(sql string, args ...interface{}) (int64, error)
	VerifyWithin(tablename string, values []std.VerifyExpression) (Valid bool, QueryOK bool, Message string)
	Escape(fv string) string
}

// ReadType - read types in data retrieval
type ReadType string

// ReturnKind - kinds of return in data retrieval
type ReturnKind string

// ReadTypes for data access
const (
	READALL     ReadType = `all`
	READBYKEY   ReadType = `key`
	READBYCODE  ReadType = `code`
	READFORFORM ReadType = `form`
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

	var paramchar string

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
			retstr = strings.Replace(retstr, match, paramchar+strconv.Itoa((i+1)), 1)
		} else {
			retstr = strings.Replace(retstr, match, paramchar, 1) // replace one at a time
		}
	}

	return retstr
}
