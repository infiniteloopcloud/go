package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/infiniteloopcloud/hyper"
	"github.com/infiniteloopcloud/log"
)

func Correlation() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationID := r.Header.Get(hyper.HeaderPrefix + string(log.CorrelationID))
			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			r = r.WithContext(
				context.WithValue(r.Context(), log.CorrelationID, correlationID),
			)

			next.ServeHTTP(w, r)
		})
	}
}
