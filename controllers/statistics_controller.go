package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func attendanceStatistics(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	categoryIDsStr := r.FormValue("categoryIDs")
	var categoryIDs []uint64 = nil
	if categoryIDsStr != "" {
		categoryIDsStrings := strings.Split(categoryIDsStr, ";")
		categoryIDs = make([]uint64, len(categoryIDsStrings))
		for idx, categoryIDStr := range categoryIDsStrings {
			categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
			if err != nil {
				return Error(*lib.InvalidParamsError(fmt.Errorf(
					"Error parsing category ID %s! %w", categoryIDStr, err,
				)))
			}
			categoryIDs[idx] = categoryID
		}
	}
	data, serr := services.Attendance(db, r.FormValue("fromDate"), r.FormValue("toDate"), categoryIDs)
	if serr != nil {
		return Error(*serr)
	}

	return File("statistic.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// RegisterStatisticsControllerRoutes Registers the functions
func RegisterStatisticsControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/statistics/attendance", serviceWrapperDBAuthenticated("AttendanceStatistics", attendanceStatistics, config))
}
