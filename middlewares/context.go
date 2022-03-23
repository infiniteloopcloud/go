package middlewares

import (
	"net/http"

	"github.com/infiniteloopcloud/hyper"
	"github.com/infiniteloopcloud/log"
)

func Context(contextKeys map[string]log.ContextField) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = hyper.FromHeader(ctx, r.Header, contextKeys)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
