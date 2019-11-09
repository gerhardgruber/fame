package models

import "time"

type DateType int

const (
	Training DateType = 1

	Other DateType = 99
)

var (
	DateT = &Date{}
)

var DateTypes = map[string]DateType{
	"Training": Training,
	"Other":    Other,
}

type Date struct {
	FameModel   `gorm:"embedded_prefix:dt_"`
	Title       string
	Description string

	CreatedByID uint64
	CreatedBy   *User

	LocationID *uint64
	Location   *Address

	DateFeedbacks []DateFeedback

	StartTime time.Time
	EndTime   time.Time

	CategoryID uint64
	Category   *DateCategory
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (dt *Date) ColumnPrefix() string {
	return "dt_"
}
