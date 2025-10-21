package datahelperlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dn "github.com/eaglebush/datainfo"
)

// DataHelperHandle manages the handle to the database connection
//
// It manages the resident database connection for proper pooling
type DataHelperHandle struct {
	db  *sql.DB
	dbi *dn.DataInfo
	err error
}

// Open connects to the database and initializes it
func (h *DataHelperHandle) Open(di *dn.DataInfo) error {
	h.dbi = di
	h.db, h.err = sql.Open(`sqlserver`, *di.ConnectionString)
	if h.err != nil {
		return fmt.Errorf("open: %w", h.err)
	}
	if di.MaxOpenConnection != nil {
		h.db.SetMaxOpenConns(*di.MaxOpenConnection)
	}
	h.db.SetMaxIdleConns(0)
	if di.MaxIdleConnection != nil {
		h.db.SetMaxIdleConns(*di.MaxIdleConnection)
	}
	if di.MaxConnectionLifetime != nil {
		h.db.SetConnMaxLifetime(time.Duration(*di.MaxConnectionLifetime))
	}
	if di.MaxConnectionIdleTime != nil {
		h.db.SetConnMaxIdleTime(time.Duration(*di.MaxConnectionIdleTime))
	}
	if err := h.db.PingContext(context.Background()); err != nil {
		h.err = fmt.Errorf("open: %w", err)
		return h.err
	}
	return nil
}

// Ping tests the database connection
func (h *DataHelperHandle) Ping() error {
	if h.db == nil {
		return fmt.Errorf("ping: no sql handle to use")
	}
	if err := h.db.PingContext(context.Background()); err != nil {
		h.err = fmt.Errorf("ping: %w", err)
		return h.err
	}
	return nil
}

// Close the database connection
func (h *DataHelperHandle) Close() error {
	if h.db == nil {
		return fmt.Errorf("ping: no sql handle to close")
	}
	if h.err = h.db.Close(); h.err != nil {
		return h.err
	}
	h.db = nil
	return nil
}
