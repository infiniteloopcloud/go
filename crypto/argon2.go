package crypto

import (
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type argon2id struct{}

// Hash provides Argon2 hashing on an arbitrary string and the corresponding salt.
// The IDKey parameter information can be found in the GoDocs of the IDKey.
func (argon2id) Hash(str, salt string) string {
	return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(str), []byte(salt), 1, 64*1024, 4, 32))
}

// Verify calls the Hash function on the arbitrary string and the corresponding salt,
// and compare with the hash
func (a argon2id) Verify(str, salt, hash string) error {
	if a.Hash(str, salt) != hash {
		return errors.New("invalid password for that hash")
	}
	return nil
}

// GenerateToken creates a new UUID and hash it
func (a argon2id) GenerateToken(tokenSalt string) string {
	return a.Hash(uuid.New().String(), tokenSalt)
}
