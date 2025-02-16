package usecase

import (
	"context"
	"fmt"
	"merch-shop/internal/domain"
)

func (u *UseCase) SendCoin(ctx context.Context, fromUserID uint64, req domain.SendCoinRequest) error {
	fromUser, err := u.repo.GetUserByID(ctx, fromUserID)
	if err != nil {
		return fmt.Errorf("repo.GetUserByID %s: %w", req.ToUser, err)
	}

	if fromUser.Coins < req.Amount {
		return ErrNoCoins
	}

	toUser, err := u.repo.GetUserByUsername(ctx, req.ToUser)
	if err != nil {
		return fmt.Errorf("repo.GetUserByUsername %s: %w", req.ToUser, err)
	}

	if fromUserID == toUser.ID {
		return ErrSendCoin
	}

	err = u.repo.TransferCoins(ctx, fromUserID, toUser.ID, req.Amount)
	if err != nil {
		return fmt.Errorf("repo.TransferCoins: %w", err)
	}

	return nil
}

func (u *UseCase) BuyMerch(ctx context.Context, userID uint64, itemName string) error {
	itemPrice, err := u.repo.GetMerchPrice(ctx, itemName)
	if err != nil {
		return fmt.Errorf("repo.GetMerchPrice: %w", err)
	}

	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("repo.GetUserByID: %w", err)
	}

	if user.Coins < itemPrice {
		return ErrNoCoins
	}

	if err = u.repo.BuyMerch(ctx, userID, itemName, itemPrice); err != nil {
		return fmt.Errorf("repo.BuyMerch: %w", err)
	}

	return nil
}
