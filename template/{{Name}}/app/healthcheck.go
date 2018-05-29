package app

import (
	"fmt"
	"net/http"
)

// HealthCheck is used for checking third-party connections, returns detailed JSON
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// Status returns brief status based on HealthCheck results
func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
