package request

import (
  "context"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "net/url"
  "testing"
  "time"
  
  "bitbucket.org/gismart/{{Name}}/config"
)

var body = "test body"

func TestQuery(t *testing.T) {
  config.Init(&config.Config)
  
  t.Run("Request get testing...", testGetRequest)
  t.Run("Request post testing...", testPostRequest)
  t.Run("Request timeout testing...", testTimeoutRequest)
}

func handler(w http.ResponseWriter, vr *http.Request) {
  time.Sleep(time.Duration(config.Config.HTTPTimeout + 1) * time.Second)
  w.WriteHeader(http.StatusOK)
}

func postHandler(w http.ResponseWriter, vr *http.Request) {
  b, err := ioutil.ReadAll(vr.Body)
  defer vr.Body.Close()
  
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(b))
}

func getHandler(w http.ResponseWriter, vr *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(body))
}

func testGetRequest(t *testing.T) {
  ctx := context.Background()
  backend := httptest.NewServer(http.Handler(http.HandlerFunc(getHandler)))
  u, _ := url.Parse(backend.URL)
  
  code, b, err := MakeRequest(ctx, u, "GET", nil)
  
  if code != http.StatusOK {
    t.Errorf("Invalid get response: \n\t %v \n\t %v", code, err)
    return
  }
  
  if string(b) != body {
    t.Errorf("Invalid response in get request: \n\t %s", b)
  }
}

func testPostRequest(t *testing.T) {
  ctx := context.Background()
  backend := httptest.NewServer(http.Handler(http.HandlerFunc(postHandler)))
  u, _ := url.Parse(backend.URL)
  
  code, b, err := MakeRequest(ctx, u, "POST", []byte(body))
  
  if code != http.StatusOK {
    t.Errorf("Invalid post response: \n\t %d \n\t %v", code, err)
  }
  
  if string(b) != body {
    t.Errorf("Invalid response in post request: \n\t %s", b)
  }
}

func testTimeoutRequest(t *testing.T) {
  ctx := context.Background()
  backend := httptest.NewServer(http.Handler(http.HandlerFunc(handler)))
  u, _ := url.Parse(backend.URL)
  
  code, body, err := MakeRequest(ctx, u, "GET", nil)
  
  if err == nil {
    t.Errorf("Invalid HTTP_TIMEOUT %d %s %d", code, body, config.Config.HTTPTimeout)
    return
  }
}