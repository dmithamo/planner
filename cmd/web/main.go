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
	f, err := os.OpenFile("./logs/logs.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoLogger = log.New(f, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errLogger = log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
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
