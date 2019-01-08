package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func createMobilePhone(r *http.Request, params map[string]string, db *gorm.DB, session *models.Session, c *lib.Config) *reply {
	phoneParams := struct {
		PhoneNumber string
		Device      string
		DeviceType  models.MobileDeviceType
		Passcode    string
		UserID      uint64
	}{}
	err := json.NewDecoder(r.Body).Decode(&phoneParams)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding MobilePhone: %s", err),
		))
	}

	phone := models.MobilePhone{
		PhoneNumber: phoneParams.PhoneNumber,
		Device:      phoneParams.Device,
		DeviceType:  phoneParams.DeviceType,
		Passcode:    phoneParams.Passcode,
		UserID:      phoneParams.UserID,
	}

	// TODO: broken, checks only for userID which is set above
	allowed, serr := services.MayUserSeePhone(session.User, &phone, db)
	if serr != nil {
		return Error(*lib.PrivilegeError(
			fmt.Errorf("Error checking privileges while creating MobilePhone: %s", serr),
		))
	}
	if !allowed {
		return Error(*lib.PrivilegeError(
			errors.New("User trying to create MobilePhone but has privilege error"),
		))
	}

	if serr := services.CreateMobilePhone(&phone, session.User, db); serr != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"MobilePhone": phone,
	})
}

func updateMobilePhone(r *http.Request, params map[string]string, db *gorm.DB, session *models.Session, c *lib.Config) *reply {
	phoneID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("MobilePhone ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	phone, serr := services.GetMobilePhone(phoneID, db)
	if serr != nil {
		return Error(*serr)
	}

	allowed, serr := services.MayUserSeePhone(session.User, phone, db)
	if serr != nil {
		return Error(*lib.PrivilegeError(
			fmt.Errorf("Error checking privileges while updating MobilePhone: %s", serr),
		))
	}
	if !allowed {
		return Error(*lib.PrivilegeError(
			errors.New("User trying to update MobilePhone but has privilege error"),
		))
	}

	req := struct {
		Device          string
		AllowedTruckIDs []uint64
	}{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("Invalid object while decoding MobilePhone: %s", err),
		))
	}

	phone.Device = req.Device

	if serr := services.UpdatePhone(phone, db); err != nil {
		return Error(*serr)
	}

	return Success(map[string]interface{}{
		"MobilePhone": phone,
	})
}

func createInitialMobilePhone(w http.ResponseWriter, r *http.Request, c *lib.Config) {
	req := struct {
		PhoneNumber string
		EMail       string
		Device      string
		Passcode    string
		DeviceType  models.MobileDeviceType
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("InvalidObjectError %+v", err)
		replyError(w, "InvalidObjectError")
		return
	}

	db, serr := c.GetDatabaseConnection()
	if serr != nil {
		log.Error("Could not get database for login", serr)
		replyError(w, "DatabaseError")
		return
	}

	phone := models.MobilePhone{
		PhoneNumber: req.PhoneNumber,
		Device:      req.Device,
		DeviceType:  req.DeviceType,
		Passcode:    req.Passcode,
	}

	email := req.EMail
	if email == "" {
		email = req.PhoneNumber
	}
	user, serr := services.GetOrCreateUserByName(email, db)
	if serr != nil {
		log.Errorf("Could not get or create user with email '%s': %s", req.EMail, serr)
		replyError(w, "UserError")
		return
	}

	serr = services.CreateMobilePhone(&phone, user, db)
	if serr != nil {
		log.Errorf("Could not create Phone %+v", serr)
		replyError(w, "CreateError")
		return
	}
	phone.User.MobilePhone = nil

	replyData(w, map[string]interface{}{
		"MobilePhone": phone,
	})
}

func getMobilePhone(r *http.Request, params map[string]string, db *gorm.DB, session *models.Session, c *lib.Config) *reply {
	phoneID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return Error(*lib.InvalidParamsError(
			fmt.Errorf("MobilePhone ID '%s' could not be parsed: %s", params["id"], err),
		))
	}

	mobilePhone, serr := services.GetMobilePhone(phoneID, db)
	if serr != nil {
		return Error(*serr)
	}

	allowed, serr := services.MayUserSeePhone(session.User, mobilePhone, db)
	if serr != nil {
		return Error(*lib.PrivilegeError(
			fmt.Errorf("Error checking privileges while getting MobilePhone: %s", serr),
		))
	}
	if !allowed {
		return Error(*lib.PrivilegeError(
			errors.New("User trying to get MobilePhone but has privilege error"),
		))
	}

	return Success(map[string]interface{}{
		"MobilePhone": mobilePhone,
	})
}

func logFromPhone(w http.ResponseWriter, r *http.Request, c *lib.Config) {
	db, serr := c.GetDatabaseConnection()
	if serr != nil {
		log.Errorf("Could not connect to database: %s", serr.String())
		replyFameError(w, *serr)
		return
	}

	params := mux.Vars(r)

	phoneID, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		log.Errorf("Invalid params! %s", err)
		replyFameError(w, *lib.InvalidParamsError(
			fmt.Errorf("MobilePhone ID '%s' could not be parsed: %s", params["id"], err),
		))
		return
	}

	req := struct {
		Message string
		Context string
		PhoneID string
	}{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("Invalid JSON data: %s", err)
		replyFameError(w, *lib.InvalidParamsError(err))
		return
	}

	mpl := &models.MobilePhoneLog{
		MobilePhoneID: phoneID,
		Message:       req.Message,
		Context:       req.PhoneID + ": " + req.Context,
	}

	err = db.Save(mpl).Error
	if err != nil {
		log.Errorf("Error saving log: %s", err)
		replyFameError(w, *lib.DataCorruptionError(err))
		return
	}

	replyData(w, map[string]interface{}{
		"ID": mpl.ID,
	})
}

// RegisterMobilePhonesControllerRoutes registers the functions
func RegisterMobilePhonesControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/mobile_phones", serviceWrapperDBAuthenticated("CreateMobilePhone", createMobilePhone, config)).Methods("POST")
	router.HandleFunc("/mobile_phones/{id:[0-9]+}", serviceWrapperDBAuthenticated("GetMobilePhone", getMobilePhone, config)).Methods("GET")
	router.HandleFunc("/mobile_phones/{id:[0-9]+}", serviceWrapperDBAuthenticated("updateMobilePhone", updateMobilePhone, config)).Methods("POST")
	router.HandleFunc("/app/v1/mobile_phones", serviceWrapper("app/CreateMobilePhone", createInitialMobilePhone, config)).Methods("POST")
	router.HandleFunc("/app/v1/mobile_phones/{id:[0-9]+}", serviceWrapperDBAuthenticated("app/GetMobilePhone", getMobilePhone, config)).Methods("GET")
	router.HandleFunc("/app/v1/mobile_phones/{id:[0-9]+}/log", serviceWrapper("logFromPhone", logFromPhone, config)).Methods("POST")
}
