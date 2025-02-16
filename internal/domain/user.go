package domain

import (
	"crypto/sha256"
	"encoding/base64"
)

type User struct {
	ID    uint64 `json:"user_id"`
	Coins uint64 `json:"coins"`
	Credentials
}

type Credentials struct {
	Username string   `json:"username" validate:"required"`
	Password Password `json:"password" validate:"required"`
}

type Password string

func (p *Password) Secure() Password {
	return Password(p.hash())
}

func (p *Password) Verify(password Password) bool {
	return p.String() == password.hash()
}

func (p *Password) String() string {
	if p == nil {
		return ""
	}
	return string(*p)
}

func (p *Password) hash() string {
	hasher := sha256.New()
	hasher.Write([]byte(p.String()))
	h := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(h)
}

type Info struct {
	UserID      uint64      `json:"user_id"`
	Coins       uint64      `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Inventory struct {
	Type     string `json:"type"`
	Quantity uint   `json:"quantity"`
}

type CoinHistory struct {
	Received []CoinTransaction `json:"received"`
	Sent     []CoinTransaction `json:"sent"`
}

type CoinTransaction struct {
	UserName string `json:"fromUser,omitempty"`
	Amount   int    `json:"amount"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount uint64 `json:"amount" validate:"required"`
}
