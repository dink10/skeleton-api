package filter

import (
  "net/http"
  "context"
  
  "bitbucket.org/gismart/{{Name}}/app/models"
  "bitbucket.org/gismart/{{Name}}/services/render"
)

func QueryCtxGenerator(model models.ModelQuery) func(next http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      q, err := QValues{r.URL.Query()}.ParseQuery(model)
      if err != nil {
        render.MustRenderJSONError(w, r, render.HTTPBadRequest(err))
        return
      }
      
      ctx := context.WithValue(r.Context(), contextQuery, q)
      next.ServeHTTP(w, r.WithContext(ctx))
    })
  }
}
