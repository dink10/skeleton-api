package authorisation

import (
  "context"
  "encoding/json"
  "net/http"
  "net/url"
  "time"
  
  "github.com/dgrijalva/jwt-go"
  "github.com/go-chi/jwtauth"
  
  "bitbucket.org/gismart/{{Name}}/app/models"
  "bitbucket.org/gismart/{{Name}}/services/request"
  "bitbucket.org/gismart/{{Name}}/config"
)

func SetCookie(w http.ResponseWriter, user *models.User) error {
  _, jwtString, err := getTokenAuth().Encode(jwt.MapClaims{"email": user.Email})
  
  if err != nil {
    return err
  }
  
  cookie := BaseCookie
  
  cookie.Value = jwtString
  cookie.MaxAge = int(cookieAge / time.Second)
  cookie.Domain = config.Config.Server.Host
  
  // i work with the http.Header because of http.SetCookie use Header().Add() instead of Header().Set()
  w.Header().Set("Set-Cookie", cookie.String())
  
  return nil
}

func DeleteCookie(w http.ResponseWriter) {
  cookie := BaseCookie
  
  cookie.Value = ""
  cookie.MaxAge = 0
  cookie.Domain = config.Config.Server.Host
  
  w.Header().Set("Set-Cookie", cookie.String())
}

func GetUserInfoFromGoogle(code string, ctx context.Context) (*models.User, error) {
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

func Verifier() func(http.Handler) http.Handler {
  return jwtauth.Verify(getTokenAuth(), func(r *http.Request) string {
    cookie, err := r.Cookie(BaseCookie.Name)
    if err != nil {
      return ""
    }
    return cookie.Value
  })
}