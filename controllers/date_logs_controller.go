package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func getDateLogs(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dateID, err := strconv.ParseUint(params["date_id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date ID '%s' could not be parsed: %s", params["date_id"], err),
		))
	}

	dateLogs, serr := services.GetDateLogs(db, dateID)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateLogs": dateLogs,
	})
}

func createDateLog(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	decoder := json.NewDecoder(r.Body)
	p := &services.CreateUpdateDateLogParams{}
	err := decoder.Decode(p)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding date log: %s", err),
		))
	}

	dtl, serr := services.CreateDateLog(
		c,
		db,
		p,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateLog": dtl,
	})
}

func getDateLog(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dateID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date Log ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	dateLog, serr := services.GetDateLogByID(db, dateID)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateLog": dateLog,
	})
}

func updateDateLog(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date Log ID could not be parsed: %s", err),
		))
	}

	decoder := json.NewDecoder(r.Body)
	p := &services.CreateUpdateDateLogParams{}
	err = decoder.Decode(p)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding date log: %s", err),
		))
	}

	dateLog, serr := services.UpdateDateLog(
		c,
		db,
		id,
		p,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateLog": dateLog,
	})
}

func deleteDateLog(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date Log ID could not be parsed: %s", err),
		))
	}

	decoder := json.NewDecoder(r.Body)
	data := models.Date{}
	err = decoder.Decode(&data)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("InvalidObjectError %+v", err),
		))
	}

	dateLog, serr := services.GetDateLogByID(db, id)
	if serr != nil {
		return Error(*serr)
	}

	err = db.Delete(dateLog).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("DataBaseError %+v", err),
		))
	}

	return Success()
}

// RegisterDateLogsControllerRoutes Registers the functions
func RegisterDateLogsControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/dates/{date_id:[0-9]+}/date_logs", serviceWrapperDBAuthenticated("getDateLogs", getDateLogs, config)).Methods("GET")
	router.HandleFunc("/date_logs", serviceWrapperDBAuthenticated("createDateLog", createDateLog, config)).Methods("POST")
	router.HandleFunc("/date_logs/{id:[0-9]+}", serviceWrapperDBAuthenticated("getDateLog", getDateLog, config)).Methods("GET")
	router.HandleFunc("/date_logs/{id:[0-9]+}", serviceWrapperDBAuthenticated("updateDateLog", updateDateLog, config)).Methods("POST")
	router.HandleFunc("/date_logs/{id:[0-9]+}/delete", serviceWrapperDBAuthenticated("deleteDateLog", deleteDateLog, config)).Methods("POST")
}
