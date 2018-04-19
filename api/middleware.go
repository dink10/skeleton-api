package api

import (
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
)

// loggerMiddleware logs all HTTP requests to output in a proper format
func loggerMiddleware() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			log.Debugf("%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				time.Since(start),
			)
		})
	}
}
