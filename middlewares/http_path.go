package middlewares

import (
	"context"
	"net/http"

	"github.com/infiniteloopcloud/log"
	"github.com/volatiletech/null/v9"
)

func HTTPPath() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !getHTTPPath(r.Context()).Valid {
				r = r.WithContext(context.WithValue(r.Context(), log.HTTPPath, r.URL.Path))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getHTTPPath(ctx context.Context) null.String {
	ctxVal := ctx.Value(log.HTTPPath)
	if v, ok := ctxVal.(string); ok {
		return null.StringFrom(v)
	}
	return null.String{}
}
