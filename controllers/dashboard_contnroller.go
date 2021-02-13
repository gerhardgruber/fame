package controllers

import (
	"net/http"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func getStatus(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dates, positive, present, ferr := services.GetUserStatus(db, *sess.UserID)
	if ferr != nil {
		return Error(*ferr)
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
