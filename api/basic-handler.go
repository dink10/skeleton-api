package api

import (
	"net/http"
)

func (srv *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	writeResponseEmpty(w, r, http.StatusOK)
}
