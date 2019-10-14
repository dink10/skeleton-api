package auth

import (
	"context"
	"net/http"

	"bitbucket.org/gismart/{{Name}}/database"
	"bitbucket.org/gismart/{{Name}}/services/render"
)

func initPsRepository(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := database.GetFromContext(r.Context())
		if err != nil {
			render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
			return
		}

		repo := postgres{db}

		ctx := context.WithValue(r.Context(), authStorage, repo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
