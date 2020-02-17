package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"middleware/handler/db"
	"net/http"

	"github.com/Jeffail/gabs/v2"
)

// UUID to represent logined user
type UUID struct {
	Val string
	// New stands for if this UUID struct is newly assigned by server
	// it's solely used by login function at now
	New bool
}

type handlerWrapper func(w http.ResponseWriter, r *http.Request) http.Handler

func (hW handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hW(w, r).ServeHTTP(w, r)
}

// Err packs error returned by handler to be responded in json
type Err struct {
	err error
}

func (e Err) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := writeJSON(&db.Response{Err: &db.JsError{Err: e.err}, Data: nil}, w)
	if err != nil {
		log.Printf("write http json response: %v\n", err)
		return
	}
	return
}

// JSONData packs meaningful data returned by handler to be responded in json
type JSONData struct {
	d interface{}
}

func (data JSONData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := writeJSON(&db.Response{Err: nil, Data: data.d}, w)
	if err != nil {
		log.Printf("write http json response: %v", err)
		return
	}
	return
}

func writeJSON(resp *db.Response, w http.ResponseWriter) error {
	w.Header().Set(`Content-Type`, `application/json`)

	jo, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal data into json: %v", err)
	}

	_, err = w.Write(jo)
	if err != nil {
		return fmt.Errorf("write http response: %v", err)
	}

	return nil
}

func parseJSONReq(r *http.Request, callback func(pJSON *gabs.Container) error) error {
	if r.Method != http.MethodPost {
		return errors.New("request is not POST")
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("request body is not application/json")
	}

	jsParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		return fmt.Errorf("parse request body as json: %v", err)
	}

	return callback(jsParsed)

}
