package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/infiniteloopcloud/go/middlewares/contaxt"
)

func Tracing() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(
				context.WithValue(r.Context(), contaxt.TracingTime, time.Now()),
			)

			next.ServeHTTP(w, r)
		})
	}
}
