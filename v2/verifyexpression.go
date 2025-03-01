package datahelperlite

// VerifyExpression - verify expression for VerifyWithin function
type VerifyExpression struct {
	Name     string `json:"name,omitempty"`     // name of the database table column
	Value    any    `json:"value,omitempty"`    // value of the column
	Operator string `json:"operator,omitempty"` // operator of the validation
}
