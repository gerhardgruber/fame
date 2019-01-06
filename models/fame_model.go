package models

import (
	"time"
)

// FameModel is used as base struct for fame models
type FameModel struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
