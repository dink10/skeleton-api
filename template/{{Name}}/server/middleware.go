package server

import (
	"net/http"
	"fmt"
	nettrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	cfg "bitbucket.org/gismart/piano2/config"
)

// Tracer enables tracing for single http request with correct resource name
func Tracer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		resource := req.URL.Path
		next = nettrace.WrapHandler(next, fmt.Sprintf("{{Name}}-%s", cfg.Config.Logger.DataDogEnv), resource)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}