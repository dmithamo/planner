package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	const port = ":3001"
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Printf("starting server at: http://127.0.0.1%v", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)

	rd, err := json.Marshal(r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(rd))
}
