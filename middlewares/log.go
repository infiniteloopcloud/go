package middlewares

import (
	"context"
	"net/http"

	"github.com/infiniteloopcloud/hyper"
	"github.com/infiniteloopcloud/log"
)

type LogOpts struct {
	LogLevel uint8
}

// CustomLog middleware is responsible for setting the logLevel,
// and also logs the response body
func CustomLog(opts ...LogOpts) func(next http.Handler) http.Handler {
	var o LogOpts
	if len(opts) == 1 {
		o = opts[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log.SetLevel(o.LogLevel)
			ctx = context.WithValue(ctx, log.HTTPPath, r.RequestURI)
			rw := hyper.NewWriter(w)
			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)
		})
	}
}

// LogResponse middleware is responsible for log the response body
func LogResponse() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw, isHyperWriter := w.(*hyper.Writer)
			if !isHyperWriter {
				rw = hyper.NewWriter(w)
			}
			next.ServeHTTP(rw, r)
			log.Debug(r.Context(), rw.ResponseBody())
		})
	}
}
