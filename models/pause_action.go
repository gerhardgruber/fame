package models

import "time"

type PauseType uint8

const (
	TrainingPause  PauseType = 0
	OperationPause PauseType = 1
)

var PauseActionT = &PauseAction{}

type PauseAction struct {
	FameModel `gorm:"embedded_prefix:pa_"`
	UserID    uint64
	User      *User
	StartTime *time.Time
	EndTime   *time.Time
	Type      PauseType
}

// ColumnPrefix implements the gorm columnPrefixer interface
// and returns the column prefix
func (pa *PauseAction) ColumnPrefix() string {
	return "pa_"
}
