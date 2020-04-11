package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	myapi "github.com/mchughanalytics/rtu_api/api"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", myapi.DefaultEndpoint).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}
