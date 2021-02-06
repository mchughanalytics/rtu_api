package rtuapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

type RmApi struct {
	BasePath           string
	ListenAddr         string
	Operations         chan *RtuOperationSet
	lookupOperationSet map[string]*RtuOperationSet
	Server             *http.Server
	state              chan *RtuInternalStateChange
	Key                string
}

func NewRmAPI(bp, la, key string) *RmApi {
	api := &RmApi{
		BasePath:   bp,
		ListenAddr: la,
		Key:        key,
	}

	return api
}

// API SETUP

func (a *RmApi) Init(c chan *RtuOperationSet, s chan *RtuInternalStateChange) {
	a.Operations = c
	a.state = s
	a.lookupOperationSet = map[string]*RtuOperationSet{}
	a.Server = &http.Server{
		Addr:    a.ListenAddr,
		Handler: nil,
	}
	a.SetupRoutes()

	if err := a.Server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func (a *RmApi) SetupRoutes() {

	// Define handler functions
	defaultHandler := http.HandlerFunc(a.HandleDefaultEndpoint)
	stopHandler := http.HandlerFunc(a.HandleStop)
	versionsHandler := http.HandlerFunc(a.HandleVersions)
	versionLatestHandler := http.HandlerFunc(a.HandleVersionLatest)
	deviceConfigureHandler := http.HandlerFunc(a.HandleDeviceConfigure)

	// Map handler functions to routes
	http.Handle(fmt.Sprintf("%s/", a.BasePath), Middleware(defaultHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "stop"), Middleware(stopHandler))
	http.Handle(fmt.Sprintf("%s/%s", a.BasePath, "versions"), Middleware(versionsHandler))
	http.Handle(fmt.Sprintf("%s/%s/%s", a.BasePath, "versions", "latest"), Middleware(versionLatestHandler))
	http.Handle(fmt.Sprintf("%s/%s/%s", a.BasePath, "device", "configure"), Middleware(deviceConfigureHandler))
}

// HANDLERS

// DefaultEndpoint returns content for ./
func (a *RmApi) HandleDefaultEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

// Versions returns content for ./versions/
func (a *RmApi) HandleStop(w http.ResponseWriter, r *http.Request) {

	if !a.verifyKey(r, a.Key) {
		invalidKey(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	a.state <- newRtuInternalState("stop")
}

// Versions returns content for ./versions/
func (a *RmApi) HandleVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	v, err := GetAllVersions()
	if err != nil {
		w.Write(jsonErrMessage(err))
	} else {

		vbytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			w.Write(jsonErrMessage(err))
		} else {
			w.Write(vbytes)
		}
	}
}

// Versions returns content for ./versions/
func (a *RmApi) HandleVersionLatest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	v, err := GetAllVersions()
	if err != nil {
		w.Write(jsonErrMessage(err))
	} else {

		latest := v.GetLatest()

		vbytes, err := json.MarshalIndent(latest, "", "  ")
		if err != nil {
			w.Write(jsonErrMessage(err))
		} else {
			w.Write(vbytes)
		}
	}
}

//
func (a *RmApi) HandleDeviceConfigure(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusMethodNotAllowed)
	case http.MethodPost:
		if !verifyContentType(r, "application/json") {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}

		q := &RmDeviceConfiguration{}

		errString := getRequestBody(r, q)
		if len(errString) > 0 {
			switch errString {
			case "default":
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			case "Request body too large":
				http.Error(w, errString, http.StatusRequestEntityTooLarge)
			default:
				http.Error(w, errString, http.StatusBadRequest)
			}
		}

		s := newRtuInternalState("device_configure")
		s.Data = q
		a.state <- s

		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// UTILITY

func jsonErrMessage(msg error) []byte {

	base_start := `{`
	base_end := `}`

	full := fmt.Sprintf("%s\n  \"response\": \"an error occured.\",\n  \"detail\": \"%s\"\n%s", base_start, msg.Error(), base_end)

	return []byte(full)
}

// INTERNAL COMMON HTTP

func interpretHttpError(err error) string {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return fmt.Sprintf("Request body contains badly-formed JSON")
	case errors.As(err, &unmarshalTypeError):
		return fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Sprintf("Request body contains unknown field %s", fieldName)
	case errors.Is(err, io.EOF):
		return "Request body must not be empty"
	case err.Error() == "http: request body too large":
		return "Request body too large"
	default:
		return "default"
	}
}

func getRequestBody(r *http.Request, output interface{}) string {

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&output)
	if err != nil {
		return interpretHttpError(err)
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return "Request body must only contain a single JSON object"
	}

	return ""
}

func verifyContentType(r *http.Request, ct string) bool {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != ct {
			return false
		}
		return true
	}
	return false
}

func (a *RmApi) findOperation(ID string) []byte {

	for i, opset := range a.lookupOperationSet {
		if i == ID {
			return opset.GetOutputJSON()
		}

		for oi, op := range opset.lookupOperation {
			if oi == ID {
				return op.GetOutputJSON()
			}
		}
	}

	return nil
}

func (a *RmApi) verifyKey(r *http.Request, key string) bool {
	if r.Header.Get("Request-ID") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Request-ID")
		if value != key {
			return false
		}
		return true
	}
	return false
}

func invalidKey(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(jsonErrMessage(fmt.Errorf("Failed to provide valid Request-ID")))
	return
}

// CUSTOM MIDDLEWARE

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler.ServeHTTP(w, r)
	})
}
