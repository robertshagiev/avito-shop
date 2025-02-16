package authorization

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"strconv"
)

type TokenManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *TokenManager {
	return &TokenManager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (m *TokenManager) NewAccessToken(userID uint64) (string, error) {
	claims := jwt.StandardClaims{
		Audience: "client_id",
		Subject:  strconv.Itoa(int(userID)),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := t.SignedString(m.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
