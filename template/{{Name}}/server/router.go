package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"

	"bitbucket.org/gismart/ddtracer"

	"bitbucket.org/gismart/{{Name}}/app/auth"
	"bitbucket.org/gismart/{{Name}}/app/health"
	"bitbucket.org/gismart/{{Name}}/config"
	log "bitbucket.org/gismart/{{Name}}/services/logger"
	"bitbucket.org/gismart/{{Name}}/services/authorisation"
	"bitbucket.org/gismart/{{Name}}/swagger"
)

func runRoute(withTracing bool) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(log.RequestLogger())
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(CORSMiddlewareGenereator(config.Config.AllowedOrigins, config.Config.AllowedMethods, config.Config.AllowedHeaders, config.Config.Credentials))

	r.Get("/health", health.Health)
	r.Mount("/swagger", swagger.Router())

	r.Group(func(r chi.Router) {
		if withTracing {
			r.Use(ddtracer.TraceMiddleware)
		}
		r.Use(storageMiddlewareGenereator(withTracing))
		r.Mount("/auth", auth.Router(authorisation.Verifier(), jwtauth.Authenticator, authorisation.AuthCtx))

		// Authentication protected routes
		r.Group(func(r chi.Router) {
			r.Use(authorisation.Verifier(), jwtauth.Authenticator, authorisation.AuthCtx)
		})
	})

	return r
}
