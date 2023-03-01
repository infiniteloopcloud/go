package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/infiniteloopcloud/hyper"
	"github.com/infiniteloopcloud/log"
	xssvalidator "github.com/infiniteloopcloud/xss-validator"
)

func XSSBodyValidation() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var body json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				hyper.ReturnBadRequest(ctx, w, hyper.InvalidRequest, errors.New(hyper.ErrBind+err.Error()))
				return
			}

			err := xssvalidator.Validate(string(body), xssvalidator.DefaultRules...)
			if err != nil {
				log.Errorf(ctx, err, "xss validation error")
				hyper.ReturnBadRequest(ctx, w, hyper.InvalidRequest, errors.New(hyper.ErrBind))
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(body))

			next.ServeHTTP(w, r)
		})
	}
}
