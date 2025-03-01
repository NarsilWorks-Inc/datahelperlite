package datahelperlite

import (
	"testing"
)

func TestConvert(t *testing.T) {

	sample := "Sample"

	type string2 string

	var sample2 string2 = string2("Hey")

	result1 := ToDBType[VarChar](sample)
	t.Logf("%T %v", result1, result1)

	result2 := ToDBType[NVarCharMax](sample2)
	t.Logf("%T %v", result2, result2)

	result3 := ToDBType[NVarCharMax](&sample2)
	t.Logf("%T %v", result3, result3)
}
