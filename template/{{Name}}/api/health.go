package api

import "net/http"

func HealthCheck(rw http.ResponseWriter, req *http.Request) {
    Status(rw, req)
}
