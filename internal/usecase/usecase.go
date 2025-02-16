package usecase

import (
	"context"
	"merch-shop/internal/domain"
)

type UseCase struct {
	auth Auth
	repo Repository
}

type Auth interface {
	NewAccessToken(userID uint64) (string, error)
}

type Repository interface {
	CreateUser(ctx context.Context, creds domain.Credentials) (uint64, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserCoin(ctx context.Context, userID uint64) (domain.User, error)
	GetUserInventory(ctx context.Context, userID uint64) ([]domain.Inventory, error)
	GetUserTransactions(ctx context.Context, userID uint64) (domain.CoinHistory, error)
	TransferCoins(ctx context.Context, fromUserID, toUserID uint64, amount uint64) error
}

func New(auth Auth, repo Repository) *UseCase {
	return &UseCase{
		auth: auth,
		repo: repo,
	}
}
