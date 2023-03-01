package middlewares

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infiniteloopcloud/hyper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXSSBodyValidationContinue(t *testing.T) {
	validation := XSSBodyValidation()
	handler := validation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var data map[string]interface{}
		require.Nil(t, hyper.Bind(ctx, r, &data))

		_, ok := data["key"]
		assert.True(t, ok)

		hyper.Success(ctx, w, nil)
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test_endpoint", bytes.NewBufferString(`{"key":"value"}`))
	handler.ServeHTTP(w, r)

	response := w.Result()
	require.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestXSSBodyValidationFail(t *testing.T) {
	validation := XSSBodyValidation()
	handler := validation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var data map[string]interface{}
		require.Nil(t, hyper.Bind(ctx, r, &data))

		_, ok := data["key"]
		assert.True(t, ok)

		hyper.Success(ctx, w, nil)
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test_endpoint", bytes.NewBufferString(`{"key":"<script>alert('hello');</script>"}`))
	handler.ServeHTTP(w, r)

	response := w.Result()
	require.NotNil(t, response)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func BenchmarkXSSBodyValidation(b *testing.B) {
	validation := XSSBodyValidation()
	handler := validation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var data map[string]interface{}
		if err := hyper.Bind(ctx, r, &data); err != nil {
			hyper.Error(ctx, w, err)
			return
		}

		if _, ok := data["key"]; !ok {
			hyper.ReturnBadRequest(ctx, w, "invalid key", errors.New("missing key"))
			return
		}

		hyper.Success(ctx, w, nil)
	}))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/test_endpoint", bytes.NewBufferString(`{"key":"value"}`)))
	}
}
