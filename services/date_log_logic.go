package services

import (
	"fmt"
	"time"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
)

type CreateUpdateDateLogParams struct {
	DateID    uint64
	UserID    uint64
	FromTime  time.Time
	UntilTime time.Time
	Present   bool
	Comment   string
}

// GetDateLogs loads all date logs
func GetDateLogs(db *gorm.DB, dateID uint64) ([]*models.DateLog, *lib.FameError) {
	dateLogs := []*models.DateLog{}

	if err := db.Model(models.DateLogT).Where(
		db.L(models.DateLogT, "DateID").Eq(dateID),
	).Preload("Date").Preload("User").Find(dateLogs).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get date logs: %s", err),
		)
	}

	return dateLogs, nil
}

func CreateDateLog(c *lib.Config, db *gorm.DB, p *CreateUpdateDateLogParams) (*models.DateLog, *lib.FameError) {
	dateLog := &models.DateLog{
		DateID:    p.DateID,
		UserID:    p.UserID,
		FromTime:  p.FromTime,
		UntilTime: p.UntilTime,
		Present:   p.Present,
		Comment:   p.Comment,
	}

	if err := db.Create(dateLog).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not create date log: %s", err),
		)
	}

	return dateLog, nil
}

func UpdateDateLog(c *lib.Config, db *gorm.DB, id uint64, p *CreateUpdateDateLogParams) (*models.DateLog, *lib.FameError) {
	dateLog, ferr := GetDateLogByID(db, id)
	if ferr != nil {
		return nil, ferr
	}

	dateLog.FromTime = p.FromTime
	dateLog.UntilTime = p.UntilTime
	dateLog.Present = p.Present
	dateLog.Comment = p.Comment

	if err := db.Save(dateLog).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not update date log: %s", err),
		)
	}

	return dateLog, nil
}

// GetDateLogByID loads a date log with the given ID
func GetDateLogByID(db *gorm.DB, id uint64) (*models.DateLog, *lib.FameError) {
	dateLog := &models.DateLog{}
	if err := db.Preload("User").Preload("Date").First(dateLog, id).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get date log %d: %s", id, err),
		)
	}

	return dateLog, nil
}
