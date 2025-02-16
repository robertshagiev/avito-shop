package usecase

import (
	"context"
	"errors"
	"merch-shop/internal/domain"
	"merch-shop/internal/usecase/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseCase_SendCoin(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name         string
		fromUser     domain.User
		toUser       domain.User
		req          domain.SendCoinRequest
		mockFromErr  error
		mockToErr    error
		mockTransErr error
		expectErr    error
	}{
		{
			name: "Successful transfer",
			fromUser: domain.User{
				ID:    1,
				Coins: 100,
			},
			toUser: domain.User{
				ID: 2,
				Credentials: domain.Credentials{
					Username: "ivanov",
				},
			},
			req: domain.SendCoinRequest{
				ToUser: "ivanov",
				Amount: 50,
			},
			mockFromErr:  nil,
			mockToErr:    nil,
			mockTransErr: nil,
			expectErr:    nil,
		},
		{
			name: "Error fetching sender",
			fromUser: domain.User{
				ID:    1,
				Coins: 100,
			},
			req: domain.SendCoinRequest{
				ToUser: "ivanov",
				Amount: 50,
			},
			mockFromErr:  errors.New("DB error"),
			mockToErr:    nil,
			mockTransErr: nil,
			expectErr:    errors.New("repo.GetUserByID"),
		},
		{
			name: "Insufficient balance",
			fromUser: domain.User{
				ID:    1,
				Coins: 10,
			},
			req: domain.SendCoinRequest{
				ToUser: "ivanov",
				Amount: 50,
			},
			mockFromErr:  nil,
			mockToErr:    nil,
			mockTransErr: nil,
			expectErr:    ErrNoCoins,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repository)
			useCase := &UseCase{repo: mockRepo}

			ctx := context.Background()

			mockRepo.ExpectedCalls = nil

			mockRepo.On("GetUserByID", ctx, tt.fromUser.ID).
				Return(tt.fromUser, tt.mockFromErr).Once()

			if tt.mockFromErr != nil {
				err := useCase.SendCoin(ctx, tt.fromUser.ID, tt.req)
				assert.ErrorContains(t, err, tt.expectErr.Error())
				mockRepo.AssertExpectations(t)
				return
			}

			mockRepo.On("GetUserByUsername", ctx, tt.req.ToUser).
				Return(tt.toUser, tt.mockToErr).Once()

			if tt.mockToErr != nil {
				err := useCase.SendCoin(ctx, tt.fromUser.ID, tt.req)
				assert.ErrorContains(t, err, tt.expectErr.Error())
				mockRepo.AssertExpectations(t)
				return
			}

			if tt.expectErr == nil || tt.expectErr == ErrSendCoin {
				mockRepo.On("TransferCoins", ctx, tt.fromUser.ID, tt.toUser.ID, tt.req.Amount).
					Return(tt.mockTransErr).Once()
			}

			err := useCase.SendCoin(ctx, tt.fromUser.ID, tt.req)

			if tt.expectErr != nil {
				assert.ErrorContains(t, err, tt.expectErr.Error())
				return
			}
			assert.NoError(t, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
