package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
)

type CreateUpdateDateParams struct {
	Title       string
	Description string
	Location    string
	StartTime   time.Time
	EndTime     time.Time
	CategoryID  uint64
	Closed      bool
}

// GetDates loads all dates
func GetDates(db *gorm.DB, loadPastDates bool, search string) (dates *[]models.Date, serr *lib.FameError) {
	dates = &[]models.Date{}

	q := db.Order(db.L(models.DateT, "StartTime").OrderAsc())
	if !loadPastDates {
		q = q.Where(
			db.L(models.DateT, "EndTime").Gt(time.Now()),
		)
	}

	search = strings.TrimSpace(search)
	if search != "" {
		search = "%" + search + "%"
		q = q.Where(
			db.L(models.DateT, "Title").Like(search).Or(
				db.L(models.DateT, "Description").Like(search).Or(
					db.L(models.DateT, "LocationStr").Like(search),
				),
			),
		)
	}

	if err := q.Preload("Category").Preload("Location").Preload("CreatedBy").Preload("DateFeedbacks").Find(dates).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get dates: %s", err),
		)
	}

	return dates, nil
}

func CreateDate(c *lib.Config, db *gorm.DB, u *models.User, p *CreateUpdateDateParams) (*models.Date, *lib.FameError) {
	adr, ferr := GetOrCreateAddress(p.Location, db, c)
	if ferr != nil {
		return nil, ferr
	}

	date := &models.Date{
		Title:       p.Title,
		Description: p.Description,
		StartTime:   p.StartTime,
		EndTime:     p.EndTime,
		LocationID:  &adr.ID,
		LocationStr: p.Location,
		CategoryID:  p.CategoryID,
		CreatedByID: u.ID,
		Closed:      p.Closed,
	}

	if err := db.Create(date).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not create date: %s", err),
		)
	}

	return date, nil
}

func UpdateDate(c *lib.Config, db *gorm.DB, id uint64, p *CreateUpdateDateParams) (*models.Date, *lib.FameError) {
	date, ferr := GetDateByID(db, id)
	if ferr != nil {
		return nil, ferr
	}

	adr, ferr := GetOrCreateAddress(p.Location, db, c)
	if ferr != nil {
		return nil, ferr
	}

	date.Title = p.Title
	date.Description = p.Description
	date.StartTime = p.StartTime
	date.EndTime = p.EndTime
	date.LocationID = &adr.ID
	date.LocationStr = p.Location
	date.CategoryID = p.CategoryID
	date.Closed = p.Closed

	if err := db.Save(date).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not update date: %s", err),
		)
	}

	return date, nil
}

// GetDateByID loads a date with the given ID
func GetDateByID(db *gorm.DB, id uint64) (date *models.Date, ferr *lib.FameError) {
	date = &models.Date{}
	if err := db.Preload("Category").
		Preload("Location").
		Preload("CreatedBy").
		Preload("DateFeedbacks").
		Preload("DateLogs").
		First(date, id).Error; err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get date %d: %s", id, err),
		)
	}

	return date, nil
}

func UpdateDateFeedback(db *gorm.DB, id uint64, userID uint64, feedback models.DateFeedbackType) (df *models.DateFeedback, ferr *lib.FameError) {
	df = &models.DateFeedback{}
	err := db.Model(df).Where(db.L(df, "UserID").Eq(userID)).
		Where(db.L(df, "DateID").Eq(id)).
		First(df).Error

	if gorm.IsRecordNotFoundError(err) {
		df = &models.DateFeedback{
			DateID: id,
			UserID: userID,
		}

	} else if err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not get date feedback for user id %d and date id %d! %s", userID, id, err),
		)
	}

	df.Feedback = feedback

	err = db.Save(df).Error
	if err != nil {
		return nil, lib.DataCorruptionError(
			fmt.Errorf("Could not update date feedback %+v! %s", df, err),
		)
	}

	return df, nil
}
