package datahelperlite

import (
	"database/sql"
	"errors"

	dn "github.com/eaglebush/datainfo"
)

var (
	Handler DataHelperHandle
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

func NewHandle() DataHelperHandle {
	return Handler
}
