package middlewares

import (
	"net/http"

	"github.com/infiniteloopcloud/hyper"
	"github.com/infiniteloopcloud/log"
)

type LogOpts struct {
	LogLevel uint8
}

// SetLogLevel middleware is responsible for set the log level
func SetLogLevel(opts ...LogOpts) func(next http.Handler) http.Handler {
	var o LogOpts
	if len(opts) == 1 {
		o = opts[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.SetLevel(o.LogLevel)
			next.ServeHTTP(w, r)
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
