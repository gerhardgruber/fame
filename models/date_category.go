package models

var (
	DateCategoryT = &DateCategory{}
)

type DateCategory struct {
	FameModel `gorm:"embedded_prefix:dc_"`
	Name      string
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (dt *DateCategory) ColumnPrefix() string {
	return "dc_"
}
