package middlewares

import (
	"context"
	"strings"

	"github.com/go-chi/cors"
	"github.com/infiniteloopcloud/log"
)

type CorsOpts struct {
	AllowedOrigins string
}

func Cors(opts ...CorsOpts) cors.Options {
	var o CorsOpts
	if len(opts) == 1 {
		o = opts[0]
		if o.AllowedOrigins == "" {
			o.AllowedOrigins = "*"
		}
	} else {
		o.AllowedOrigins = "*"
	}

	c := cors.Options{
		AllowedOrigins:   strings.Split(o.AllowedOrigins, ";"),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           400,
	}
	log.Debugf(context.Background(), "Cors config: %+v", c)
	return c
}
