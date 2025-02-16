package api

import (
	"encoding/json"
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
	useCase  *usecase.UseCase
}

func NewHTTPHandler(useCase *usecase.UseCase) *HTTPHandler {
	validate := validator.New()

	return &HTTPHandler{
		validate: validate,
		useCase:  useCase,
	}
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
		slog.Error("Failed to generate access token", "error", err)
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
		slog.Error("Failed to generate access token", "error", err)
		apierror.WriteError(w, err)
		return
	}

	apierror.RenderJSONWithStatus(w, info, http.StatusOK)
}

type sendCoinResponse struct {
	Message string `json:"message"`
}

func (h *HTTPHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	var (
		body domain.SendCoinRequest
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

	fromUserID, ok := shopcontext.UserID(ctx)
	if !ok {
		slog.Error("Failed to get user ID")
		apierror.WriteError(w, apierror.ErrAuthorizationRequired)
		return
	}

	err = h.useCase.SendCoin(ctx, fromUserID, body)
	if err != nil {

		slog.Error("Failed to send coin", "error", err)
		apierror.WriteError(w, err)
		return
	}

	// Отправляем успешный ответ
	apierror.RenderJSONWithStatus(w, sendCoinResponse{Message: "Монеты успешно отправлены"}, http.StatusOK)

}
