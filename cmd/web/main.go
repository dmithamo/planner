package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dmithamo/todolist/pkg/logservice"
	dbservice "github.com/dmithamo/todolist/pkg/mysqlservice"
)

type data interface{}
type application struct {
	infoLogger      *log.Logger
	errLogger       *log.Logger
	port            *string
	staticResServer http.Handler
	mux             *http.ServeMux
	projects        dbservice.Projects // data: projects
}

// make this global to both init() and main() to keep main() short
var app *application

// init *injects dependencies into an instance of the app
func init() {
	port := flag.String("p", ":3001", "http address where server will run")
	dsn := flag.String("db", "", "data source name of the db to be used")
	flag.Parse()

	logFile := fmt.Sprintf("./logs/log_%v.log", time.Now())
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logservice.AttachLogFile(f)
	log.Printf("logging to file %v", logFile)

	staticDir := "./views/static"
	staticResServer := http.FileServer(http.Dir(staticDir))

	mux := http.NewServeMux()

	// instantiate db, and create tables.
	// db errs are logged from the db service, but using main's loggers
	db, _ := dbservice.OpenDB(*dsn, logservice.InforLogger, logservice.ErrorLogger)
	// defer db.Close() <- This is not necessary. Or is it? TODO
	dbservice.DropTables(db)
	dbservice.CreateTables(db)

	app = &application{
		port:            port,
		infoLogger:      logservice.InforLogger,
		errLogger:       logservice.ErrorLogger,
		mux:             mux,
		staticResServer: staticResServer,
		projects:        dbservice.Projects{DB: db},
	}
}

// main runs an instance of the app
func main() {
	app.mux.HandleFunc("/", app.landingPage)
	app.mux.HandleFunc("/list/", app.listOfProjects)
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
