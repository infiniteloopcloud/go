package middlewares

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/infiniteloopcloud/hyper"
	"github.com/infiniteloopcloud/log"
	"github.com/stretchr/testify/assert"
)

func TestCorrelation(t *testing.T) {
	r := chi.NewRouter()
	r.Use(Correlation())

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cid := ctx.Value(log.CorrelationID)
		assert.Equal(t, "test_correlation_id", cid)

		w.WriteHeader(http.StatusOK)
	})

	srv := http.Server{
		Addr:    ":11112",
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(10 * time.Millisecond)

	req, err := http.NewRequest("GET", "http://localhost:11112/hello", nil)
	if err != nil {
		t.Error(err)
	}
	ctx := context.WithValue(context.Background(), log.CorrelationID, "test_correlation_id")
	req.Header = hyper.IntoHeader(ctx, req.Header, map[string]log.ContextField{string(log.CorrelationID): log.CorrelationID})

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
}
