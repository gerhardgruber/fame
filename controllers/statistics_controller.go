package controllers

import (
	"net/http"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func attendanceStatistics(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	data, serr := services.Attendance(db, r.FormValue("fromDate"), r.FormValue("toDate"))
	if serr != nil {
		return Error(*serr)
	}

	return File("statistic.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// RegisterStatisticsControllerRoutes Registers the functions
func RegisterStatisticsControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/statistics/attendance", serviceWrapperDBAuthenticated("AttendanceStatistics", attendanceStatistics, config))
}
