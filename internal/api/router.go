package api

import (
	"crypto/rsa"
	"github.com/go-chi/chi/v5"
	"merch-shop/internal/api/middlewares"
	"net/http"
)

func NewRouter(
	handler *HTTPHandler,
	publicKey *rsa.PublicKey,
) (http.Handler, error) {
	r := chi.NewRouter()

	mid := middlewares.New(publicKey)

	r.Route("/api", func(r chi.Router) {

		r.Post("/auth", handler.Auth)
		r.With(mid.JWTToken).Get("/info", handler.Info)
		r.With(mid.JWTToken).Post("/sendCoin", handler.SendCoin)

	})

	return r, nil
}
