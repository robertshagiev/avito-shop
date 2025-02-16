package usecase

import "errors"

var (
	ErrUnauthorized  = errors.New("invalid username or password")
	ErrNotFound      = errors.New("item not found")
	ErrNoCoins       = errors.New("have not coins")
	ErrSendCoin      = errors.New("can't send coins to yourself")
	PasswordNotValid = errors.New("password not valid")
	UsernameNotValid = errors.New("username not valid")
)
