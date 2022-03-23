package middlewares

import (
	"context"
	"net/http"

	"github.com/infiniteloopcloud/go/middlewares/contaxt"
)

func HTTPPath() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !contaxt.GetHTTPPath(r.Context()).Valid {
				r = r.WithContext(context.WithValue(r.Context(), contaxt.HTTPPath, r.URL.Path))
			}

			next.ServeHTTP(w, r)
		})
	}
}
