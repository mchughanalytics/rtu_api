package main

import (
	//"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	myapi "github.com/mchughanalytics/rtu_api/api"

	//mycommon "github.com/mchughanalytics/rtu_api/common"
	myui "github.com/mchughanalytics/rtu_api/ui"
)

func main() {
	// main router
	r := mux.NewRouter()

	// mapping endpoints/functions
	r.HandleFunc("/", myapi.DefaultEndpoint).Methods(http.MethodGet)
	r.HandleFunc("/versions", myapi.Versions).Methods(http.MethodGet)
	r.HandleFunc("/versions/latest", myapi.VersionLatest).Methods(http.MethodGet)

	myui.OpenBrowser("index.html")
	// cli, err := mycommon.NewRmClient("10.11.99.1", "root", "mypassword")
	// if err != nil {
	// 	panic(err)
	// }

	//clijson, _ := json.MarshalIndent(cli, "", "  ")
	//log.Printf("client: \n%s", clijson)

	// start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
