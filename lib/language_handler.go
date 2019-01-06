package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// pathCleaner ensures the name is eg `en-US` or `de`
	pathCleaner = regexp.MustCompile(`^[a-z]{2}(-[a-zA-Z]{2})?$`)
)

// Language is a collection of caption mappings
type Language struct {
	Captions map[string]string
	LastUsed time.Time
}

// LanguageHandler holds all languages loaded from the I18nBasePath
type LanguageHandler struct {
	I18nBasePath string
	languages    map[string]Language
}

// GetLanguage loads a Lanugage from the LanguageHandler
func (l *LanguageHandler) GetLanguage(name string) (map[string]string, *FameError) {
	if !pathCleaner.MatchString(name) {
		return nil, InternalError(
			fmt.Errorf("InvalidLanguageError:Invalid language name %s", name),
		)
	}

	language := Language{}
	if lang, ok := l.languages[name]; ok {
		language = lang
	} else if lang, ok := l.languages[name[:2]]; ok {
		language = lang
	}

	if language.Captions != nil {
		return language.Captions, nil
	}

	lang, langName, err := l.loadLanguage(name)
	if err != nil {
		return nil, err
	}

	language.Captions = lang
	l.languages[langName] = language

	return lang, nil
}

func (l *LanguageHandler) loadLanguage(name string) (map[string]string, string, *FameError) {
	var langName string

	name = strings.ToLower(name)

	i18nFile := filepath.Join(l.I18nBasePath, name+".json")
	log.Info("Try to load i18n file from ", i18nFile)
	if _, err := os.Stat(i18nFile); !os.IsNotExist(err) {
		langName = name
	} else if _, err := os.Stat(filepath.Join(l.I18nBasePath, name[:2]+".json")); !os.IsNotExist(err) {
		langName = name[:2]
	} else {
		return nil, "", InternalError(
			fmt.Errorf("InvalidLanguageError:Language '%s' not found", name),
		)
	}

	file, err := os.Open(filepath.Join(l.I18nBasePath, langName+".json"))
	if err != nil {
		return nil, "", InternalError(err)
	}

	decoder := json.NewDecoder(file)

	language := make(map[string]string)
	err = decoder.Decode(&language)
	if err != nil {
		return nil, "", InternalError(err)
	}

	return language, langName, nil
}
