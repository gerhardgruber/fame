package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func getStatus(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	startDate := time.Now().Add(-1 * time.Hour * 24 * 30)

	dates := 0
	err := db.Model(models.DateT).Where(
		db.L(models.DateT, "StartTime").Ge(startDate),
	).Where(
		db.L(models.DateT, "StartTime").Le(time.Now()),
	).Count(&dates).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("Error getting number of dates since %v! %w", startDate, err),
		))
	}

	positive := 0
	err = db.Model(models.DateT).Where(
		db.L(models.DateT, "StartTime").Ge(startDate),
	).Where(
		db.L(models.DateT, "StartTime").Le(time.Now()),
	).Joins(
		db.InnerJoin(models.DateFeedbackT).On(db.L(models.DateT, "ID"), db.L(models.DateFeedbackT, "DateID")),
	).Where(
		db.L(models.DateFeedbackT, "Feedback").Eq(models.DateFeedbackTypeYes),
	).Where(
		db.L(models.DateFeedbackT, "UserID").Eq(*sess.UserID),
	).Count(&positive).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("Error getting number of positive dates since %v! %w", startDate, err),
		))
	}

	present := 0
	err = db.Model(models.DateT).Where(
		db.L(models.DateT, "StartTime").Ge(startDate),
	).Where(
		db.L(models.DateT, "StartTime").Le(time.Now()),
	).Joins(
		db.InnerJoin(models.DateLogT).On(db.L(models.DateT, "ID"), db.L(models.DateLogT, "DateID")),
	).Where(
		db.L(models.DateLogT, "Present").Eq(true),
	).Where(
		db.L(models.DateLogT, "UserID").Eq(*sess.UserID),
	).Count(&present).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("Error getting number of present dates since %v! %w", startDate, err),
		))
	}

	return Success(map[string]interface{}{
		"Dates":    dates,
		"Positive": positive,
		"Present":  present,
	})
}

// RegisterDashboardControllerRoutes Registers the functions
func RegisterDashboardControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/dashboard/status", serviceWrapperDBAuthenticated("status", getStatus, config)).Methods("GET")
}
