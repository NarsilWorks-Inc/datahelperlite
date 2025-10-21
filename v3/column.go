package datahelperlite

import "reflect"

type Column interface {
	Name() string
	DatabaseTypeName() string
	ScanType() reflect.Type
}
