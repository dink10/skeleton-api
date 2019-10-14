package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"bitbucket.org/gismart/{{Name}}/app/models"
	"bitbucket.org/gismart/{{Name}}/config"
)

const (
	oauthInfoURL    = "https://www.googleapis.com/oauth2/v2/userinfo"
	cookieAge       = 7 * 24 * time.Hour
	ContextAuthUser = models.ContextKey("authUser")
	authStorage     = models.StorageKey("auth")
)

var (
	tokenAuth   *jwtauth.JWTAuth
	oAuthConfig *oauth2.Config
	BaseCookie  = http.Cookie{
		Name:     "_gt",
		HttpOnly: true,
		Path:     "/",
		// Secure:   true,
	}
)

func GetTokenAuth() *jwtauth.JWTAuth {
	if tokenAuth == nil {
		tokenAuth = jwtauth.New("HS256", []byte(config.Config.Oauth.JWTSecret), nil)
	}

	return tokenAuth
}

func getOAuthConfig() *oauth2.Config {
	if oAuthConfig == nil {
		oAuthConfig = &oauth2.Config{
			// google oauth2 hack
			RedirectURL:  "postmessage",
			ClientID:     config.Config.Oauth.ClientID,
			ClientSecret: config.Config.Oauth.ClientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		}
	}

	return oAuthConfig
}

func GetFromContext(ctx context.Context) (*models.User, error) {
	if user, ok := ctx.Value(ContextAuthUser).(*models.User); ok {
		return user, nil
	}

	return nil, errors.New("can not cast to *models.User")
}

func getStorageFromContext(ctx context.Context) (*postgres, error) {
	if repo, ok := ctx.Value(authStorage).(postgres); ok {
		return &repo, nil
	}

	return nil, errors.New("can not cast to repository")
}
