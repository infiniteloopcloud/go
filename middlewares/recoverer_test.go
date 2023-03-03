package middlewares

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func panicingHandler(http.ResponseWriter, *http.Request) { panic("foo") }

func TestRecoverer(t *testing.T) {
	r := chi.NewRouter()

	oldRecovererErrorWriter := recovererErrorWriter
	defer func() { recovererErrorWriter = oldRecovererErrorWriter }()
	buf := &bytes.Buffer{}
	recovererErrorWriter = buf

	r.Use(Recoverer)
	r.Get("/", panicingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/", nil)
	assertEqual(t, res.StatusCode, http.StatusInternalServerError)

	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "->") {
			if !strings.Contains(line, "panicingHandler") {
				t.Fatalf("First func call line should refer to panicingHandler, but actual line:\n%v\n", line)
			}
			return
		}
	}
	t.Fatal("First func call line should start with ->.")
}

func TestRecovererAbortHandler(t *testing.T) {
	defer func() {
		rcv := recover()
		if rcv != http.ErrAbortHandler {
			t.Fatalf("http.ErrAbortHandler should not be recovered")
		}
	}()

	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(w, req)
}

func TestRecovererCustomParser(t *testing.T) {
	r := chi.NewRouter()

	oldRecovererErrorWriter := recovererErrorWriter
	defer func() { recovererErrorWriter = oldRecovererErrorWriter }()
	buf := &bytes.Buffer{}
	recovererErrorWriter = buf

	SetRecoverParser(custom{})
	r.Use(Recoverer)
	r.Get("/", panicingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/", nil)
	assertEqual(t, res.StatusCode, http.StatusInternalServerError)

	lines := strings.Split(buf.String(), "\n")
	if len(lines) > 0 {
		if !strings.Contains(lines[0], "goroutine") {
			t.Fatalf("First func call line should refer to goroutine, but actual line:\n%v\n", lines[0])
		}
		return
	}
	t.Fatal("First func call line should be something like 'goroutine 18 [running]:'")
}

type custom struct{}

func (c custom) Parse(debugStack []byte, rvr interface{}) ([]byte, error) {
	stack := strings.Split(string(debugStack), "\n")

	if len(stack) > 0 {
		return []byte(stack[0]), nil
	}
	return debugStack, nil
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func assertEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("expecting values to be equal but got: '%v' and '%v'", a, b)
	}
}
