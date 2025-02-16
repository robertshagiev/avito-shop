package usecase

import (
	"context"
	"merch-shop/internal/domain"
)

type UseCase struct {
	auth Auth
	repo Repository
}

//go:generate mockery --name=Auth --output=./mocks --filename=auth.go --structname=Auth
type Auth interface {
	NewAccessToken(userID uint64) (string, error)
}

//go:generate mockery --name=Repository --output=./mocks --filename=repository.go --structname=Repository
type Repository interface {
	CreateUser(ctx context.Context, creds domain.Credentials) (uint64, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByID(ctx context.Context, userID uint64) (domain.User, error)
	GetUserInventory(ctx context.Context, userID uint64) ([]domain.Inventory, error)
	GetUserTransactions(ctx context.Context, userID uint64) (domain.CoinHistory, error)
	TransferCoins(ctx context.Context, fromUserID, toUserID uint64, amount uint64) error
	BuyMerch(ctx context.Context, userID uint64, itemName string, itemPrice uint64) error
	GetMerchPrice(ctx context.Context, itemName string) (uint64, error)
}

func New(auth Auth, repo Repository) *UseCase {
	return &UseCase{
		auth: auth,
		repo: repo,
	}
}
