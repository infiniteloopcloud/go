# Migration

Commonly used migration tools

### Usage

```go
package main

import (
	"context"
	"embed"

	"github.com/infiniteloopcloud/migration"
)

//go:embed files
var Files embed.FS

func main() {
	ctx := context.Background()

	// Library usage
	migration.Get(ctx, "database connection string", Files)

	// Executable usage
	migration.MigrationTool("database connection string", false, Files)
}
```