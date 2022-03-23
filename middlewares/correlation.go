package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/infiniteloopcloud/hyper"
)

func Correlation() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationID := r.Header.Get(hyper.HeaderPrefix + string(hyper.CorrelationID))
			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			r = r.WithContext(
				context.WithValue(r.Context(), hyper.CorrelationID, correlationID),
			)

			next.ServeHTTP(w, r)
		})
	}
}
