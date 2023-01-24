package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"os"
)

func GenerateCurve() (string, error) {
	var curveLength = "256"
	if cl := os.Getenv("CURVE_LENGTH"); cl != "" {
		curveLength = cl
	}

	var ellipticCurve elliptic.Curve
	switch curveLength {
	case "224":
		ellipticCurve = elliptic.P224()
	case "256":
		ellipticCurve = elliptic.P256()
	case "384":
		ellipticCurve = elliptic.P384()
	case "521":
		ellipticCurve = elliptic.P521()
	default:
		return "", errors.New("invalid curve length [224, 256, 384, 521]")
	}

	privateKey, err := ecdsa.GenerateKey(ellipticCurve, rand.Reader)
	if err != nil {
		return "", err
	}
	marshalledPrivateKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(marshalledPrivateKey), nil
}
