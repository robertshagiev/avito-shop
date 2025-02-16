package middlewares

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"merch-shop/internal/api/apierror"
	shopcontext "merch-shop/internal/api/context"
	"net/http"
	"strconv"
	"strings"
)

const authHeader = "Authorization"

func (m *Middlewares) JWTToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getTokenString(r)
		if err != nil {
			apierror.WriteError(w, err)
			return
		}

		token, err := getToken(m.publicKey, tokenString)
		if err != nil {
			apierror.WriteError(w, fmt.Errorf("getToken: %w", err))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, ok := claims["sub"]
			if ok && userID != nil {
				userIDStr, isString := userID.(string)
				if !isString {
					apierror.WriteError(w, apierror.ErrInvalidToken)
					return
				}

				userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
				if err != nil {
					apierror.WriteError(w, apierror.ErrInvalidToken)
					return
				}

				ctx := shopcontext.WithUserID(r.Context(), userIDUint)
				r = r.WithContext(ctx)
			}
		} else {
			apierror.WriteError(w, apierror.ErrInvalidToken)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getToken(key *rsa.PublicKey, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, apierror.ErrParsingToken
			}

			return key, nil
		},
	)
	if err != nil {
		return nil, apierror.ErrParsingToken
	}

	if !token.Valid {
		return nil, apierror.ErrInvalidToken
	}

	return token, nil
}

func getTokenString(r *http.Request) (string, error) {
	bearer := r.Header.Get(authHeader)
	if bearer == "" {
		return "", apierror.ErrInvalidAuthHeader

	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(bearer, "Bearer"))
	if tokenString == "" {
		return "", apierror.ErrInvalidAuthHeader
	}

	return tokenString, nil
}
