package cfg

import (
	"context"
	"log"
	"os"
)

var instance = must()

const (
	cfgDebug     = "CFG_DEBUG"
	defaultPath  = "/config.json"
	pathEnv      = "CFG_LOCATION"
	hotReloadEnv = "CFG_SKIP_HOT_RELOAD"
)

func Get(key string) string {
	return instance.Get(key)
}

func SetErrorf(fn func(ctx context.Context, err error, format string, args ...interface{})) {
	instance.opts.Errorf = fn
}

func SetInfof(fn func(ctx context.Context, format string, args ...interface{})) {
	instance.opts.Infof = fn
}

func SetWarnf(fn func(ctx context.Context, format string, args ...interface{})) {
	instance.opts.Warnf = fn
}

func Rebuild() {
	instance = must()
}

func must() *cfg {
	var path = defaultPath
	if p := os.Getenv(pathEnv); p != "" {
		path = p
	}

	var hotReload = true
	if h := os.Getenv(hotReloadEnv); h == "true" {
		hotReload = false
	}
	c, err := New(Opts{
		Path:      path,
		HotReload: hotReload,
	})
	if err != nil {
		if os.Getenv(cfgDebug) == "true" {
			log.Printf("cfg::New: %s", err.Error())
		}
	}
	return &c
}
