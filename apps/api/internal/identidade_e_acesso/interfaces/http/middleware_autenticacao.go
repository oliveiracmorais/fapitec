package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

type contextKey string

const UsuarioContextKey contextKey = "usuario_casdoor"

type ValidatorJWT interface {
	ValidarJWT(token string) (*casdoorsdk.Claims, error)
}

func AutenticacaoMiddleware(validator ValidatorJWT) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"erro":"token ausente"}`))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"erro":"formato invalido"}`))
				return
			}

			claims, err := validator.ValidarJWT(parts[1])
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"erro":"token invalido"}`))
				return
			}

			ctx := context.WithValue(r.Context(), UsuarioContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
