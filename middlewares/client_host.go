package middlewares

import (
	"context"
	"net/http"

	"github.com/infiniteloopcloud/go/middlewares/contaxt"
)

func ClientHost() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !contaxt.GetClientHost(r.Context()).Valid {
				r = r.WithContext(context.WithValue(r.Context(), contaxt.ClientHost, r.Host))
			}

			next.ServeHTTP(w, r)
		})
	}
}
