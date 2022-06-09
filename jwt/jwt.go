// Package jwt provides secure signature-based token sharing with the client-side.
// Implements a JWT standard for providing this functionality.
package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"time"

	"github.com/infiniteloopcloud/go/crypto"
	"github.com/pascaldekloe/jwt"
)

type Metadata struct {
	PrivateKey string
	Issuer     string
	ClientHost string
}

// DefaultPrivateKey is for NON-PRODUCTION use
//nolint:lll
const DefaultPrivateKey = "MHcCAQEEIMFk+5s8b1hHM6W98IqPgMzUSc8LDJTHqQ2kdYskM4ASoAoGCCqGSM49AwEHoUQDQgAEP0qVlnDM/hlIAg7pzixSgtZfpH2f5C9lw7B/ZZb0bTB2QzyrRCyKeIwzQ0oqMVTLJ7oVvZZshsCDmuv4vDjZQw=="

var (
	ErrJWTExpired      = errors.New("jwt expired")
	ErrInvalidAudience = errors.New("invalid audience")
)

// claimsParser defines an interface which makes the creation of the JWT more flexible,
// every type which implements the claimsParser can be added as a JWT claim
type claimsParser interface {
	ClaimsParse() map[string]interface{}
}

// Create actually creates a valid jwt with all necessary standards, and put
// all the claims into the payload. Also, setup jwt standard claims as well.
func Create(ctx context.Context, meta Metadata, claimsContainer ...claimsParser) ([]byte, error) {
	c := jwt.Claims{}
	if c.Set == nil {
		c.Set = make(map[string]interface{})
	}

	for _, claims := range claimsContainer {
		if claims != nil {
			for key, claim := range claims.ClaimsParse() {
				c.Set[key] = claim
			}
		}
	}

	if idUntyped, ok := c.Set["user_id"]; ok {
		if id, ok := idUntyped.(string); ok {
			c.Subject = id
		}
	}
	if tokenUntyped, ok := c.Set["token"]; ok {
		if t, ok := tokenUntyped.(string); ok {
			c.ID = t
		}
	}
	if expirationUntyped, ok := c.Set["token_expires_at"]; ok {
		if tokenExpiration, ok := expirationUntyped.(time.Time); ok {
			c.Expires = jwt.NewNumericTime(tokenExpiration.Round(time.Second))
		}
	}

	c.Audiences = []string{meta.ClientHost}

	c.Issuer = meta.Issuer

	c.NotBefore = jwt.NewNumericTime(time.Now().UTC())
	c.Issued = jwt.NewNumericTime(time.Now().UTC())

	privateKey, err := getPrivateKey(meta)
	if err != nil {
		return nil, err
	}

	return c.ECDSASign(jwt.ES256, privateKey)
}

// Verify starts with the signature verification, then checking the
// jwt standard claims, finally returns the claims.
func Verify(ctx context.Context, meta Metadata, t []byte) (*jwt.Claims, error) {
	privateKey, err := getPrivateKey(meta)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.ECDSACheck(t, &privateKey.PublicKey) // TODO
	if err != nil {
		return nil, err
	}

	if !claims.Valid(time.Now().UTC()) {
		return nil, ErrJWTExpired
	}

	if !claims.AcceptAudience(meta.ClientHost) {
		return nil, ErrInvalidAudience
	}

	return claims, nil
}

func getPrivateKey(meta Metadata) (*ecdsa.PrivateKey, error) {
	var privateKey *ecdsa.PrivateKey
	var err error
	if pkStr := meta.PrivateKey; pkStr != "" {
		privateKey, err = parsePrivateKey(pkStr)
	} else {
		privateKey, err = parsePrivateKey(DefaultPrivateKey)
	}

	return privateKey, err
}

func parsePrivateKey(pkStr string) (*ecdsa.PrivateKey, error) {
	privateKeyParsed, err := base64.StdEncoding.DecodeString(pkStr)
	if err != nil {
		return nil, err
	}
	return x509.ParseECPrivateKey(privateKeyParsed)
}

//nolint:deadcode,unused
func parsePublicKey(pkStr string) (ecdsa.PublicKey, error) {
	return crypto.UnmarshalECPublicKey(pkStr)
}
