package models

import "time"

var DateLogT = &DateLog{}

type DateLog struct {
	FameModel `gorm:"embedded_prefix:dtl_"`

	DateID    uint64
	Date      *Date
	UserID    uint64
	User      *User
	FromTime  time.Time
	UntilTime time.Time
	Present   bool
	Comment   string `gorm:"size:4096"`
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (dtl *DateLog) ColumnPrefix() string {
	return "dtl_"
}
