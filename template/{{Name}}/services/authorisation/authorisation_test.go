package authorisation

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/dgrijalva/jwt-go"
	
	"bitbucket.org/gismart/{{Name}}/app/models"
	"bitbucket.org/gismart/{{Name}}/config"
)

var (
	user = models.User{
		Name: "test",
		Email: "test@gmail.com",
	}
)

func TestAuthorisation(t *testing.T) {
	config.Init(&config.Config)
	
	t.Run("Testing cookie setting...", testCookieSetting)
	t.Run("Testing cookie deleting...", testCookieDeleting)
}

func testCookieSetting(t *testing.T) {
	w := httptest.NewRecorder()
	
	if err := SetCookie(w, &user); err != nil {
		t.Error("incorrect work of SetCookie", err)
		return
	}
	
	r := &http.Request{Header: http.Header{"Cookie": []string{w.Header().Get("Set-Cookie")}}}
	cookie, err := r.Cookie(BaseCookie.Name)
	if err != nil {
		t.Error("invalid cookie name", err)
		return
	}
	
	_, jwt, err := getTokenAuth().Encode(jwt.MapClaims{"email": user.Email})
	if err != nil {
		t.Error("encode jwt error", err)
	}
	
	if jwt != cookie.Value {
		t.Errorf("invalid cookie \n\t %+v \n\t %+v", jwt, cookie.Value)
	}
}

func testCookieDeleting(t *testing.T) {
	w := httptest.NewRecorder()
	
	if err := SetCookie(w, &user); err != nil {
		t.Error("incorrect work of SetCookie", err)
		return
	}
	
	DeleteCookie(w)
	
	r := &http.Request{Header: http.Header{"Cookie": []string{w.Header().Get("Set-Cookie")}}}
	cookie, err := r.Cookie(BaseCookie.Name)
	if err != nil {
		t.Error("invalid cookie name", err)
		return
	}
	
	if cookie.Value != "" {
		t.Errorf("invalid cookie \n\t %+v", cookie.Value)
	}
}
