package usecase

import (
	"context"
	"merch-shop/internal/domain"
)

func (u *UseCase) GetInfo(ctx context.Context, userID uint64) (domain.Info, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return domain.Info{}, err
	}

	inventory, err := u.repo.GetUserInventory(ctx, userID)
	if err != nil {
		return domain.Info{}, err
	}

	history, err := u.repo.GetUserTransactions(ctx, userID)
	if err != nil {
		return domain.Info{}, err
	}

	return domain.Info{
		UserID:      userID,
		Coins:       user.Coins,
		Inventory:   inventory,
		CoinHistory: history,
	}, nil
}
