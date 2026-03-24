package datahelperlite

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	dn "github.com/eaglebush/datainfo"
)

var (
	Handler map[string]DataHelperHandle
)

// DataHelperHandle manages the handle to the database connection
//
// It manages the resident database connection for proper pooling
type DataHelperHandle interface {
	// Open a database handle by providing database connection info
	Open(di *dn.DataInfo) error
	// Ping if the database connection is alive
	Ping() error
	// Returns the database handle
	DB() *sql.DB
	// Returns the database info
	DI() *dn.DataInfo
	// Close the handle
	Close() error
	// Err retrieves the last handle error
	Err() error
}

// Errors
var (
	ErrHandleNoConn    error = errors.New("no connection of the object was initialized")
	ErrHandleNoConnStr error = errors.New("connection string not set")
	ErrHandleNoHandle  error = errors.New("no sql handle")
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

// Reconnect allows reconnection to stabilize the handler.
//
// It returns a function to close the timer.
func Reconnect(
	hndl DataHelperHandle,
	interval time.Duration,
	mu sync.Locker,
	logf func(string, ...any),
) func() {
	stopCh := make(chan struct{})

	logger := func(format string, args ...any) {
		if logf != nil {
			logf(format, args...)
		}
	}

	lock := func() {
		if mu != nil {
			mu.Lock()
		}
	}

	unlock := func() {
		if mu != nil {
			mu.Unlock()
		}
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		justConnected := false
		connCount := 0

		for {
			select {
			case <-ticker.C:
				// Try connecting if there is no db handle
				if hndl.DB() == nil {
					di := hndl.DI()
					if di == nil {
						logger("Database error: nil DataInfo")
						continue
					}

					if err := hndl.Open(di); err != nil {
						logger("Database error: %v", err)
						continue
					}

					lock()
					connCount++
					unlock()

					justConnected = true
				}

				// Ping to check connection
				if err := hndl.Ping(); err != nil {
					logger("Database error: %v", err)

					lock()
					_ = hndl.Close()
					unlock()

					continue
				}

				if justConnected {
					lock()
					if connCount == 1 {
						logger("Database connection successful!")
					} else {
						logger("Database re-connection successful!")
					}
					unlock()

					justConnected = false
				}

			case <-stopCh:
				logger("Re-connection ticker stopped!")
				return
			}
		}
	}()

	return func() {
		close(stopCh)
	}
}
