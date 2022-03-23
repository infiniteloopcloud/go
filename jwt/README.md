# JWT

JWT wrapper library which makes it simple to user ECDSA based JWT signing.

### Usage

```go
package main

import (
	"context"

	"github.com/infiniteloopcloud/jwt"
)

type token struct {
	token string `json:"token"`
}

func (t token) ClaimsParse() map[string]interface{} {
	return map[string]interface{}{
		"token": t.token,
    }
}

func main() {
	ctx := context.Background()
	t := token{
		token: "random_token",
    }

	m := jwt.Metadata{
		PrivateKey: "private_key",
		Issuer: "some issuer",
	}

	signed, _ := jwt.Create(ctx, m, t)
	claims, _ := jwt.Verify(ctx, m, signed)
	// Use the claims
}
```