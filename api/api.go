package api

import (
	"encoding/json"
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
	if err != nil {
		w.Write([]byte(`{"message": "an error occured."}`))
	} else {

		vbytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			w.Write([]byte(`{"message": "an error occured."}`))
		} else {
			w.Write(vbytes)
		}
	}
}

// Versions returns content for ./versions/
func VersionLatest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	v, err := common.GetAllVersions()
	if err != nil {
		w.Write([]byte(`{"message": "an error occured."}`))
	} else {

		latest := v.GetLatest()

		vbytes, err := json.MarshalIndent(latest, "", "  ")
		if err != nil {
			w.Write([]byte(`{"message": "an error occured."}`))
		} else {
			w.Write(vbytes)
		}
	}
}
