package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/kelseyhightower/envconfig"
	"log/slog"
	"strings"
)

type Config struct {
	ServerPort  string `envconfig:"SERVER_PORT" default:"8080"`
	DatabaseURL string `envonfig:"DATABASE_URL" required:"true"`
	PrivateKey  string `envconfig:"PRIVATE_KEY" required:"true"`
	PublicKey   string `envconfig:"PUBLIC_KEY" required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		slog.Error("Failed to load config", err)
		return nil, err
	}
	return &cfg, nil
}

func ParsePrivateKey(base64Key string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(base64Key))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 private key: %w", err)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("failed to parse private key as RSA")
	}

	return rsaKey, nil
}

func ParsePublicKey(base64Key string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 public key: %w", err)
	}

	pk, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return pk, nil
}
