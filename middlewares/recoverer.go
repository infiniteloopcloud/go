package middlewares

// The original work was derived from chi's middleware, source:
// https://github.com/go-chi/chi/tree/master/web/middleware

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/infiniteloopcloud/log"
)

type RecoverParserIface interface {
	Parse(ctx context.Context, debugStack []byte, rvr interface{}) ([]byte, error)
}

var RecoverParser RecoverParserIface = prettyStack{}

type PrintPrettyStackFn func(ctx context.Context, rvr interface{})

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/go-chi/httplog middleware pkgs.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				//nolint:errorlint
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					// nolint:forbidigo
					printPrettyStackFn(r.Context(), rvr)
				}

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// for ability to test the PrintPrettyStack function
var recovererErrorWriter io.Writer = os.Stderr

var printPrettyStackFn = PrintPrettyStack

func PrintPrettyStack(ctx context.Context, rvr interface{}) {
	debugStack := debug.Stack()
	out, err := RecoverParser.Parse(ctx, debugStack, rvr)
	if err == nil {
		//nolint:errcheck
		recovererErrorWriter.Write(out)
	} else {
		// print stdlib output as a fallback
		os.Stderr.Write(debugStack)
	}
}

func MultilinePrettyPrintStack(ctx context.Context, rvr interface{}) {
	debugStack := debug.Stack()
	parsed, err := panicParser{}.Parse(debugStack)
	if err != nil {
		os.Stderr.Write(debugStack)
		return
	}
	multilinePrettyPrint(ctx, rvr, parsed)
}

func multilinePrettyPrint(ctx context.Context, rvr interface{}, parsed []stackLine) {
	slices.Reverse(parsed)
	panicID := uuid.NewString()
	for _, line := range parsed {
		logLine := log.Parse(ctx, log.ErrorLevelString, "panic happen",
			fmt.Errorf("panic happen: %v", rvr),
			log.Field{
				Key:   "panic_id",
				Value: panicID,
			},
			log.Field{
				Key:   "function_name",
				Value: line.FunctionName,
			},
			log.Field{
				Key:   "file_path",
				Value: line.FilePath,
			},
			log.Field{
				Key:   "line_number",
				Value: strconv.Itoa(line.LineNumber),
			},
		)
		// nolint: errcheck
		recovererErrorWriter.Write([]byte(logLine))
	}
}

func SetRecovererErrorWriter(w io.Writer) {
	recovererErrorWriter = w
}

func SetRecoverParser(p RecoverParserIface) {
	RecoverParser = p
}

// nolint:forbidigo
func SetPrintPrettyStack(fn PrintPrettyStackFn) {
	printPrettyStackFn = fn
}

type prettyStack struct {
}

func (s prettyStack) Parse(_ context.Context, debugStack []byte, rvr interface{}) ([]byte, error) {
	var err error
	useColor := true
	buf := &bytes.Buffer{}

	cW(buf, false, bRed, "\n")
	cW(buf, useColor, bCyan, " panic: ")
	cW(buf, useColor, bBlue, "%v", rvr)
	cW(buf, false, bWhite, "\n \n")

	// process debug stack info
	stack := strings.Split(string(debugStack), "\n")
	lines := []string{}

	// locate panic line, as we may have nested panics
	for i := len(stack) - 1; i > 0; i-- {
		lines = append(lines, stack[i])
		if strings.HasPrefix(stack[i], "panic(") {
			lines = lines[0 : len(lines)-2] // remove boilerplate
			break
		}
	}

	// reverse
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opp := len(lines) - 1 - i
		lines[i], lines[opp] = lines[opp], lines[i]
	}

	// decorate
	for i, line := range lines {
		lines[i], err = s.decorateLine(line, useColor, i)
		if err != nil {
			return nil, err
		}
	}

	for _, l := range lines {
		fmt.Fprintf(buf, "%s", l)
	}
	return buf.Bytes(), nil
}

func (s prettyStack) decorateLine(line string, useColor bool, num int) (string, error) {
	line = strings.TrimSpace(line)
	switch {
	case strings.HasPrefix(line, "\t") || strings.Contains(line, ".go:"):
		return s.decorateSourceLine(line, useColor, num)
	case strings.HasSuffix(line, ")"):
		return s.decorateFuncCallLine(line, useColor, num)
	default:
		if strings.HasPrefix(line, "\t") {
			return strings.Replace(line, "\t", "      ", 1), nil
		} else {
			return fmt.Sprintf("    %s\n", line), nil
		}
	}
}

