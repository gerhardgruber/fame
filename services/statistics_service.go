package services

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/jinzhu/gorm"
	"github.com/tealeg/xlsx"
)

type userStat struct {
	user    *models.User
	present int
	partial int
	absent  int
}

func Attendance(db *gorm.DB, fromDate string, toDate string) ([]uint8, *lib.FameError) {
	const timeFormat = "Mon Jan 02 2006 15:04:05 GMT-0700"

	fromTime, err := time.Parse(timeFormat, fromDate)
	if err != nil {
		return nil, lib.InvalidParamsError(fmt.Errorf("Error parsing from date %s! %w", fromDate, err))
	}

	toTime, err := time.Parse(timeFormat, toDate)
	if err != nil {
		return nil, lib.InvalidParamsError(fmt.Errorf("Error parsing to date %s! %w", toDate, err))
	}

	xls := xlsx.NewFile()
	sht, err := xls.AddSheet("fame Statistik")
	if err != nil {
		return nil, lib.InternalError(fmt.Errorf("Error adding sheet to excel file! %w", err))
	}

	stat := map[uint64]*userStat{}
	users := []*models.User{}
	err = db.Model(models.UserT).Find(&users).Error
	if err != nil {
		return nil, lib.DataCorruptionError(fmt.Errorf("Error fetching users from database! %w", err))
	}
	for _, u := range users {
		stat[u.ID] = &userStat{
			user:    u,
			present: 0,
			partial: 0,
			absent:  0,
		}
	}

	dates := []*models.Date{}
	err = db.Model(models.DateT).
		Preload("DateLogs").
		Where(db.L(models.DateT, "StartTime").Ge(fromTime)).
		Where(db.L(models.DateT, "EndTime").Le(toTime)).
		Find(&dates).Error
	if err != nil {
		return nil, lib.DataCorruptionError(fmt.Errorf("Error fetching dates from database! %w", err))
	}

	row := sht.AddRow()
	row.AddCell().SetString("Name")
	row.AddCell().SetString("Anwesend")
	row.AddCell().SetString("Teilweise anwesend")
	row.AddCell().SetString("Abwesend")
	row.AddCell().SetString("Unbekannt")

	for _, dt := range dates {
		for _, dl := range dt.DateLogs {
			s, ok := stat[dl.UserID]
			if ok {
				if !dl.Present {
					s.absent++
				} else if dl.FromTime.After(dt.StartTime) || dl.UntilTime.Before(dt.EndTime) {
					s.partial++
				} else {
					s.present++
				}
			}
		}
	}

	for _, s := range stat {
		row := sht.AddRow()
		row.AddCell().SetString(s.user.FullName())
		row.AddCell().SetInt(s.present)
		row.AddCell().SetInt(s.partial)
		row.AddCell().SetInt(s.absent)
		row.AddCell().SetInt(len(dates) - (s.present + s.partial + s.absent))
	}

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err = xls.Write(writer)
	if err != nil {
		return nil, lib.InternalError(fmt.Errorf("Error writing XLSX data! %w", err))
	}

	return buf.Bytes(), nil
}
