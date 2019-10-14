package render

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func MustRenderJSONError(w http.ResponseWriter, r *http.Request, err *ErrorResponse) {
	json, e := json.Marshal(err)
	if e != nil {
		log.Error(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	_, e = w.Write(json)
	if e != nil {
		log.Error(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func MustRenderJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
