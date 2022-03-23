# Middlewares

Commonly used middlewares in microservices

### Correlation

Try to access the `Ctx_correlation_id` in the header. If it's not exists it will create a new one and attach it to context.

```go
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/infiniteloopcloud/middlewares"
)

func main() {
	r := chi.NewRouter()
	r.Use(middlewares.Correlation())
}
```

### Context

Context middleware load the values sent in the header from the sender's context into the request's context.

```go
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/infiniteloopcloud/middlewares"
)

func main() {
	r := chi.NewRouter()
	r.Use(middlewares.Context())
}
```

### CORS

Setup CORS's server-side settings

```go
package main

import (
	"github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
	"github.com/infiniteloopcloud/middlewares"
)

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(middlewares.Cors()))
}
```

### Log

Add custom request/response logging

```go
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/infiniteloopcloud/log"
	"github.com/infiniteloopcloud/middlewares"
)

func main() {
	r := chi.NewRouter()
	r.Use(middlewares.CustomLog(middlewares.LogOpts{LogLevel: log.LevelToUint()}))
}
```

### Prometheus

Collect metrics of request/responses

```go
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/infiniteloopcloud/middlewares"
)

func main() {
	r := chi.NewRouter()
	r.Use(middlewares.Prometheus)
}
```