func (s prettyStack) decorateFuncCallLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, "(")
	if idx < 0 {
		return "", errors.New("not a func call line")
	}

	buf := &bytes.Buffer{}
	pkg := line[0:idx]
	// addr := line[idx:]
	method := ""

	if idx := strings.LastIndex(pkg, string(os.PathSeparator)); idx < 0 {
		if idx := strings.Index(pkg, "."); idx > 0 {
			method = pkg[idx:]
			pkg = pkg[0:idx]
		}
	} else {
		method = pkg[idx+1:]
		pkg = pkg[0 : idx+1]
		if idx := strings.Index(method, "."); idx > 0 {
			pkg += method[0:idx]
			method = method[idx:]
		}
	}
	pkgColor := nYellow
	methodColor := bGreen

	if num == 0 {
		cW(buf, useColor, bRed, " -> ")
		pkgColor = bMagenta
		methodColor = bRed
	} else {
		cW(buf, useColor, bWhite, "    ")
	}
	cW(buf, useColor, pkgColor, "%s", pkg)
	cW(buf, useColor, methodColor, "%s\n", method)
	// cW(buf, useColor, nBlack, "%s", addr)
	return buf.String(), nil
}

func (s prettyStack) decorateSourceLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, ".go:")
	if idx < 0 {
		return "", errors.New("not a source line")
	}

	buf := &bytes.Buffer{}
	path := line[0 : idx+3]
	lineno := line[idx+3:]

	idx = strings.LastIndex(path, string(os.PathSeparator))
	dir := path[0 : idx+1]
	file := path[idx+1:]

	idx = strings.Index(lineno, " ")
	if idx > 0 {
		lineno = lineno[0:idx]
	}
	fileColor := bCyan
	lineColor := bGreen

	if num == 1 {
		cW(buf, useColor, bRed, " ->   ")
		fileColor = bRed
		lineColor = bMagenta
	} else {
		cW(buf, false, bWhite, "      ")
	}
	cW(buf, useColor, bWhite, "%s", dir)
	cW(buf, useColor, fileColor, "%s", file)
	cW(buf, useColor, lineColor, "%s", lineno)
	if num == 1 {
		cW(buf, false, bWhite, "\n")
	}
	cW(buf, false, bWhite, "\n")

	return buf.String(), nil
}

var (
	nYellow  = []byte{'\033', '[', '3', '3', 'm'}
	bRed     = []byte{'\033', '[', '3', '1', ';', '1', 'm'}
	bGreen   = []byte{'\033', '[', '3', '2', ';', '1', 'm'}
	bBlue    = []byte{'\033', '[', '3', '4', ';', '1', 'm'}
	bMagenta = []byte{'\033', '[', '3', '5', ';', '1', 'm'}
	bCyan    = []byte{'\033', '[', '3', '6', ';', '1', 'm'}
	bWhite   = []byte{'\033', '[', '3', '7', ';', '1', 'm'}

	reset = []byte{'\033', '[', '0', 'm'}
)

var IsTTY bool

func init() {
	// This is sort of cheating: if stdout is a character device, we assume
	// that means it's a TTY. Unfortunately, there are many non-TTY
	// character devices, but fortunately stdout is rarely set to any of
	// them.
	//
	// We could solve this properly by pulling in a dependency on
	// code.google.com/p/go.crypto/ssh/terminal, for instance, but as a
	// heuristic for whether to print in color or in black-and-white, I'd
	// really rather not.
	fi, err := os.Stdout.Stat()
	if err == nil {
		m := os.ModeDevice | os.ModeCharDevice
		IsTTY = fi.Mode()&m == m
	}
}

// colorWrite
func cW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	if IsTTY && useColor {
		//nolint:errcheck
		w.Write(color)
	}
	fmt.Fprintf(w, s, args...)
	if IsTTY && useColor {
		//nolint:errcheck
		w.Write(reset)
	}
}

type RecovererLogStack struct {
	DebugStackLengthLimit int
}

func (l RecovererLogStack) Parse(ctx context.Context, debugStack []byte, rvr interface{}) ([]byte, error) {
	debugStackMsgWrapped := strings.Join(strings.Split(string(debugStack), "\n"), "|")
	if l.DebugStackLengthLimit > 0 {
		if lenDebugStack := len(debugStackMsgWrapped); lenDebugStack > l.DebugStackLengthLimit {
			debugStackMsgWrapped = debugStackMsgWrapped[:l.DebugStackLengthLimit]
		}
	}
	parsedLog := log.Parse(ctx, log.ErrorLevelString, "panic happen",
		fmt.Errorf("panic happen: %v", rvr),
		log.Field{
			Key:   "debug_stack",
			Value: debugStackMsgWrapped,
		})
	return []byte(parsedLog), nil
}
