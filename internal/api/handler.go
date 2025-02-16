package api

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"merch-shop/internal/api/apierror"
	shopcontext "merch-shop/internal/api/context"
	"merch-shop/internal/domain"
	"merch-shop/internal/usecase"
	"net/http"
)

type HTTPHandler struct {
	validate *validator.Validate
	useCase  UseCase
}

func NewHTTPHandler(useCase *usecase.UseCase) *HTTPHandler {
	validate := validator.New()

	return &HTTPHandler{
		validate: validate,
		useCase:  useCase,
	}
}

//go:generate mockery --name=UseCase --output=./mocks --filename=useCase.go --structname=UseCase
type UseCase interface {
	GetInfo(ctx context.Context, userID uint64) (domain.Info, error)
	Login(ctx context.Context, creds domain.Credentials) (string, error)
	CheckCredentials(ctx context.Context, creds domain.Credentials) (uint64, error)
	SendCoin(ctx context.Context, fromUserID uint64, req domain.SendCoinRequest) error
	BuyMerch(ctx context.Context, userID uint64, itemName string) error
}

type authReq struct {
	domain.Credentials
}

type authResp struct {
	Token string `json:"token"`
}

func (h *HTTPHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var (
		body authReq
		err  error
		ctx  = r.Context()
	)

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.WriteError(w, apierror.ErrParsingBody)
		return
	}
	defer r.Body.Close()

	if err = h.validate.Struct(body); err != nil {
		apierror.WriteError(w, apierror.ErrValidatingBody)
		return
	}

	token, err := h.useCase.Login(ctx, body.Credentials)
	if err != nil {
		slog.Error("useCase.Login", "error", err)
		apierror.WriteError(w, err)
		return
	}

	apierror.RenderJSONWithStatus(w, authResp{Token: token}, http.StatusOK)
}

func (h *HTTPHandler) Info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := shopcontext.UserID(ctx)
	if !ok {
		slog.Error("Failed to get user ID")
		apierror.WriteError(w, apierror.ErrAuthorizationRequired)
		return
	}

	info, err := h.useCase.GetInfo(ctx, userID)
	if err != nil {
		slog.Error("useCase.GetInfo", "error", err)
		apierror.WriteError(w, err)
		return
	}

	apierror.RenderJSONWithStatus(w, info, http.StatusOK)
}

func (h *HTTPHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	var (
		body domain.SendCoinRequest
		err  error
		ctx  = r.Context()
	)

	fromUserID, ok := shopcontext.UserID(ctx)
	if !ok {
		slog.Error("Failed to get user ID")
		apierror.WriteError(w, apierror.ErrAuthorizationRequired)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.WriteError(w, apierror.ErrParsingBody)
		return
	}
	defer r.Body.Close()

	if err = h.validate.Struct(body); err != nil {
		apierror.WriteError(w, apierror.ErrValidatingBody)
		return
	}

	if err = h.useCase.SendCoin(ctx, fromUserID, body); err != nil {
		slog.Error("useCase.SendCoin", "error", err)
		apierror.WriteError(w, err)
		return
	}

	apierror.RenderJSONWithStatus(w, apierror.JSON{}, http.StatusOK)

}

func (h *HTTPHandler) BuyMerch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	item := chi.URLParam(r, "item")

	if item == "" {
		apierror.WriteError(w, apierror.ErrInvalidRequest)
		return
	}

	userID, ok := shopcontext.UserID(ctx)
	if !ok {
		slog.Error("Failed to get user ID")
		apierror.WriteError(w, apierror.ErrAuthorizationRequired)
		return
	}

	if err := h.useCase.BuyMerch(ctx, userID, item); err != nil {
		slog.Error("useCase.BuyMerch", "error", err)
		apierror.WriteError(w, err)
		return
	}

	apierror.RenderJSONWithStatus(w, apierror.JSON{}, http.StatusOK)
}
