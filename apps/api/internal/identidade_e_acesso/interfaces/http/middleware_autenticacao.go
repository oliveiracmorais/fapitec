package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

type contextKey string

const UsuarioContextKey contextKey = "usuario_casdoor"
const AuthOrigemContextKey contextKey = "auth_origem"

type ValidatorJWT interface {
	ValidarJWT(token string) (*casdoorsdk.Claims, error)
}

type internalClaims struct {
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
}

func validarInternalJWT(token string) (*casdoorsdk.Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, http.ErrNoCookie
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var ic internalClaims
	if err := json.Unmarshal(payload, &ic); err != nil {
		return nil, err
	}

	if time.Now().Unix() > ic.Exp {
		return nil, http.ErrNoCookie
	}

	return &casdoorsdk.Claims{
		User: casdoorsdk.User{
			Name:        ic.Sub,
			DisplayName: ic.Name,
			Email:       ic.Email,
			Type:        "proponente",
		},
	}, nil
}

func AutenticacaoMiddleware(validator ValidatorJWT) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token := ""

			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					token = parts[1]
				}
			}

			if token == "" {
				cookie, err := r.Cookie("fapitec_token")
				if err == nil {
					token = cookie.Value
				}
			}

			if token == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"erro":"token ausente"}`))
				return
			}

			claims, err := validator.ValidarJWT(token)
			origem := "casdoor"
			if err != nil {
				claims, err = validarInternalJWT(token)
				origem = "internal"
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"erro":"token invalido"}`))
					return
				}
			}

			ctx := context.WithValue(r.Context(), UsuarioContextKey, claims)
			ctx = context.WithValue(ctx, AuthOrigemContextKey, origem)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
