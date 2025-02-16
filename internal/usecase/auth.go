package usecase

import (
	"context"
	"errors"
	"fmt"
	"merch-shop/internal/domain"
)

func (u *UseCase) Login(ctx context.Context, creds domain.Credentials) (string, error) {
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
