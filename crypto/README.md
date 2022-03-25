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

#### Encrypt/Decrypt

This functionality makes it easier to encrypt arbitrary length strings with AES-128

```go
package main

import (
	"log"

	"github.com/infiniteloopcloud/crypto"
)

func main() {
	dataToEncrypt := "testdatatestdatatestdatatestdatatestdatatestdata1234"

	// NOTE: The key and the IV always should be 16 length
	e := crypto.NewEncypter("bhXRirFB8IaQxxjm", "oZDRWCHryRtlVA1I")
	resultData, err := e.Encrypt(dataToEncrypt)
	if err != nil {
		// handle error
	}

	decrypted, err := e.Decrypt(resultData)
	if err != nil {
		// handle error
	}
	if decrypted != dataToEncrypt {
		log.Printf("Decrypted is not the original data, original: %s, decrypted: %s", dataToEncrypt, decrypted)
	}
}
```

#### Elliptic curve helpers

- `MarshalECPublicKey` accept an ecdsa.PublicKey and marshal it to a compressed shareable format
- `UnmarshalECPublicKey` accept a compressed format and parse to an ecdsa.PublicKey