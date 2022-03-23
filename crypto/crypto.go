// Package crypto provides cryptographic solutions for the user package.
// The main responsibility to define one or more hashing methods to
// store passwords and other sensitive data securely.
package crypto

import (
	crand "crypto/rand"
	"errors"
	"math/big"
	"math/rand"
	"time"
)

const (
	// Argon2id defines a type for argon2 hashing as enum
	Argon2id uint8 = iota
)

const (
	// DefaultSaltLength defines the default length when salt generation happening
	DefaultSaltLength = 18

	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-=_+{}[]:';,./<>?!@#$%^&*()"
)

// Descriptor defines what a hashing implementation should be able to do
type Descriptor interface {
	// Hash should be able to do a hashing on an arbitrary string
	Hash(str, salt string) string

	// Verify should be able to verify a previously created hash
	// by compering with the arbitrary string
	Verify(str, salt, hash string) error

	// GenerateToken should be able to generate a token
	// which hardened with the hashing algorithm
	GenerateToken(tokenSalt string) string
}

// Get will return a hash algorithm for use, or returns error
func Get(typ uint8) (Descriptor, error) {
	if typ == Argon2id {
		return argon2id{}, nil
	}
	return nil, errors.New("undefined crypto algorithm")
}

// RandomString generates random bytes and returns as string
func RandomString(n int) string {
	b := make([]byte, n)
	var err error
	for i := range b {
		var j int64
		j, err = genRandNum(0, int64(len(letterBytes)-1))
		if err != nil {
			break
		}
		b[i] = letterBytes[j]
	}

	// In case of error we return a non-crypto safe string
	if err != nil {
		rand.Seed(time.Now().UTC().UnixNano())
		bm := make([]byte, n)
		for i := range b {
			bm[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		return string(bm)
	}

	return string(b)
}

func genRandNum(min, max int64) (int64, error) {
	bg := big.NewInt(max - min)

	n, err := crand.Int(crand.Reader, bg)
	if err != nil {
		return 0, err
	}

	return n.Int64() + min, nil
}
