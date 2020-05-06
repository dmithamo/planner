package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLogger *log.Logger
	errLogger  *log.Logger
	port       string
	staticDir  string
	fileServer http.Handler
	mux        *http.ServeMux
}

var app *application

// init *injects dependencies into an instance of the app
func init() {
	port := flag.String("p", ":3001", "http address where server will run")
	staticDir := flag.String("s", "./ui/static", "location of static resources")
	flag.Parse()

	f, err := os.OpenFile("./logs/logs.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("logging to file /logs/log.log")

	fileServer := http.FileServer(http.Dir(*staticDir))
	infoLogger := log.New(f, "INFO \t", log.Ldate|log.Ltime|log.Lshortfile)
	errLogger := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	app = &application{
		infoLogger: infoLogger,
		errLogger:  errLogger,
		port:       *port,
		fileServer: fileServer,
		mux:        mux,
	}
}

// main runs an instance of the app
func main() {
	app.mux.HandleFunc("/", app.home)
	app.mux.HandleFunc("/list/", app.list)
	app.mux.HandleFunc("/list/add/", app.add)
	app.mux.HandleFunc("/settings/", app.settings)
	app.mux.Handle("/static/", http.StripPrefix("/static", app.fileServer))

	srv := &http.Server{
		Addr:     app.port,
		ErrorLog: app.errLogger,
		Handler:  app.mux,
	}

	app.infoLogger.Printf("starting server [127.0.0.1%v]", app.port)
	app.errLogger.Fatal(srv.ListenAndServe())
}
