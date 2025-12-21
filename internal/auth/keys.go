package auth

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

func LoadKeys() error {
	privData, err := os.ReadFile("keys/private.pem")
	if err != nil {
		return err
	}

	pubData, err := os.ReadFile("keys/public.pem")
	if err != nil {
		return err
	}

	PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privData)
	if err != nil {
		return err
	}

	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubData)
	if err != nil {
		return err
	}

	return nil
}
