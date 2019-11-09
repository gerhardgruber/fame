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

func getDateCategories(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dateCategories, serr := services.GetDateCategories(db)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateCategories": dateCategories,
	})
}

func getAppDateCategories(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dateCategories, serr := services.GetDateCategories(db)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateCategories": dateCategories,
	})
}

func getDateCategory(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	dateCategoryID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date category ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	dateCategory, serr := services.GetDateCategoryByID(db, dateCategoryID)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateCategory": dateCategory,
	})
}

func createDateCategory(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	decoder := json.NewDecoder(r.Body)
	p := &models.DateCategory{}
	err := decoder.Decode(p)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding date category: %s", err),
		))
	}

	dc, serr := services.CreateDateCategory(
		c,
		db,
		sess.User,
		p,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateCategory": dc,
	})
}

func deleteDateCategory(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date Category ID could not be parsed: %s", err),
		))
	}

	dateCategory, serr := services.GetDateCategoryByID(db, id)
	if serr != nil {
		return Error(*serr)
	}

	err = db.Delete(dateCategory).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("DataBaseError %+v", err),
		))
	}

	return Success()
}

func updateDateCategory(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Date Category ID could not be parsed: %s", err),
		))
	}

	decoder := json.NewDecoder(r.Body)
	p := &models.DateCategory{}
	err = decoder.Decode(p)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding date category: %s", err),
		))
	}

	dateCategory, serr := services.UpdateDateCategory(
		c,
		db,
		id,
		p,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"DateCategory": dateCategory,
	})
}

// RegisterDateCategoriesControllerRoutes Registers the functions
func RegisterDateCategoriesControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/date_categories", serviceWrapperDBAuthenticated("getDateCategories", getDateCategories, config)).Methods("GET")
	router.HandleFunc("/date_categories", serviceWrapperDBAuthenticated("createDateCategory", createDateCategory, config)).Methods("POST")
	router.HandleFunc("/date_categories/{id:[0-9]+}", serviceWrapperDBAuthenticated("getDateCategory", getDateCategory, config)).Methods("GET")
	router.HandleFunc("/date_categories/{id:[0-9]+}", serviceWrapperDBAuthenticated("updateDateCategory", updateDateCategory, config)).Methods("POST")
	router.HandleFunc("/date_categories/{id:[0-9]+}/delete", serviceWrapperDBAuthenticated("deleteDateCategory", deleteDateCategory, config)).Methods("POST")

	router.HandleFunc("/app/v1/date_categories", serviceWrapperDBAuthenticated("getDateCategories", getAppDateCategories, config)).Methods("GET")
}
