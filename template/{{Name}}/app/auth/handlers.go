package auth

import (
	"net/http"

	"bitbucket.org/gismart/{{Name}}/services/render"
	"bitbucket.org/gismart/{{Name}}/services/authorisation"
)

// Create godoc
// @Tags auth
// @Produce json
// @Param Cookie header string true "_gt=jwttoken"
// @Success 200 {object} models.User
// @Failure 403 {object} render.ErrorResponse
// @Failure 401 {object} render.ErrorResponse
// @Failure 500 {object} render.ErrorResponse
// @Router /auth/me [get]
func me(w http.ResponseWriter, r *http.Request) {
	user, err := GetFromContext(r.Context())
	if err != nil {
		render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
		return
	}

	db, err := getStorageFromContext(r.Context())
	if err != nil {
		render.MustRenderJSONError(w, r, render.HTTPBadRequest(err))
		return
	}

	if err := authorisation.SetCookie(w, user); err != nil {
		render.MustRenderJSONError(w, r, render.HTTPBadRequest(err))
		return
	}

	if err := db.Commit(); err != nil {
		render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
		return
	}

	render.MustRenderJSON(w, r, user)
}

// Create godoc
// @Tags auth
// @Produce json
// @Param code query string true "google authorisation code"
// @Success 200 {object} models.User
// @Header 200 {string} Set-Cookie "jwt"
// @Failure 403 {object} render.ErrorResponse
// @Failure 401 {object} render.ErrorResponse
// @Failure 500 {object} render.ErrorResponse
// @Router /auth/login [get]
func login(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	data, err := authorisation.GetUserInfoFromGoogle(code, r.Context())

	if err != nil {
		render.MustRenderJSONError(w, r, render.HTTPBadRequest(err))
		return
	}

	db, err := getStorageFromContext(r.Context())
	if err != nil {
		render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
		return
	}

	if err := db.Upsert(data); err != nil {
		render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
		return
	}

	if err := authorisation.SetCookie(w, data); err != nil {
		render.MustRenderJSONError(w, r, render.HTTPBadRequest(err))
		return
	}

	if err := db.Commit(); err != nil {
		render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
		return
	}

	render.MustRenderJSON(w, r, data)
}

// Create godoc
// @Tags auth
// @Produce json
// @Param Cookie header string true "_gt=jwttoken"
// @Success 200 {object} models.User
// @Failure 500 {object} render.ErrorResponse
// @Router /auth/logout [get]
func logout(w http.ResponseWriter, r *http.Request) {
	authorisation.DeleteCookie(w)

	db, err := getStorageFromContext(r.Context())
	if err != nil {
		render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
		return
	}

	if err := db.Commit(); err != nil {
		render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
