package server

import (
	"context"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/gismart/{{Name}}/database"
	"bitbucket.org/gismart/{{Name}}/services/render"
)

func CORSMiddlewareGenereator(allowedOrigins, allowedMethods, allowedHeaders, credentials string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientOrigin := strings.ToLower(r.Header.Get("Origin"))
			origins := strings.Split(allowedOrigins, ",")

			for _, origin := range origins {
				if strings.TrimSpace(strings.ToLower(origin)) == clientOrigin {
					w.Header().Set("Access-Control-Allow-Origin", clientOrigin)
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Allow-Credentials", credentials)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func storageMiddlewareGenereator(withTracing bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := database.PostgresTransaction(r.Context(), withTracing)
			if err != nil {
				render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
				return
			}
			defer func() {
				if err := tx.Rollback(); err != nil {
					log.Errorf("Transaction was rolled back: %s", err)
				}
			}()

			ctx := context.WithValue(r.Context(), database.ContextStorage, tx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
