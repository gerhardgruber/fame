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

func getOperations(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	operations, serr := services.GetOperations(db)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"Operations": operations,
	})
}

func createOperation(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	type createOperationParams struct {
		Title     string
		FirstName string
		LastName  string
	}

	decoder := json.NewDecoder(r.Body)
	cup := createOperationParams{}
	err := decoder.Decode(&cup)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding Operation: %s", err),
		))
	}

	o, serr := services.CreateOperation(
		cup.Title,
		cup.FirstName,
		cup.LastName,
		db,
	)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"Operation": o,
	})
}

func getOperation(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	operationID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Operation ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	operation, serr := services.GetOperationByID(operationID, db)
	if serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"Operation": operation,
	})
}

func deleteOperation(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Operation ID could not be parsed: ", err),
		))
	}

	decoder := json.NewDecoder(r.Body)
	data := models.Operation{}
	err = decoder.Decode(&data)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("InvalidObjectError %+v", err),
		))
	}

	operation, serr := services.GetOperationByID(id, db)
	if serr != nil {
		return Error(*serr)
	}

	err = db.Delete(operation).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("DataBaseError %+v", err),
		))
	}

	return Success()
}

func updateOperationAPI(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Operation ID could not be parsed: ", err),
		))
	}

	decoder := json.NewDecoder(r.Body)
	data := models.Operation{}
	err = decoder.Decode(&data)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("InvalidObjectError %+v", err),
		))
	}

	operation, serr := services.GetOperationByID(id, db)
	if serr != nil {
		return Error(*serr)
	}

	operation.Title = data.Title
	operation.FirstName = data.FirstName
	operation.LastName = data.LastName

	err = db.Save(operation).Error
	if err != nil {
		return Error(*lib.DataCorruptionError(
			fmt.Errorf("DataBaseError %+v", err),
		))
	}

	return Success(map[string]interface{}{
		"operation": operation,
	})
}

// RegisterOperationsControllerRoutes Registers the functions
func RegisterOperationsControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/operations", serviceWrapperDBAuthenticated("getOperations", getOperations, config)).Methods("GET")
	router.HandleFunc("/operations", serviceWrapperDBAuthenticated("createOperation", createOperation, config)).Methods("POST")
	router.HandleFunc("/operations/{id:[0-9]+}", serviceWrapperDBAuthenticated("getOperation", getOperation, config)).Methods("GET")
	router.HandleFunc("/operations/{id:[0-9]+}", serviceWrapperDBAuthenticated("updateOperation", updateOperationAPI, config)).Methods("POST")
	router.HandleFunc("/operations/{id:[0-9]+}/delete", serviceWrapperDBAuthenticated("deleteOperation", deleteOperation, config)).Methods("POST")
}
