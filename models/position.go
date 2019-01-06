package models

var (
	// PositionT TODO: comment
	PositionT = &Position{}
)

// Position TODO: comment
type Position struct {
	FameModel     `gorm:"embedded_prefix:pos_"`
	MobilePhoneID uint64
	MobilePhone   MobilePhone
	Longitude     float64
	Latitude      float64
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (u *Position) ColumnPrefix() string {
	return "pos_"
}
