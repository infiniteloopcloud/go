package contaxt

import (
	"context"

	"gopkg.in/volatiletech/null.v6"

	"github.com/infiniteloopcloud/log"
)

const (
	ClientHost  log.ContextField = "client_host"
	HTTPPath    log.ContextField = "http_path"
	TracingTime log.ContextField = "tracing_time"
)

func GetClientHost(ctx context.Context) null.String {
	ctxVal := ctx.Value(ClientHost)
	if v, ok := ctxVal.(string); ok {
		return null.StringFrom(v)
	}
	return null.String{}
}

func GetHTTPPath(ctx context.Context) null.String {
	ctxVal := ctx.Value(HTTPPath)
	if v, ok := ctxVal.(string); ok {
		return null.StringFrom(v)
	}
	return null.String{}
}

func GetTracingTime(ctx context.Context) null.String {
	ctxVal := ctx.Value(TracingTime)
	if v, ok := ctxVal.(string); ok {
		return null.StringFrom(v)
	}
	return null.String{}
}
