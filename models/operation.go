package models

var (
	OperationT = &Operation{}
)

type Operation struct {
	FameModel `gorm:"embedded_prefix:opr_"`
	Title     string
	FirstName string
	LastName  string
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (o *Operation) ColumnPrefix() string {
	return "opr_"
}
