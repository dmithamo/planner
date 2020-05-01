package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var infoLogger *log.Logger
var errLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errLogger = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	port := flag.String("p", ":3001", "http address where server will run")
	staticDir := flag.String("s", "./ui/static", "location of static resources")
	flag.Parse()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(*staticDir))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/list/", list)
	mux.HandleFunc("/list/add/", add)
	mux.HandleFunc("/settings/", settings)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLogger.Printf("starting server [127.0.0.1%v]", *port)
	errLogger.Fatal(http.ListenAndServe(*port, mux))
}
