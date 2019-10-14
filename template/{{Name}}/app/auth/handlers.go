package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"

	"bitbucket.org/gismart/{{Name}}/app/models"
	"bitbucket.org/gismart/{{Name}}/config"
	"bitbucket.org/gismart/{{Name}}/services/render"
	"bitbucket.org/gismart/{{Name}}/services/request"
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

	if err := setCookie(w, user); err != nil {
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
	data, err := getUserInfoFromGoogle(code, r.Context())

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

	if err := setCookie(w, data); err != nil {
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
	deleteCookie(w)

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

func setCookie(w http.ResponseWriter, user *models.User) error {
	_, jwtString, err := GetTokenAuth().Encode(jwt.MapClaims{"email": user.Email})

	if err != nil {
		return err
	}

	cookie := BaseCookie

	cookie.Value = jwtString
	cookie.MaxAge = int(cookieAge / time.Second)
	cookie.Domain = config.Config.Server.Host

	http.SetCookie(w, &cookie)

	return nil
}

func deleteCookie(w http.ResponseWriter) {
	cookie := BaseCookie

	cookie.Value = ""
	cookie.MaxAge = 0
	cookie.Domain = config.Config.Server.Host

	http.SetCookie(w, &cookie)
}

func getUserInfoFromGoogle(code string, ctx context.Context) (*models.User, error) {
	token, err := getOAuthConfig().Exchange(ctx, code)

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(oauthInfoURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("access_token", token.AccessToken)
	u.RawQuery = q.Encode()

	_, resp, err := request.MakeRequest(ctx, u, "GET", nil)
	if err != nil {
		return nil, err
	}

	var userInfo models.User

	if err := json.Unmarshal(resp, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, err
}
