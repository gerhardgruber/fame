package lib

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jinzhu/gorm"
	// Loading all different dialects for gorm
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")

	executable, err := os.Executable()
	if err != nil {
		log.Error("An error occured while reading the location of the executable! ", err)
	}

	configDir := filepath.Join(filepath.Dir(executable), "..", "config")
	viper.AddConfigPath(configDir)
	viper.SetConfigName("fame")
	log.Debug("Added config path ", configDir)
}

// Config contains config information, e.g. the database connection string
type Config struct {
	DatabaseType             string
	DatabaseConnectionString string
	databaseConnection       *gorm.DB
	GoogleAPIKey             string
	languageHander           *LanguageHandler
	I18nBasePath             string
	MailUser                 string
	MailPassword             string
	WebserverPort            int
}

// NewConfig creates a new Config struct
// If fileName is empty, the function will search for a config file with the name "config" (and extension .json, .yaml, ...)
// in a "config" directory parallel to the current executable directory.
func NewConfig(fileName string) (*Config, error) {
	if fileName != "" {
		viper.SetConfigFile(fileName)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not read config file %s! %s", viper.ConfigFileUsed(), err)
	}

	log.Info("Using config file ", viper.ConfigFileUsed())

	dBType := viper.GetString("database.type")
	dBConnection := viper.GetString("database.connection")
	if dBType == "" {
		return nil, fmt.Errorf("Error while reading config file %s! No database.type specified", viper.ConfigFileUsed())
	}
	if dBConnection == "" {
		return nil, fmt.Errorf("Error while reading config file %s! No database.connection specified", viper.ConfigFileUsed())
	}

	i18nBasePath := viper.GetString("path.i18n")
	if i18nBasePath == "" {
		i18nBasePath = "i18n"
	}

	i18nBasePath, err = filepath.Abs(i18nBasePath)
	if err != nil {
		return nil, fmt.Errorf("Error while reading config file %s! Could not read i18n base path! %s", viper.ConfigFileUsed(), err)
	}

	MailUser := viper.GetString("mail.user")
	if MailUser == "" {
		MailUser = "service@fame.com"
	}

	MailPassword := viper.GetString("mail.password")
	if MailPassword == "" {
		MailPassword = "1234"
	}

	WebserverPort := viper.GetString("webserver.port")
	if WebserverPort == "" {
		WebserverPort = "9000"
	}

	port, err := strconv.Atoi(WebserverPort)
	if err != nil {
		port = 9000
	}

	GoogleAPIKey := viper.GetString("google.apikey")

	return &Config{DatabaseType: dBType, DatabaseConnectionString: dBConnection, GoogleAPIKey: GoogleAPIKey, I18nBasePath: i18nBasePath, MailUser: MailUser, MailPassword: MailPassword, WebserverPort: port}, nil
}

// SetDefaultArguments sets the default programm arguments and returns pointers to the set arguments.
// flag.Parse() must becalled before the arguments can be accessed.
func SetDefaultArguments() (*string, *string) {
	configFile := flag.String("config", "", "The config file")
	logLevel := flag.String("loglevel", log.InfoLevel.String(), "The log level")

	return configFile, logLevel
}

// SetLogLevel sets the log level with the given string.
// If the string could not be parsed, log.InfoLevel will be set.
func SetLogLevel(logLevel *string) {
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Errorf("Unknown error level %s was set. Falling back to level Info.", *logLevel)
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

// GetDatabaseConnection opens a new database connection,
// or returns a already opened database connection
func (c *Config) GetDatabaseConnection() (*gorm.DB, *FameError) {
	if c.databaseConnection != nil {
		return c.databaseConnection, nil
	}

	log.Debug("Opening new database connection...")

	db, err := gorm.Open(c.DatabaseType, c.DatabaseConnectionString)
	if err != nil {
		return nil, InternalError(
			fmt.Errorf("Could not connect to database with type \"%s\" and connection string \"%s\"! %s", c.DatabaseType, c.DatabaseConnectionString, err),
		)
	}
	c.databaseConnection = db //.Debug()
	return c.databaseConnection, nil
}

// GetGoogleAPIKey returns the GoogleAPIKey
func (c *Config) GetGoogleAPIKey() string {
	return c.GoogleAPIKey
}

// GetLanguageHandler loads the language handler
func (c *Config) GetLanguageHandler() (LanguageHandler, *FameError) {
	if c.languageHander != nil {
		return *c.languageHander, nil
	}

	return LanguageHandler{c.I18nBasePath, make(map[string]Language)}, nil
}
