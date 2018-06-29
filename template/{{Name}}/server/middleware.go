package server

import (
	"net/http"
	nettrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

// Tracer enables tracing for single http request with correct resource name
func Tracer(service string) (func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			resource := req.URL.Path
			next = nettrace.WrapHandler(next, service, resource)
			next.ServeHTTP(w, req)
		})
	}
}