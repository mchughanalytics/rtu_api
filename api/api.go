package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mchughanalytics/rtu_api/common"
)

// DefaultEndpoint returns content for ./
func DefaultEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

// Versions returns content for ./versions/
func Versions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	v, err := common.GetAllVersions()

	fmt.Printf("\n\nv: %s\nerr: %s", v, err)

	vbytes, err := json.MarshalIndent(v, "", "  ")

	fmt.Printf("\n\nvbytes: %s", vbytes)

	if err != nil {
		w.Write([]byte(`{"message": "an error occured."}`))
	} else {
		w.Write(vbytes)
	}
}
