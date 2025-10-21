package datahelperlite

// Row datahelperlite row interface
type Row interface {
	Scan(dest ...any) error
}
