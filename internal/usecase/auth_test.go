package usecase

import (
	"context"
	"errors"
	"merch-shop/internal/domain"
	"merch-shop/internal/usecase/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Login(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedToken := "mocked_token"

	for _, tt := range []struct {
		name          string
		creds         domain.Credentials
		mockUser      domain.User
		mockUserErr   error
		mockUserID    uint64
		mockUserIDErr error
		expectToken   string
		expectErr     bool
	}{
		{
			name: "Successful authentication",
			creds: domain.Credentials{
				Username: "testuser",
				Password: "TestPassword1",
			},
			mockUser: domain.User{
				ID: 1,
				Credentials: domain.Credentials{
					Username: "testuser",
					Password: "TestPassword1",
				},
			},
			mockUserErr:   nil,
			mockUserID:    1,
			mockUserIDErr: nil,
			expectToken:   expectedToken,
			expectErr:     false,
		},
		{
			name: "Creating a new user if it is not in the database",
			creds: domain.Credentials{
				Username: "newuser",
				Password: "NewPass123",
			},
			mockUser:      domain.User{},
			mockUserErr:   ErrNotFound,
			mockUserID:    2,
			mockUserIDErr: nil,
			expectToken:   expectedToken,
			expectErr:     false,
		},
		{
			name: "Error creating new user",
			creds: domain.Credentials{
				Username: "newuser",
				Password: "NewPass123",
			},
			mockUser:      domain.User{},
			mockUserErr:   ErrNotFound,
			mockUserID:    0,
			mockUserIDErr: errors.New("DB error"),
			expectToken:   "",
			expectErr:     true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repository)
			mockAuth := new(mocks.Auth)
			useCase := &UseCase{repo: mockRepo, auth: mockAuth}

			mockRepo.Test(t)
			mockAuth.Test(t)

			mockRepo.On("GetUserByUsername", ctx, tt.creds.Username).
				Return(tt.mockUser, tt.mockUserErr).Once()

			if errors.Is(tt.mockUserErr, ErrNotFound) {
				mockRepo.On("CreateUser", ctx, mock.Anything).
					Return(tt.mockUserID, tt.mockUserIDErr).Once()
			}

			if tt.mockUserIDErr == nil {
				mockAuth.On("NewAccessToken", tt.mockUserID).
					Return(expectedToken, nil).Once()
			}

			token, err := useCase.Login(ctx, tt.creds)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectToken, token)
		})
	}
}

func TestUseCase_CheckCredentials(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name       string
		creds      domain.Credentials
		mockReturn domain.User
		mockError  error
		expectedID uint64
		expectErr  bool
	}{
		{
			name: "Valid credentials",
			creds: domain.Credentials{
				Username: "testuser",
				Password: "testpassword",
			},
			mockReturn: domain.User{
				ID: 1,
				Credentials: domain.Credentials{
					Username: "testuser",
					Password: "testpassword",
				},
			},
			expectedID: 1,
			expectErr:  false,
		},
		{
			name: "User not found",
			creds: domain.Credentials{
				Username: "testuser",
				Password: "testpassword",
			},
			mockReturn: domain.User{},
			mockError:  ErrNotFound,
			expectedID: 0,
			expectErr:  true,
		},
		{
			name: "Database error",
			creds: domain.Credentials{
				Username: "testuser",
				Password: "testpassword",
			},
			mockReturn: domain.User{},
			mockError:  errors.New("DB error"),
			expectedID: 0,
			expectErr:  true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repository)
			useCase := &UseCase{repo: mockRepo}

			ctx := context.Background()

			mockRepo.ExpectedCalls = nil

			mockRepo.On("GetUserByUsername", ctx, tt.creds.Username).
				Return(tt.mockReturn, tt.mockError).Once()

			userID, err := useCase.CheckCredentials(ctx, tt.creds)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, uint64(0), userID)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedID, userID)

			mockRepo.AssertExpectations(t)
		})
	}
}
