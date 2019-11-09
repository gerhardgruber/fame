package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// serviceFunction is a function that can be converted to a wrappedServiceFunction by serviceWrapper
type serviceFunction func(http.ResponseWriter, *http.Request, *lib.Config)

// serviceFunctionDBAuthenticated is a function that can converted to a wrappedServiceFunction by serviceWrapperDBAuthenticated
type serviceFunctionDBAuthenticated func(*http.Request, map[string]string, *gorm.DB, *models.Session, *lib.Config) *reply

// wrappedServiceFunction is a function that is ready to be passed to mux.Router.HandleFunc
type wrappedServiceFunction func(http.ResponseWriter, *http.Request)

// serviceWrapper adds the current configuration to the web request
func serviceWrapper(name string, fun serviceFunction, conf *lib.Config) wrappedServiceFunction {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fun(w, r, conf)

	}
}

// serviceWrapperDBAuthenticated adds the current configuration and a database connection to the web requests
// and checks if both the current session is correct and the user has the necessary priviledges to do the request
func serviceWrapperDBAuthenticated(name string, fun serviceFunctionDBAuthenticated, conf *lib.Config) wrappedServiceFunction {
	tmpFunc := func(w http.ResponseWriter, r *http.Request, c *lib.Config) {
		db, err := c.GetDatabaseConnection()
		if err != nil {
			log.Error(fmt.Sprintf("Could not get database for %s: ", name), err)
			replyError(w, "DatabaseError")
			return
		}

		tx := db.Begin()
		if tx.Error != nil {
			replyError(w, "DBBeginError")
			return
		}

		var repl *reply

		defer func() {
			panic := recover()
			if panic != nil {
				log.Errorf("PANIC: %s", panic)
				debug.PrintStack()
			}

			if repl == nil || !repl.Success {
				w.WriteHeader(500)

				rxrb := tx.Rollback()
				if rxrb.Error != nil {
					replyError(w, "DBRollbackError")
					return
				}

				if repl == nil {
					repl = &reply{false, "ServerError", nil, "", nil, nil}
				} else {
					repl.Message = name + ": " + repl.Message
				}
			} else {
				rxcm := tx.Commit()
				if rxcm.Error != nil {
					replyError(w, "DBCommitError")
					return
				}
			}

			sendReply(w, repl)
		}()

		vars := mux.Vars(r)

		session, err := services.CheckSession(r.FormValue("session"), tx)
		if err != nil {
			log.Warn(fmt.Sprintf("Authentication Error: Could not authenticate for %s: ", name), err)
			replyError(w, "AuthenticationError")
			return
		}

		repl = fun(r, vars, tx, session, c)
	}

	return serviceWrapper(name, tmpFunc, conf)
}

type file struct {
	filename string
	mimeType string
	data     []byte
}

type reply struct {
	Success     bool                   `json:"success"`
	Caption     string                 `json:"caption,omitempty"`
	CaptionData map[string]interface{} `json:"captionData,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	File        *file                  `json:"-"`
}

type REPLY = reply

// replyError writes success: false as well as a custom error message to the output stream
func replyError(w http.ResponseWriter, message string) {
	w.WriteHeader(500)
	sendReply(w, Error(*lib.InternalError(
		errors.New(message),
	)))
}

// replyFameError writes success: false as well as a custom error message to the output stream
func replyFameError(w http.ResponseWriter, err lib.FameError) {
	w.WriteHeader(500)
	sendReply(w, Error(err))
}

// replyData writes success: true as well as data to the output stream
func replyData(w http.ResponseWriter, data map[string]interface{}) {
	sendReply(w, Success(data))
}

// Error fills the error message into a reply struct representing a failed API call
func Error(e lib.FameError) *reply {
	log.Errorf("%s: %s", e.ErrorCode, e.ErrorMessage)
	if len(e.StackTrace) != 0 {
		log.Error(string(e.StackTrace))
	}
	return &reply{
		Success:     false,
		Caption:     e.Caption,
		CaptionData: e.CaptionData,
		Message:     e.ErrorMessage,
		Data:        nil,
		File:        nil,
	}
}

// Success fills the data into a reply struct representing a successfull API call
func Success(data ...map[string]interface{}) *reply {
	if len(data) == 0 {
		return &reply{true, "", nil, "", map[string]interface{}{}, nil}
	} else if len(data) > 1 {
		log.Errorf("Trying to write %d results on success, can only take one!", len(data))
	}

	return &reply{true, "", nil, "", data[0], nil}
}

func File(filename string, mimeType string, data []byte) *reply {
	return &reply{true, "", nil, "", nil, &file{
		filename, mimeType, data,
	}}
}

// sendReply should not be used outside utils
func sendReply(w http.ResponseWriter, r *reply) {
	if r.File != nil {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", r.File.filename))
		w.Header().Set("Content-Type", r.File.mimeType)
		n, err := w.Write(r.File.data)
		if err != nil {
			log.Errorf("Error while sending result to client after %d bytes: %s", n, err)
		}
	} else {
		out, err := json.Marshal(r)
		if err != nil {
			log.Errorf("Error while encoding reply %+v to json! %s", r, err)
			out, err = json.Marshal(&reply{false, "ERR_INTERNAL", nil, "MarshalDataError", nil, nil})
			if err != nil {
				log.Fatalf("Can not marshal fallback reply data")
			}
		}

		// log.Infof("data: %s", string(out))
		n, err := fmt.Fprintf(w, string(out))
		if err != nil {
			log.Errorf("Error while sending result to client after %d bytes: %s", n, err)
		}
	}
}
