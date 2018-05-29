package server

import (
    "net/http"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "bitbucket.org/gismart/{{Name}}/api"
    "github.com/go-chi/render"
    log "bitbucket.org/gismart/marketingtool/logger"
)

func runRoute() http.Handler {
    r := chi.NewRouter()

    r.Use(render.SetContentType(render.ContentTypeJSON))

    r.Use(middleware.RequestID)
    r.Use(log.RequestLogger())
    r.Use(middleware.Recoverer)
    r.Use(middleware.StripSlashes)

    r.Get("/health", api.HealthCheck)
    r.Get("/status", api.Status)

    return r
}
