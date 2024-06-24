package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type panicParserWriter struct {
	b *bytes.Buffer
}

func (w *panicParserWriter) Write(b []byte) (n int, err error) {
	w.b.Write(b)
	w.b.WriteByte(byte('\n'))
	return len(b) + 1, nil
}

func TestPanicParser_Parse(t *testing.T) {
	r := chi.NewRouter()

	var parsed []stackLine
	oldRecovererErrorWriter := recovererErrorWriter
	defer func() { recovererErrorWriter = oldRecovererErrorWriter }()
	w := panicParserWriter{b: bytes.NewBuffer(nil)}
	recovererErrorWriter = &w

	r.Use(Recoverer)
	SetPrintPrettyStack(func(ctx context.Context, rvr interface{}) {
		debugStack := debug.Stack()
		var err error
		parsed, err = panicParser{}.Parse(debugStack)
		if err != nil {
			os.Stderr.Write(debugStack)
			return
		}
		multilinePrettyPrint(ctx, rvr, parsed)
	})
	r.Get("/", panicingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/", nil)
	assertEqual(t, res.StatusCode, http.StatusInternalServerError)

	assert.Len(t, parsed, 13)

	var latestPanicID string
	for _, line := range strings.Split(w.b.String(), "\n") {
		if line == "" {
			continue
		}
		var logLine map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &logLine))
		panicID, ok := logLine["panic_id"]
		require.True(t, ok)
		if latestPanicID == "" {
			// nolint: errcheck
			latestPanicID = panicID.(string)
		} else {
			assert.Equal(t, latestPanicID, panicID)
		}
		lineNum, hasLineField := logLine["line_number"]
		require.True(t, hasLineField)
		// nolint: errcheck
		_, err := strconv.Atoi(lineNum.(string))
		require.NoError(t, err)
		_, hasFnName := logLine["function_name"]
		require.True(t, hasFnName)
		_, hasFilePath := logLine["file_path"]
		require.True(t, hasFilePath)
	}
}
