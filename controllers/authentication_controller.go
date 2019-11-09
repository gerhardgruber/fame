package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// TODO: Use transactions here as well

func authenticationLogin(w http.ResponseWriter, r *http.Request, c *lib.Config) {
	r.ParseForm()
	req := struct {
		Name string
		PW   string
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

	loggedInUser, ses, err := services.LoginHandler(c, db, req.Name, req.PW)

	if err != nil {
		log.Error("Error logging in!", err)
		replyError(w, "AuthenticationError")
		return
	}

	replyData(w, map[string]interface{}{
		"user":    loggedInUser,
		"session": ses.Key,
	})
}

func authenticationLoginApp(w http.ResponseWriter, r *http.Request, c *lib.Config) {
	r.ParseForm()
	req := struct {
		Name     string
		PW       string
		Passcode string
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

	loggedInUser, ses, err := services.LoginHandler(c, db, req.Name, req.PW)
	if err != nil {
		log.Error("Error logging in!", err)
		replyError(w, "AuthenticationError")
		return
	}

	// Check if phone exists and create one if not
	mobilePhone, ferr := services.GetMobilePhoneByUserID(db, loggedInUser.ID)
	if ferr != nil {
		log.Error("Error logging in!", ferr)
		replyError(w, "AuthenticationError")
		return
	}
	ph, err := models.HashPassword(req.Passcode)
	if err != nil {
		log.Error("Error hashing phone passcode!", err)
		replyError(w, "AuthenticationError")
		return
	}
	mobilePhone.PasscodeHash = ph
	db.Save(mobilePhone)

	mobilePhone.User = loggedInUser

	replyData(w, map[string]interface{}{
		"MobilePhone": mobilePhone,
		"session":     ses.Key,
	})
}

func authenticationLogout(r *http.Request, params map[string]string, db *gorm.DB, sess *models.Session, c *lib.Config) *reply {
	err := services.CloseSession(sess, db)
	if err != nil {
		return Error(lib.FameError{
			ErrorCode:    "LogoutError",
			Caption:      "ERR_LOGOUT",
			CaptionData:  nil,
			ErrorMessage: fmt.Sprintf("Could not close Session %d: %s", sess.ID, err),
		})
	}

	return Success()
}

func authenticateViaMobilePhone(w http.ResponseWriter, r *http.Request, c *lib.Config) {
	r.ParseForm()

	phone := &models.MobilePhone{}
	err := json.NewDecoder(r.Body).Decode(phone)
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

	ses, phone, err := services.MobilePhoneLoginHandler(c, db, phone.ID, phone.Passcode)
	if err != nil {
		log.Errorf("Error logging in! %+v", err)
		replyError(w, "AuthenticationError")
		return
	}

	replyData(w, map[string]interface{}{
		"session":     ses.Key,
		"MobilePhone": phone,
	})
}

// RegisterAuthenticationControllerRoutes Registers the functions
func RegisterAuthenticationControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/authentication/login", serviceWrapper("AuthenticationLogin", authenticationLogin, config)).Methods("POST")
	router.HandleFunc("/authentication/logout", serviceWrapperDBAuthenticated("AuthenticationLogout", authenticationLogout, config)).Methods("POST")
	router.HandleFunc("/app/v1/authentication/login", serviceWrapper("AuthenticationLoginApp", authenticationLoginApp, config)).Methods("POST")
	router.HandleFunc("/app/v1/authentication/mobile_phone", serviceWrapper("AuthenticationMobilePhone", authenticateViaMobilePhone, config)).Methods("POST")
}
