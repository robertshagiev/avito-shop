package apierror

import (
	"errors"
	"merch-shop/internal/usecase"
	"net/http"
)

var (
	ErrParsingBody           = errors.New("failed to parse the request body")
	ErrValidatingBody        = errors.New("failed to validate the structure of request body")
	ErrInvalidToken          = errors.New("invalid token")
	ErrParsingToken          = errors.New("failed to parse the JWT token")
	ErrInvalidAuthHeader     = errors.New("the Authorization header is empty or does not contain Bearer token")
	ErrAuthorizationRequired = errors.New("authorization required")
	ErrInvalidRequest        = errors.New("invalid request")
)

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, err error) {
	e := FromError(err)
	RenderJSONWithStatus(w, JSON{"error": e.Message}, e.Code)
}

func FromError(err error) *Err {
	var (
		code    int
		message string
	)

	switch {
	default:
		code = http.StatusInternalServerError
		message = "internal server error"

	case errors.Is(err, ErrParsingBody):
		code = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, ErrValidatingBody):
		code = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, usecase.ErrUnauthorized):
		code = http.StatusUnauthorized
		message = err.Error()
	case errors.Is(err, ErrInvalidToken):
		code = http.StatusUnauthorized
		message = err.Error()
	case errors.Is(err, ErrInvalidAuthHeader):
		code = http.StatusUnauthorized
		message = err.Error()
	case errors.Is(err, ErrParsingToken):
		code = http.StatusUnauthorized
		message = err.Error()
	case errors.Is(err, ErrAuthorizationRequired):
		code = http.StatusUnauthorized
		message = err.Error()
	case errors.Is(err, usecase.ErrNoCoins):
		code = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, usecase.ErrSendCoin):
		code = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, ErrInvalidRequest):
		code = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, usecase.PasswordNotValid):
		code = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, usecase.UsernameNotValid):
		code = http.StatusBadRequest
		message = err.Error()
	}

	return &Err{
		Code:    code,
		Message: message,
	}
}
