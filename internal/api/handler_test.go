package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	shopcontext "merch-shop/internal/api/context"
	"merch-shop/internal/api/mocks"
	"merch-shop/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := shopcontext.WithUserID(r.Context(), 1)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestBuyMerch(t *testing.T) {
	t.Parallel()

	// Создаём мок UseCase
	mockUseCase := new(mocks.UseCase)
	handler := &HTTPHandler{useCase: mockUseCase}

	tests := []struct {
		name           string
		item           string
		authHeader     string
		mockUseCaseErr error
		expectedStatus int
	}{
		{
			name:           "Unauthorized request (no token)",
			item:           "hat",
			authHeader:     "",
			mockUseCaseErr: nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Database error",
			item:           "bag",
			authHeader:     "Bearer valid_token",
			mockUseCaseErr: errors.New("DB error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUseCase.ExpectedCalls = nil

			mockUseCase.On("BuyMerch", mock.Anything, mock.Anything, tt.item).
				Return(tt.mockUseCaseErr).Maybe()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/buy/%s", tt.item), nil)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			r := chi.NewRouter()
			r.With(mockJWTMiddleware).Get("/buy/{item}", handler.BuyMerch)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestSendCoin(t *testing.T) {
	t.Parallel()

	mockUseCase := new(mocks.UseCase)
	mockValidator := validator.New()

	handler := &HTTPHandler{
		useCase:  mockUseCase,
		validate: mockValidator,
	}

	tests := []struct {
		name           string
		requestBody    any
		authHeader     string
		mockUseCaseErr error
		expectedStatus int
		expectMockCall bool
	}{
		{
			name:           "Invalid JSON",
			requestBody:    `{invalid_json`,
			authHeader:     "Bearer valid_token",
			mockUseCaseErr: nil,
			expectedStatus: http.StatusBadRequest,
			expectMockCall: false,
		},
		{
			name: "Unauthorized request (no token)",
			requestBody: domain.SendCoinRequest{
				ToUser: "recipient",
				Amount: 10,
			},
			authHeader:     "",
			mockUseCaseErr: nil,
			expectedStatus: http.StatusUnauthorized,
			expectMockCall: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUseCase.ExpectedCalls = nil

			var reqBody []byte
			var err error

			switch body := tt.requestBody.(type) {
			case domain.SendCoinRequest:
				reqBody, err = json.Marshal(body)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			case string:
				reqBody = []byte(body)
			}

			if tt.expectMockCall {
				if validReq, ok := tt.requestBody.(domain.SendCoinRequest); ok {
					mockUseCase.On("SendCoin", mock.Anything, uint64(1), validReq).
						Return(tt.mockUseCaseErr).Once()
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			r := chi.NewRouter()
			r.With(mockJWTMiddleware).Post("/sendCoin", handler.SendCoin)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectMockCall {
				mockUseCase.AssertExpectations(t)
			}
		})
	}
}
