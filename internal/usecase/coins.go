package usecase

import (
	"context"
	"fmt"
	"merch-shop/internal/domain"
)

func (u *UseCase) SendCoin(ctx context.Context, fromUserID uint64, req domain.SendCoinRequest) error {
	user, err := u.repo.GetUserByUsername(ctx, req.ToUser)
	if err != nil {
		return fmt.Errorf("ошибка поиска пользователя %s: %w", req.ToUser, err)
	}

	if fromUserID == user.ID {
		return ErrSendCoin
	}

	if user.Coins < req.Amount {
		return ErrNoCoins
	}

	err = u.repo.TransferCoins(ctx, fromUserID, user.ID, req.Amount)
	if err != nil {
		return fmt.Errorf("ошибка перевода монет: %w", err)
	}

	return nil
}
