package middlewares

import "crypto/rsa"

type Middlewares struct {
	publicKey *rsa.PublicKey
}

func New(
	publicKey *rsa.PublicKey,
) *Middlewares {
	return &Middlewares{
		publicKey: publicKey,
	}
}
