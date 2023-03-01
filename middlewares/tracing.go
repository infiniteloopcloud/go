package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/infiniteloopcloud/log"
)

func Tracing() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(
				context.WithValue(r.Context(), log.TracingTime, time.Now()),
			)

			next.ServeHTTP(w, r)
		})
	}
}
