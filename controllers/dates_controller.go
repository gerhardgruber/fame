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

func getDates(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dates, serr := services.GetDates(db, r.FormValue("loadPastDates") == "true")
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"Dates": dates,
	})
}

func getAppDates(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dates, serr := services.GetDates(db, r.FormValue("pastDates") == "true")
	if serr != nil {
		return Error(*serr)
	}

	users, ferr := services.GetUsers(db)
	if ferr != nil {
		return Error(*ferr)
	}

	return Success(map[string]interface{}{
		"Dates": dates,
		"Users": users,
	})
}

func createDate(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	decoder := json.NewDecoder(r.Body)
	p := &services.CreateUpdateDateParams{}
	err := decoder.Decode(p)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding date: %s", err),
		))
	}

	dt, serr := services.CreateDate(
		c,
		db,
		sess.User,
		p,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"Date": dt,
	})
}

func getDate(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dateID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	date, serr := services.GetDateByID(db, dateID)
	if serr != nil {
		return Error(*serr)
	}

	users, ferr := services.GetUsers(db)
	if ferr != nil {
		return Error(*ferr)
	}

	return Success(map[string]interface{}{
		"Date":  date,
		"Users": users,
	})
}

func getDateApp(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dateID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("DAte ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	date, serr := services.GetDateByID(db, dateID)
	if serr != nil {
		return Error(*serr)
	}

	users, ferr := services.GetUsers(db)
	if ferr != nil {
		return Error(*ferr)
	}

	return Success(map[string]interface{}{
		"Date":  date,
		"Users": users,
	})
}

func deleteDate(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date ID could not be parsed: %s", err),
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

	date, serr := services.GetDateByID(db, id)
	if serr != nil {
		return Error(*serr)
	}

	err = db.Delete(date).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("DataBaseError %+v", err),
		))
	}

	return Success()
}

func updateDate(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date ID could not be parsed: %s", err),
		))
	}

	decoder := json.NewDecoder(r.Body)
	p := &services.CreateUpdateDateParams{}
	err = decoder.Decode(p)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding date: %s", err),
		))
	}

	date, serr := services.UpdateDate(
		c,
		db,
		id,
		p,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"Date": date,
	})
}

func updateFeedback(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date ID could not be parsed: %s", err),
		))
	}

	type updateFeedbackParams struct {
		UserID   uint64
		Feedback models.DateFeedbackType
	}
	decoder := json.NewDecoder(r.Body)
	p := &updateFeedbackParams{}
	err = decoder.Decode(p)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding feedback update: %s", err),
		))
	}

	df, serr := services.UpdateDateFeedback(
		db,
		id,
		p.UserID,
		p.Feedback,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateFeedback": df,
	})
}

func updateFeedbackApp(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date ID could not be parsed: %s", err),
		))
	}

	type updateFeedbackParams struct {
		Feedback models.DateFeedbackType
	}
	decoder := json.NewDecoder(r.Body)
	p := &updateFeedbackParams{}
	err = decoder.Decode(p)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding feedback update: %s", err),
		))
	}

	df, serr := services.UpdateDateFeedback(
		db,
		id,
		*sess.UserID,
		p.Feedback,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateFeedback": df,
	})
}

// RegisterDatesControllerRoutes Registers the functions
func RegisterDatesControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/dates", serviceWrapperDBAuthenticated("getDates", getDates, config)).Methods("GET")
	router.HandleFunc("/dates", serviceWrapperDBAuthenticated("createDate", createDate, config)).Methods("POST")
	router.HandleFunc("/dates/{id:[0-9]+}", serviceWrapperDBAuthenticated("getDate", getDate, config)).Methods("GET")
	router.HandleFunc("/dates/{id:[0-9]+}", serviceWrapperDBAuthenticated("updateDate", updateDate, config)).Methods("POST")
	router.HandleFunc("/dates/{id:[0-9]+}/delete", serviceWrapperDBAuthenticated("deleteDate", deleteDate, config)).Methods("POST")
	router.HandleFunc("/dates/{id:[0-9]+}/feedback", serviceWrapperDBAuthenticated("updateFeedback", updateFeedback, config)).Methods("POST")

	router.HandleFunc("/app/v1/dates", serviceWrapperDBAuthenticated("getDates", getAppDates, config)).Methods("GET")
	router.HandleFunc("/app/v1/dates/{id:[0-9]+}", serviceWrapperDBAuthenticated("getDateApp", getDateApp, config)).Methods("GET")
	router.HandleFunc("/app/v1/dates/{id:[0-9]+}/feedback", serviceWrapperDBAuthenticated("updateFeedbackApp", updateFeedbackApp, config)).Methods("POST")
}
