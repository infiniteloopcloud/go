# cfg

Yet another config package

### Features

- Read from file
- Read from environment variable
- Hot reload of the file

### Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/infiniteloopcloud/cfg"
)

var config cfg.Config

func main() {
	var err error
	config, err = cfg.New(cfg.Opts{
		Path:      "./config.json",
		HotReload: true,
		Infof: func(ctx context.Context, format string, args ...interface{}) {
			fmt.Printf(format+"\n", args...)
		},
	})
	if err != nil {
		// TODO err
    }

	// First it reads from the ./config.json (it's a map[string]string under the hood)
	// After it reads with os.Getenv
	config.Get("something")
}
```
