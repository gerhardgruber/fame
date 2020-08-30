package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/NYTimes/gziphandler"
	logrusmiddleware "github.com/bakins/logrus-middleware"
	"github.com/kardianos/service"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/gerhardgruber/fame/controllers"
	"github.com/gerhardgruber/fame/lib"
	"github.com/gorilla/mux"
)

// ProgramName contains the name of the program
// Will be output when starting the program
const ProgramName = "fame_server"

type serviceFunction func(http.ResponseWriter, *http.Request, *lib.Config)
type wrappedServiceFunction func(http.ResponseWriter, *http.Request)

func serviceWrapper(fun serviceFunction, conf *lib.Config) wrappedServiceFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		fun(w, r, conf)
	}
}

const port int = 80002

type server struct {
	configFile string
	config     *lib.Config
}

// Start starts the fame_server system service
func (s *server) Start(srv service.Service) error {
	// Start should not block. Do the actual work async.
	log.Info("Starting...")
	go s.run()
	return nil
}

// Stop stops the fame_server system service
func (s *server) Stop(srv service.Service) error {
	// Stop should not block. Return with a few seconds.
	log.Info("Stopping...")

	return nil
}

func serveFile(webserverDir string, file string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := filepath.Join(webserverDir, "webapp", file)

		if _, err := os.Stat(fileName); err != nil {
			if os.IsNotExist(err) {
				log.Errorf("Error serving file %s! %s", fileName, err)
			}
		}

		http.ServeFile(w, r, fileName)
	})
}

// run starts the webserver
func (s *server) run() {
	config, err := lib.NewConfig(s.configFile)
	if err != nil {
		log.Error("An error occured while creating the config! ", err)
		os.Exit(2)
	}
	s.config = config

	db, serr := config.GetDatabaseConnection()
	if serr != nil {
		log.Error("An error occured while connecting to database!", serr)
		os.Exit(2)
	}
	log.Info("Migrating the database...")
	lib.MigrateDatabase(db)

	webserverDir, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), ".."))
	log.Infof("Webserver dir: %s", webserverDir)

	log.Info("Started webserver on port ", config.WebserverPort)
	log.Infof("Serving static files from %s", webserverDir)

	router := mux.NewRouter()
	webappRouter := router.PathPrefix("/api/").Subrouter()

	controllers.RegisterAuthenticationControllerRoutes(webappRouter, config)
	controllers.RegisterDateCategoriesControllerRoutes(webappRouter, config)
	controllers.RegisterDatesControllerRoutes(webappRouter, config)
	controllers.RegisterDateLogsControllerRoutes(webappRouter, config)
	controllers.RegisterMobilePhonesControllerRoutes(webappRouter, config)
	controllers.RegisterOperationsControllerRoutes(webappRouter, config)
	controllers.RegisterPositionControllerRoutes(webappRouter, config)
	controllers.RegisterStaticsControllerRoutes(webappRouter, config)
	controllers.RegisterUsersControllerRoutes(webappRouter, config)
	//router.PathPrefix("/dist").Handler(http.StripPrefix("/dist", http.FileServer(http.Dir(filepath.Join(webserverDir, "dist")))))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(filepath.Join(webserverDir, "webapp")))))
	router.Handle("/index.html", serveFile(webserverDir, "index.html"))
	router.Handle("/", serveFile(webserverDir, "index.html"))
	router.NotFoundHandler = serveFile(webserverDir, "index.html")

	log.Debug("Registered routes:")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		log.Debugf("Route: %s", path)
		return nil
	})

	handler := cors.Default().Handler(gziphandler.GzipHandler(router))
	if log.GetLevel() == log.DebugLevel {
		l := logrusmiddleware.Middleware{
			Name:   "fame",
			Logger: log.StandardLogger(),
		}

		handler = l.Handler(handler, "App")
	}

	err = http.ListenAndServe(":"+strconv.Itoa(config.WebserverPort), handler)

	if err != nil {
		log.Errorf("Error while starting fame server! %s", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Printf("Starting %s\n", ProgramName)
	configFile, logLevel := lib.SetDefaultArguments()
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()
	lib.SetLogLevel(logLevel)

	server := &server{configFile: *configFile}

	svcConfig := &service.Config{
		Name:        "fame_server",
		DisplayName: "fame server",
		Description: "fame server by Gerhard Gruber",
	}

	if len(*svcFlag) != 0 && *svcFlag == "install" {
		username := lib.ReadValue("username the service should run as")
		svcConfig.UserName = username
	}

	s, err := service.New(server, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
