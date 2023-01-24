package crypto_test

import (
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/infiniteloopcloud/go/crypto"
)

// DefaultPrivateKey is for NON-PRODUCTION use
//
//nolint:lll
const DefaultPrivateKey = "MHcCAQEEIMFk+5s8b1hHM6W98IqPgMzUSc8LDJTHqQ2kdYskM4ASoAoGCCqGSM49AwEHoUQDQgAEP0qVlnDM/hlIAg7pzixSgtZfpH2f5C9lw7B/ZZb0bTB2QzyrRCyKeIwzQ0oqMVTLJ7oVvZZshsCDmuv4vDjZQw=="

func TestMarshalPublicKey(t *testing.T) {
	privateKeyParsed, err := base64.StdEncoding.DecodeString(DefaultPrivateKey)
	if err != nil {
		t.Error(err)
	}
	key, err := x509.ParseECPrivateKey(privateKeyParsed)
	if err != nil {
		t.Error(err)
	}
	compressed := crypto.MarshalECPublicKey(key.PublicKey)
	publicK, err := crypto.UnmarshalECPublicKey(compressed)
	if err != nil {
		t.Error(err)
	}

	if publicK.X.Int64() != key.PublicKey.X.Int64() {
		t.Errorf("Unmarshaled public key %d, is not the same as %d", publicK.X.Int64(), key.PublicKey.X.Int64())
	}
}
