package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// weiteResponseJSON sends a response with status and data
func writeResponseJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Can't marshal to JSON, will send response as text/html. %s", err.Error())
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = fmt.Fprint(w, "Failed build response")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, err = fmt.Fprint(w, string(jsonData))
	}

	if err != nil {
		log.Error(err)
	}
}

// writeResponseEmpty sends a response only with status
func writeResponseEmpty(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}
