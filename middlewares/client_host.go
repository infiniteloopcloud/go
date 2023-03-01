package middlewares

import (
	"context"
	"net/http"

	"github.com/infiniteloopcloud/log"
	"github.com/volatiletech/null/v9"
)

func ClientHost() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !getClientHost(r.Context()).Valid {
				r = r.WithContext(context.WithValue(r.Context(), log.ClientHost, r.Host))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientHost(ctx context.Context) null.String {
	ctxVal := ctx.Value(log.ClientHost)
	if v, ok := ctxVal.(string); ok {
		return null.StringFrom(v)
	}
	return null.String{}
}
