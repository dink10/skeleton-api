package authorisation

import (
  "context"
  "errors"
  "net/http"
  
  "github.com/go-chi/jwtauth"
  
  "bitbucket.org/gismart/{{Name}}/services/render"
  "bitbucket.org/gismart/{{Name}}/app/models"
  "bitbucket.org/gismart/{{Name}}/database"
)

func AuthCtx(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    _, userInfo, err := jwtauth.FromContext(r.Context())
    if err != nil {
      render.MustRenderJSONError(w, r, render.HTTPBadRequest(err))
      return
    }
    
    userEmail, ok := userInfo["email"].(string)
    if !ok {
      render.MustRenderJSONError(w, r, render.HTTPBadRequest(errors.New("invalid cookie")))
      return
    }
    
    model := &models.User{}
    
    db, err := database.GetFromContext(r.Context())
    if err != nil {
      render.MustRenderJSONError(w, r, render.HTTPInternalServerError(err))
      return
    }
    
    if err := db.GetBy(map[string]interface{}{"email": userEmail}, model); err != nil {
      render.MustRenderJSONError(w, r, render.HTTPUnauthorized)
      return
    }
    
    ctx := context.WithValue(r.Context(), ContextAuthUser, model)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}