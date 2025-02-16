package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"merch-shop/internal/domain"
	"merch-shop/internal/usecase"
)

const createUser = `INSERT INTO public.users (username, password) VALUES ($1, $2) RETURNING id`

func (r *Repository) CreateUser(ctx context.Context, creds domain.Credentials) (uint64, error) {
	var userID uint64

	err := r.db.QueryRowContext(
		ctx,
		createUser,
		creds.Username,
		creds.Password,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

const getUserByUsername = `SELECT id, username, password  FROM public.users WHERE username = $1`

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var (
		result         domain.User
		storedPassword string
	)

	if err := r.db.QueryRowContext(ctx, getUserByUsername, username).Scan(&result.ID, &result.Credentials.Username, &storedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, usecase.ErrNotFound
		}
		return domain.User{}, err
	}

	result.Credentials.Password = domain.Password(storedPassword)

	return result, nil
}

const getUserByID = `SELECT coins, username FROM public.users WHERE id = $1`

func (r *Repository) GetUserByID(ctx context.Context, userID uint64) (domain.User, error) {
	var result domain.User

	fmt.Println(userID)

	if err := r.db.QueryRowContext(ctx, getUserByID, userID).Scan(&result.Coins, &result.Credentials.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, usecase.ErrNotFound
		}
		return domain.User{}, err
	}

	return result, nil
}

const getUserInventory = `
	SELECT m.name, i.quantity
	FROM public.inventory i
	JOIN public.merch m ON i.merch_id = m.id
	WHERE i.user_id = $1`

func (r *Repository) GetUserInventory(ctx context.Context, userID uint64) ([]domain.Inventory, error) {
	rows, err := r.db.QueryContext(ctx, getUserInventory, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Inventory{}, nil
		}
		return nil, fmt.Errorf("ошибка получения инвентаря: %w", err)
	}
	defer rows.Close()

	inventory := make([]domain.Inventory, 0)
	for rows.Next() {
		var item domain.Inventory
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, fmt.Errorf("ошибка обработки строки: %w", err)
		}
		inventory = append(inventory, item)
	}

	return inventory, nil
}

const getUserTransactions = `
	SELECT u1.username AS from_user, u2.username AS to_user, t.quantity
	FROM public.transactions t
	LEFT JOIN public.users u1 ON t.from_user_id = u1.id
	LEFT JOIN public.users u2 ON t.to_user_id = u2.id
	WHERE t.from_user_id = $1 OR t.to_user_id = $1`

func (r *Repository) GetUserTransactions(ctx context.Context, userID uint64) (domain.CoinHistory, error) {
	rows, err := r.db.QueryContext(ctx, getUserTransactions, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CoinHistory{}, nil
		}
		return domain.CoinHistory{}, fmt.Errorf("ошибка получения истории транзакций: %w", err)
	}
	defer rows.Close()

	history := domain.CoinHistory{
		Received: []domain.CoinTransaction{},
		Sent:     []domain.CoinTransaction{},
	}

	for rows.Next() {
		var (
			fromUser, toUser sql.NullString
			amount           int
		)

		if err := rows.Scan(&fromUser, &toUser, &amount); err != nil {
			return domain.CoinHistory{}, fmt.Errorf("ошибка обработки строки: %w", err)
		}

		if fromUser.Valid && fromUser.String != "" && toUser.String == "" {
			history.Sent = append(history.Sent, domain.CoinTransaction{
				UserName: fromUser.String,
				Amount:   amount,
			})
		} else if toUser.Valid && toUser.String != "" && fromUser.String == "" {
			history.Received = append(history.Received, domain.CoinTransaction{
				UserName: toUser.String,
				Amount:   amount,
			})
		}
	}

	return history, nil
}
