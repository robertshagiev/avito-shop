package usecase

import (
	"context"
	"errors"
	"merch-shop/internal/domain"
	"merch-shop/internal/usecase/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseCase_GetInfo(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name           string
		mockUser       domain.User
		mockInventory  []domain.Inventory
		mockHistory    domain.CoinHistory
		mockUserErr    error
		mockInvErr     error
		mockHistErr    error
		expectErr      bool
		expectedResult domain.Info
	}{
		{
			name: "Successful retrieval",
			mockUser: domain.User{
				ID:    1,
				Coins: 100,
			},
			mockInventory: []domain.Inventory{
				{Type: "pen", Quantity: 2},
				{Type: "cup", Quantity: 5},
			},
			mockHistory: domain.CoinHistory{
				Received: []domain.CoinTransaction{
					{UserName: "bob", Amount: 10},
				},
				Sent: []domain.CoinTransaction{
					{UserName: "bob", Amount: 5},
				},
			},
			expectErr: false,
			expectedResult: domain.Info{
				UserID: 1,
				Coins:  100,
				Inventory: []domain.Inventory{
					{Type: "pen", Quantity: 2},
					{Type: "cup", Quantity: 5},
				},
				CoinHistory: domain.CoinHistory{
					Received: []domain.CoinTransaction{
						{UserName: "bob", Amount: 10},
					},
					Sent: []domain.CoinTransaction{
						{UserName: "bob", Amount: 5},
					},
				},
			},
		},
		{
			name:        "Error fetching user",
			mockUserErr: errors.New("DB error"),
			expectErr:   true,
		},
		{
			name: "Error fetching inventory",
			mockUser: domain.User{
				ID:    1,
				Coins: 100,
			},
			mockInvErr: errors.New("DB error"),
			expectErr:  true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repository)
			useCase := &UseCase{repo: mockRepo}
			ctx := context.Background()
			userID := uint64(1)

			mockRepo.ExpectedCalls = nil

			mockRepo.On("GetUserByID", ctx, userID).
				Return(tt.mockUser, tt.mockUserErr).Once()

			if tt.mockUserErr != nil {
				info, err := useCase.GetInfo(ctx, userID)
				assert.ErrorContains(t, err, "DB error")
				assert.Empty(t, info)
				mockRepo.AssertExpectations(t)
				return
			}

			mockRepo.On("GetUserInventory", ctx, userID).
				Return(tt.mockInventory, tt.mockInvErr).Once()

			if tt.mockInvErr != nil {
				info, err := useCase.GetInfo(ctx, userID)
				assert.ErrorContains(t, err, "DB error")
				assert.Empty(t, info)
				mockRepo.AssertExpectations(t)
				return
			}

			mockRepo.On("GetUserTransactions", ctx, userID).
				Return(tt.mockHistory, tt.mockHistErr).Once()

			if tt.mockHistErr != nil {
				info, err := useCase.GetInfo(ctx, userID)
				assert.ErrorContains(t, err, "DB error")
				assert.Empty(t, info)
				mockRepo.AssertExpectations(t)
				return
			}

			info, err := useCase.GetInfo(ctx, userID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, info)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, info)

			mockRepo.AssertExpectations(t)
		})
	}
}
