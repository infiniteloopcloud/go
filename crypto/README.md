# Crypto

Crypto package implements state-of-the-art hashing algorithms and elliptic helper functions.

### Usage

#### Random string

```go
package main

import "github.com/infiniteloopcloud/crypto"

func main() {
	// Generates a cryptographically secure random string
	crypto.RandomString(10)
}
```

#### Hashing

```go
package main

import "github.com/infiniteloopcloud/crypto"

func main() {
	// Get a crypto algorithm, options: Argon2id
	alg, _ := crypto.Get(crypto.Argon2id)
	
	// Generate salt and hash a string
	salt := crypto.RandomString(10)
	hash := alg.Hash("data", salt)
	
	// Check data against the hash
	alg.Verify("data", salt, hash)
}
```

#### Token generation

```go
package main

import (
	"fmt"
	"github.com/infiniteloopcloud/crypto"
)

func main() {
	// Get a crypto algorithm, options: Argon2id
	alg, _ := crypto.Get(crypto.Argon2id)

	// Generate cryptographically secure token
	token := alg.GenerateToken("salt")
	fmt.Println(token)
}
```

#### Elliptic curve helpers

- `MarshalECPublicKey` accept an ecdsa.PublicKey and marshal it to a compressed shareable format
- `UnmarshalECPublicKey` accept a compressed format and parse to an ecdsa.PublicKey