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

type application struct {
	infoLogger      *log.Logger
	errLogger       *log.Logger
	port            *string
	staticResServer http.Handler
	mux             *http.ServeMux
	data            dbservice.Data
}

var app *application

// init *injects dependencies into an instance of the app
func init() {
	port := flag.String("p", ":3001", "http address where server will run")
	dsn := flag.String("db", "", "data source name of the db to be used")
	flag.Parse()

	db, err := dbservice.OpenDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logFile := fmt.Sprintf("./logs/log_%v.log", time.Now())
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logservice.AttachLogFile(f)
	log.Printf("logging to file %v", logFile)

	staticDir := "./ui/static"
	staticResServer := http.FileServer(http.Dir(staticDir))

	mux := http.NewServeMux()

	app = &application{
		port:            port,
		infoLogger:      logservice.InforLogger,
		errLogger:       logservice.ErrorLogger,
		mux:             mux,
		staticResServer: staticResServer,
		data:            dbservice.Data{},
	}
}

// main runs an instance of the app
func main() {
	app.mux.HandleFunc("/", app.landingPage)
	app.mux.HandleFunc("/list/", app.listOfTodos)
	app.mux.HandleFunc("/settings/", app.settings)

	app.mux.Handle("/static/", http.StripPrefix("/static", app.staticResServer))

	srv := &http.Server{
		Addr:     *app.port,
		ErrorLog: app.errLogger,
		Handler:  app.mux,
	}

	app.infoLogger.Printf("starting server [127.0.0.1%v]", *app.port)
	app.infoLogger.Println("db awaiting requests")
	app.errLogger.Fatal(srv.ListenAndServe())
}
