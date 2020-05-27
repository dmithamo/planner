package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ilog "github.com/dmithamo/planner/pkg/log" // custom logger
	"github.com/dmithamo/planner/pkg/mysql"
	"github.com/dmithamo/planner/pkg/projects"
)

type application struct {
	infoLogger      *log.Logger
	errLogger       *log.Logger
	port            *string
	staticResServer http.Handler
	mux             *http.ServeMux
	projects        *projects.Projects
}

// make this global to both init() and main() to keep main() short
var app *application

// init *injects dependencies into an instance of the app
func init() {
	port := flag.String("p", ":3001", "http address where server will run")
	dsn := flag.String("db", "", "data source name of the db to be used")
	flag.Parse()

	// instantiate loggers
	logFile := fmt.Sprintf("./logs/log_%v.log", time.Now())
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logservice := &ilog.Log{
		LogFile: f,
	}
	logservice.AttachLogFile()

	dbservice := &mysql.Mysql{}
	db, _ := dbservice.OpenDB(*dsn)
	dbservice.IDB = db
	// defer db.Close() <- This is not necessary. Or is it? TODO
	dbservice.DropTables()
	dbservice.CreateTables()

	// instantiate projects model
	projects := projects.Projects{
		IDB:         db,
		InfoLogger:  logservice.InfoLogger,
		ErrorLogger: logservice.ErrorLogger,
	}

	staticDir := "./views/static"
	staticResServer := http.FileServer(http.Dir(staticDir))

	mux := http.NewServeMux()

	app = &application{
		port:            port,
		infoLogger:      logservice.InfoLogger,
		errLogger:       logservice.ErrorLogger,
		mux:             mux,
		staticResServer: staticResServer,
		projects:        &projects,
	}
}

// main runs an instance of the app
func main() {
	app.mux.HandleFunc("/", app.landingPage)
	app.mux.HandleFunc("/projects/", app.listOfProjects)
	app.mux.HandleFunc("/projects/details/", app.singleProject)

	app.mux.HandleFunc("/settings/", app.settings)
	app.mux.Handle("/static/", http.StripPrefix("/static", app.staticResServer))

	srv := &http.Server{
		Addr:     *app.port,
		ErrorLog: app.errLogger,
		Handler:  app.mux,
	}

	app.infoLogger.Printf("starting server [127.0.0.1%v]", *app.port)
	app.infoLogger.Println("awaiting requests")
	app.errLogger.Fatal(srv.ListenAndServe())
}
