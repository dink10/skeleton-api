package auth

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Router returns a custom router
func Router(authMiddlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(initPsRepository)

	r.Get("/login", login)

	r.Group(func(r chi.Router) {
		r.Use(authMiddlewares...)
		r.Get("/me", me)
		r.Post("/logout", logout)
	})

	return r
}
