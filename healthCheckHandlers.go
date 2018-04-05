package main

import (
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeResponseEmpty(w, r, http.StatusOK)
}
