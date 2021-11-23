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
	user     *models.User
	present  int
	partial  int
	absent   int
	overlaps int
}

func Attendance(db *gorm.DB, fromDate string, toDate string, categoryIDs []uint64) ([]uint8, *lib.FameError) {
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

	userToSheet := map[uint64]*xlsx.Sheet{}
	stat := map[uint64]*userStat{}
	users := []*models.User{}
	err = db.Model(models.UserT).Find(&users).Error
	if err != nil {
		return nil, lib.DataCorruptionError(fmt.Errorf("Error fetching users from database! %w", err))
	}
	for _, u := range users {
		stat[u.ID] = &userStat{
			user:     u,
			present:  0,
			partial:  0,
			absent:   0,
			overlaps: 0,
		}

		userSheet, err := xls.AddSheet(u.FullName())
		if err != nil {
			return nil, lib.InternalError(fmt.Errorf("error adding sheet to excel file: %w", err))
		}
		userToSheet[u.ID] = userSheet

		row := userSheet.AddRow()
		row.AddCell().SetString("Termin")
		row.AddCell().SetString("Anwesend")
	}

	dates := []*models.Date{}
	query := db.Model(models.DateT).
		Preload("DateLogs").
		Preload("Category").
		Where(db.L(models.DateT, "StartTime").Ge(fromTime)).
		Where(db.L(models.DateT, "EndTime").Le(toTime)).
		Order(db.L(models.DateT, "StartTime").OrderAsc())

	if categoryIDs != nil {
		query = query.Where(
			db.L(models.DateT, "CategoryID").In(categoryIDs),
		)
	}

	err = query.Find(&dates).Error
	if err != nil {
		return nil, lib.DataCorruptionError(fmt.Errorf("Error fetching dates from database! %w", err))
	}

	row := sht.AddRow()
	row.AddCell().SetString("Name")
	row.AddCell().SetString("Anwesend")
	row.AddCell().SetString("Teilweise anwesend")
	row.AddCell().SetString("Abwesend")
	row.AddCell().SetString("Unbekannt")

	absent := map[uint64]map[uint64]bool{}
	for i, dt := range dates {
		overlap := false
		var prevDateID uint64
		if i > 0 {
			prevDateID = dates[i-1].ID
			prevDay, prevMon, prevYear := dates[i-1].StartTime.Date()
			curDay, curMon, curYear := dt.StartTime.Date()
			if prevDay == curDay && prevMon == curMon && prevYear == curYear {
				if !dt.EndTime.Before(dates[i-1].StartTime) && !dates[i-1].EndTime.Before(dt.StartTime) {
					overlap = true
				}
			}
		}

		logged := map[uint64]bool{}
		for _, dl := range dt.DateLogs {
			s, ok := stat[dl.UserID]
			logged[dl.UserID] = true
			if absent[dl.UserID] == nil {
				absent[dl.UserID] = map[uint64]bool{}
			}

			userSheet := userToSheet[dl.UserID]
			prevRow := userSheet.Rows[len(userSheet.Rows)-1]
			userRow := userSheet.AddRow()
			userRow.AddCell().SetString(dt.Title + " " + dt.StartTime.Format("2006-01-02"))
			if ok {
				if !dl.Present {
					if overlap && !absent[dl.UserID][prevDateID] {
						// user was at overlapped date
						s.overlaps++
						userRow.AddCell().SetString("")

					} else if overlap && absent[dl.UserID][prevDateID] {
						// user was not at overlapped date and not at dt
						s.overlaps++
						userRow.AddCell().SetString("")

					} else {
						s.absent++
						userRow.AddCell().SetInt(0)
					}
					absent[dl.UserID][dt.ID] = true

				} else if (dt.Category == nil || dt.Category.Name != models.OperationName) && (dl.FromTime.After(dt.StartTime) || dl.UntilTime.Before(dt.EndTime)) {
					if overlap && absent[dl.UserID][prevDateID] {
						s.absent--
						s.overlaps++
						prevRow.Cells[1].SetString("")
					}
					s.partial++
					userRow.AddCell().SetInt(1)
				} else {
					if overlap && absent[dl.UserID][prevDateID] {
						s.absent--
						s.overlaps++
						prevRow.Cells[1].SetString("")
					}
					s.present++
					userRow.AddCell().SetInt(1)
				}
			}
		}

		for userID, userSheet := range userToSheet {
			if !logged[userID] {
				userRow := userSheet.AddRow()
				userRow.AddCell().SetString(dt.Title + " " + dt.StartTime.Format("2006-01-02"))
				userRow.AddCell().SetString("?")
			}
		}
	}

	for _, s := range stat {
		row := sht.AddRow()
		row.AddCell().SetString(s.user.FullName())
		row.AddCell().SetInt(s.present)
		row.AddCell().SetInt(s.partial)
		row.AddCell().SetInt(s.absent)
		row.AddCell().SetInt(len(dates) - (s.present + s.partial + s.absent + s.overlaps))
	}

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err = xls.Write(writer)
	if err != nil {
		return nil, lib.InternalError(fmt.Errorf("Error writing XLSX data! %w", err))
	}

	return buf.Bytes(), nil
}
