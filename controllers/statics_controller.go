package controllers

import (
	"net/http"

	"github.com/gerhardgruber/fame/lib"
	"github.com/gerhardgruber/fame/models"
	"github.com/gerhardgruber/fame/services"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func statics(w http.ResponseWriter, r *http.Request, c *lib.Config) {
	db, err := c.GetDatabaseConnection()
	if err != nil {
		log.Error("Could not get database for rootStatic: ", err)
		replyError(w, "DatabaseError")
	}

	lang := "en"
	if r.FormValue("language") != "" {
		lang = r.FormValue("language")
	} else {
		lang = r.Header.Get("Accept-Language")
	}

	session, err := services.CheckSession(r.FormValue("session"), db)
	if err != nil {
		translationData, lang, err := services.GetTranslationData(lang, db, c)
		if err != nil {
			log.Errorf("Error while loading translation data: %s ", err)
			replyError(w, "TranslationError")
			return
		}
		replyData(w, map[string]interface{}{
			"logged_in":           false,
			"translation_data":    translationData,
			"language":            lang,
			"date_feedback_types": models.DateFeedbackTypes,
		})
		return
	}

	user := session.User

	translationData, language, err := services.GetTranslationData(user.Lang, db, c)
	if err != nil {
		log.Errorf("Error while loading translation data: %s ", err)
		replyError(w, "TranslationError")
		return
	}
	replyData(w, map[string]interface{}{
		"logged_in":           true,
		"login_name":          user.Name,
		"login_time":          session.CreatedAt,
		"translation_data":    translationData,
		"language":            language,
		"UserID":              user.ID,
		"date_feedback_types": models.DateFeedbackTypes,
	})
}

// RegisterStaticsControllerRoutes Registers the functions
func RegisterStaticsControllerRoutes(router *mux.Router, config *lib.Config) {
	router.HandleFunc("/statics", serviceWrapper("Static", statics, config))
}
