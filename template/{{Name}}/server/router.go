package server

import (
    "net/http"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "bitbucket.org/gismart/{{Name}}/api"
    "github.com/go-chi/render"
)

func runRoute() http.Handler {
    r := chi.NewRouter()

    r.Use(render.SetContentType(render.ContentTypeJSON))

    r.Use(middleware.RequestID)
    r.Use(RequestLogger())
    r.Use(middleware.Recoverer)

    r.Get("/health/", api.HealthCheck)

    return r
}
