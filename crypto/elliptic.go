package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"strconv"
)

// MarshalECPublicKey accept an ecdsa.PublicKey and marshal it to a compressed shareable format
func MarshalECPublicKey(publicKey ecdsa.PublicKey) string {
	return base64.StdEncoding.EncodeToString(
		append(
			[]byte(strconv.FormatInt(int64(publicKey.Params().BitSize), 10)),
			elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)...,
		),
	)
}

// UnmarshalECPublicKey accept a compressed format and parse to an ecdsa.PublicKey
func UnmarshalECPublicKey(compressed string) (ecdsa.PublicKey, error) {
	compressedPublicKey, err := base64.StdEncoding.DecodeString(compressed)
	if err != nil {
		return ecdsa.PublicKey{}, err
	}

	bitSize := compressedPublicKey[:3]
	data := compressedPublicKey[3:]

	var curve elliptic.Curve
	switch string(bitSize) {
	case "256":
		curve = elliptic.P256()
	case "384":
		curve = elliptic.P384()
	case "521":
		curve = elliptic.P521()
	default:
		return ecdsa.PublicKey{}, errors.New("invalid bit size")
	}
	var publicKey ecdsa.PublicKey
	publicKey.Curve = curve
	publicKey.X, publicKey.Y = elliptic.Unmarshal(curve, data)

	return publicKey, nil
}
