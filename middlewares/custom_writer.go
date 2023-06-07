package middlewares

import (
	"net/http"

	"github.com/infiniteloopcloud/hyper"
)

// CustomWriter wraps the http.ResponseWriter into a new hyper.Writer
func CustomWriter() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(hyper.NewWriter(w), r)
		})
	}
}
