package models

import "time"

var (
	DateT = &Date{}
)

type Date struct {
	FameModel   `gorm:"embedded_prefix:dt_"`
	Title       string
	Description string

	CreatedByID uint64
	CreatedBy   *User

	LocationID *uint64
	Location   *Address

	LocationStr string

	DateFeedbacks []DateFeedback

	StartTime time.Time
	EndTime   time.Time

	CategoryID uint64
	Category   *DateCategory

	Closed bool

	DateLogs []*DateLog
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (dt *Date) ColumnPrefix() string {
	return "dt_"
}
