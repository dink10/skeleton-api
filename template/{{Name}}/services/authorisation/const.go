package authorisation

import (
  "net/http"
  "time"
  
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "github.com/go-chi/jwtauth"
  
  "bitbucket.org/gismart/{{Name}}/app/models"
  "bitbucket.org/gismart/{{Name}}/config"
)

const (
  ContextAuthUser = models.ContextKey("authUser")
  oauthInfoURL    = "https://www.googleapis.com/oauth2/v2/userinfo"
  cookieAge       = 7 * 24 * time.Hour
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

func getTokenAuth() *jwtauth.JWTAuth {
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



