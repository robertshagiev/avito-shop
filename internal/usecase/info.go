package usecase

import (
	"context"
	"merch-shop/internal/domain"
)

func (u *UseCase) GetInfo(ctx context.Context, userID uint64) (domain.Info, error) {
	var (
		info domain.Info
		err  error
	)

	user, err := u.repo.GetUserCoin(ctx, userID)
	if err != nil {
		return domain.Info{}, err
	}

	info.Coins = user.Coins

	info.Inventory, err = u.repo.GetUserInventory(ctx, userID)
	if err != nil {
		return domain.Info{}, err
	}

	info.CoinHistory, err = u.repo.GetUserTransactions(ctx, userID)
	if err != nil {
		return domain.Info{}, err
	}

	return info, nil
}
