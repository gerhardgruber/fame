package services

import (
	"strings"

	"github.com/gerhardgruber/fame/lib"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func GetTranslationData(lang string, db *gorm.DB, c *lib.Config) (map[string]string, *string, *lib.FameError) {
	if lang == "" {
		lang = "de"
	}

	languageHandler, err := c.GetLanguageHandler()
	if err != nil {
		return nil, nil, err
	}

	var captions map[string]string

	for _, l := range strings.Split(lang, ",") {
		lang := strings.Split(l, ";")[0]
		captions, err = languageHandler.GetLanguage(lang)
		if err == nil {
			lang = strings.Split(lang, "-")[0]
			return captions, &lang, nil
		} else if strings.Split(err.String(), ":")[0] != "InvalidLanguageError" {
			log.Warn(err)
		}
	}

	lang = "de"
	captions, err = languageHandler.GetLanguage(lang)
	if err == nil {
		lang = strings.Split(lang, "-")[0]
		return captions, &lang, nil
	} else if strings.Split(err.String(), ":")[0] != "InvalidLanguageError" {
		log.Warn(err)
	}

	return nil, nil, err
}
