package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	myapi "github.com/mchughanalytics/rtu_api/api"
)

func main() {
	// main router
	r := mux.NewRouter()

	// mapping endpoints/functions
	r.HandleFunc("/", myapi.DefaultEndpoint).Methods(http.MethodGet)
	r.HandleFunc("/", myapi.Versions).Methods(http.MethodGet)

	// start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
