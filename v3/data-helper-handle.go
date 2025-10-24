package datahelperlite

import (
	"database/sql"
	"errors"
	"fmt"

	dn "github.com/eaglebush/datainfo"
)

var (
	Handler map[string]DataHelperHandle
)

// DataHelperHandle manages the handle to the database connection
//
// It manages the resident database connection for proper pooling
type DataHelperHandle interface {
	Open(di *dn.DataInfo) error
	Ping() error
	DB() *sql.DB
	DI() *dn.DataInfo
	Close() error
	Err() error
}

// Errors
var (
	ErrHandleNoConn    error = errors.New(`no connection of the object was initialized`)
	ErrHandleNoConnStr error = errors.New(`connection string not set`)
	ErrHandleNoHandle  error = errors.New(`no sql handle`)
)

// New creates new datahelper lite if the dhl parameter is null.
func NewHandle(helperId string) (DataHelperHandle, error) {
	hnd, present := Handler[helperId]
	if !present {
		return nil, fmt.Errorf("'%s' helper name is invalid", helperId)
	}
	return hnd, nil
}

// SetHandler sets the internal handler object
func SetHandler(name string, hndl DataHelperHandle) {
	if Handler == nil {
		Handler = make(map[string]DataHelperHandle)
	}
	Handler[name] = hndl
}
