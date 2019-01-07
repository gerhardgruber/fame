package controllers

import (
	"net/http"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func getOperations(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	operations, serr := services.GetOperations(db)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"Operations": operations,
	})
}

// RegisterOperationsControllerRoutes Registers the functions
func RegisterOperationsControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/operations", serviceWrapperDBAuthenticated("getOperations", getOperations, config)).Methods("GET")
}
