package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "bitbucket.org/gismart/{{Name}}/docs"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"

		referer := r.Header.Get("Referer")
		if referer != "" {
			scheme = strings.Split(referer, "://")[0]
		}

		httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf(
			"%s://%s/swagger/doc.json",
			scheme,
			r.Host,
		)))(w, r)
	})

	return r
}
