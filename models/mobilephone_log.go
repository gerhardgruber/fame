package models

var (
	// MobilePhoneLogT TODO: comment
	MobilePhoneLogT = &MobilePhoneLog{}
)

// MobilePhoneLog TODO: comment
type MobilePhoneLog struct {
	FameModel     `gorm:"embedded_prefix:mpl_"`
	MobilePhoneID uint64
	Message       string
	Context       string
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (ml *MobilePhoneLog) ColumnPrefix() string {
	return "mpl_"
}
