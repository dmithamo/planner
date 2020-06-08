package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	ilog "github.com/dmithamo/planner/pkg/log" // custom logger
	"github.com/dmithamo/planner/pkg/mysql"
	"github.com/dmithamo/planner/pkg/projects"
	"github.com/justinas/alice"
)

type application struct {
	infoLogger      *log.Logger
	errLogger       *log.Logger
	port            *string
	staticResServer http.Handler
	mux             *http.ServeMux
	templates       map[string]*template.Template
	templateData    templateData
	isDevEnv        bool
}

type templateData struct {
	projects *projects.Projects
}

// make this global to both init() and main() to keep main() short
var app *application

// init *injects dependencies into an instance of the app
func init() {
	port := flag.String("p", ":3001", "http address where server will run")
	dsn := flag.String("db", "", "data source name of the db to be used")
	recreateDB := flag.Bool("rdb", false, "drop db tables and recreate them")
	logToFile := flag.Bool("ltf", true, "toggle whether to persist logs to file")
	isDevEnv := flag.Bool("dev", true, "toggle whether app is running in developemnt mode")
	flag.Parse()

	// instantiate loggers
	logservice := &ilog.Log{}
	if *logToFile {
		year, month, date := time.Now().Date()
		logFile := fmt.Sprintf("./log-%v-%v-%v.log", year, month, date)
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

		checkFatalErrorsHelper(err, logservice.ErrorLogger, "app init::logToFile config::fail")
		logservice = &ilog.Log{
			LogFile: f,
		}
		logservice.Initialize()
	}

	// instantiate db
	dbservice := &mysql.Mysql{}
	db, err := dbservice.OpenDB(fmt.Sprintf("%v?parseTime=true", *dsn))
	checkFatalErrorsHelper(err, logservice.ErrorLogger, "app init::db open::fail")
	logservice.InfoLogger.Println("app init::db open::success")
	dbservice.IDB = db
	// defer db.Close() <- This is not necessary. Or is it? TODO. After reading

	// recreate tables in db if need be
	if *recreateDB {
		err := dbservice.DropTables()
		checkFatalErrorsHelper(err, logservice.ErrorLogger, "app init::db drop tables::fail")
		logservice.InfoLogger.Println("app init::db drop tables::success")

		err = dbservice.CreateTables()
		checkFatalErrorsHelper(err, logservice.ErrorLogger, "app init::db create tables::fail")
		logservice.InfoLogger.Println("app init::db create tables::success")
	}

	// instantiate data: projects model
	projects := projects.Projects{
		IDB: db,
	}

	// instantiate static resources server
	staticDir := "./views/static"
	staticResServer := http.FileServer(http.Dir(staticDir))

	// instantiate template cache
	templateCache, err := buildTemplatesCache("./views/html")
	checkFatalErrorsHelper(err, logservice.ErrorLogger, "app init::templates build cache::fail")
	logservice.InfoLogger.Println("app init::templates build cache::success")

	// Avengers, Assemble! <assemble all the things>
	app = &application{
		port:            port,
		infoLogger:      logservice.InfoLogger,
		errLogger:       logservice.ErrorLogger,
		mux:             http.NewServeMux(),
		staticResServer: staticResServer,
		templateData:    templateData{projects: &projects},
		templates:       templateCache,
		isDevEnv:        *isDevEnv,
	}

	// Yes I know this is superflous
	logservice.InfoLogger.Println("::::::::::::::::::::::::::::::::")
	logservice.InfoLogger.Println("{re}starting application")
	logservice.InfoLogger.Println("::::::::::::::::::::::::::::::::")
	logservice.InfoLogger.Println("app init::success")
}

// main runs an instance of the app
func main() {
	app.mux.HandleFunc("/", app.landingPage)
	app.mux.HandleFunc("/projects/", app.listOfProjects)
	app.mux.HandleFunc("/projects/details/", app.singleProject)

	app.mux.HandleFunc("/settings/", app.settings)
	app.mux.Handle("/static/", http.StripPrefix("/static", app.staticResServer))

	standardMiddleware := alice.New(app.panicRecovery, app.requestLogger, app.auth, app.secureHeaders)
	srv := &http.Server{
		Addr:     *app.port,
		ErrorLog: app.errLogger,
		Handler:  standardMiddleware.Then(app.mux),
	}

	app.infoLogger.Printf("app start::start server [127.0.0.1%v]::success", *app.port)
	app.infoLogger.Println("app start::ready for requests::success")
	app.errLogger.Fatal("app start::fail", srv.ListenAndServe())
}
