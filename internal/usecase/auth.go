package usecase

import (
	"context"
	"errors"
	"fmt"
	"merch-shop/internal/domain"
	"regexp"
)

func (u *UseCase) Login(ctx context.Context, creds domain.Credentials) (string, error) {
	if !validationUsername(creds.Username) {
		return "", UsernameNotValid
	}

	if !validationPassword(creds.Password) {
		return "", PasswordNotValid
	}

	userID, err := u.CheckCredentials(ctx, creds)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return "", err
		}

		creds.Password = creds.Password.Secure()

		userID, err = u.repo.CreateUser(ctx, creds)
		if err != nil {
			return "", err
		}
	}

	return u.auth.NewAccessToken(userID)
}

func (u *UseCase) CheckCredentials(ctx context.Context, creds domain.Credentials) (uint64, error) {
	fmt.Println(creds.Username)
	userInfo, err := u.repo.GetUserByUsername(ctx, creds.Username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("repo.GetUserByUsername error: %w", err)
	}

	if userInfo.Password.Verify(creds.Password) {
		return 0, ErrUnauthorized
	}

	return userInfo.ID, nil
}

func validationPassword(password domain.Password) bool {
	if len(password) < 8 {
		return false
	}

	if len(password) > 64 {
		return false
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(string(password))
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(string(password))
	hasDigit := regexp.MustCompile(`\d`).MatchString(string(password))

	if !hasUppercase || !hasLowercase || !hasDigit {
		return false
	}

	return true
}

func validationUsername(username string) bool {
	loginRegex := regexp.MustCompile(`^[a-zA-Z0-9]{3,15}$`)

	if !loginRegex.MatchString(username) {
		return false
	}

	return true
}
