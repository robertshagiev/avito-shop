package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"merch-shop/internal/usecase"
)

const transferCoins = `
	WITH updated_sender AS (
		UPDATE users
		SET coins = coins - $1
		WHERE id = $2 AND coins >= $1
		RETURNING id
	),
	updated_receiver AS (
		UPDATE users
		SET coins = coins + $1
		WHERE id = $3
		RETURNING id
	)
	INSERT INTO public.transactions (from_user_id, to_user_id, quantity)
	SELECT $2, $3, $1
	WHERE EXISTS (SELECT 1 FROM updated_sender) AND EXISTS (SELECT 1 FROM updated_receiver)
	`

func (r *Repository) TransferCoins(ctx context.Context, fromUserID, toUserID uint64, amount uint64) error {

	_, err := r.db.ExecContext(ctx, transferCoins, amount, fromUserID, toUserID)
	if err != nil {
		return fmt.Errorf("ошибка выполнения перевода монет: %w", err)
	}

	return nil
}

const buyMerchQuery = `
WITH deducted AS (
	UPDATE public.users 
	SET coins = coins - $1
	WHERE id = $2 AND coins >= $1
	RETURNING id
)
INSERT INTO public.inventory (user_id, merch_id, quantity)
SELECT $2, m.id, 1
FROM public.merch m, deducted
WHERE m.name = $3
ON CONFLICT (user_id, merch_id) 
DO UPDATE SET quantity = inventory.quantity + 1;
`

func (r *Repository) BuyMerch(ctx context.Context, userID uint64, itemName string, itemPrice uint64) error {
	result, err := r.db.ExecContext(ctx, buyMerchQuery, itemPrice, userID, itemName)
	if err != nil {
		return fmt.Errorf("ошибка при покупке товара: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке обновления: %w", err)
	}

	if rowsAffected == 0 {
		return usecase.ErrNoCoins
	}

	return nil
}

const getMerchPriceQuery = `SELECT price FROM public.merch WHERE name = $1`

func (r *Repository) GetMerchPrice(ctx context.Context, itemName string) (uint64, error) {
	var price uint64
	err := r.db.QueryRowContext(ctx, getMerchPriceQuery, itemName).Scan(&price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, usecase.ErrNotFound
		}
		return 0, fmt.Errorf("ошибка получения цены товара: %w", err)
	}

	return price, nil
}
