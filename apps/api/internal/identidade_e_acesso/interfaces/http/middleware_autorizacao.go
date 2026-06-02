package http

import (
	"context"
	"net/http"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

type PermissionChecker interface {
	VerificarPermissao(ctx context.Context, usuarioID, perfil, modulo, operacao string) (bool, error)
}

func AutorizacaoMiddleware(checker PermissionChecker, modulo, operacao string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UsuarioContextKey).(*casdoorsdk.Claims)
			if !ok || claims == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"erro":"usuario nao autenticado"}`))
				return
			}

			perfil := claims.Type
			if perfil == "" {
				perfil = "proponente"
			}

			permitido, err := checker.VerificarPermissao(r.Context(), claims.Name, perfil, modulo, operacao)
			if err != nil || !permitido {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"erro":"acesso negado"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
