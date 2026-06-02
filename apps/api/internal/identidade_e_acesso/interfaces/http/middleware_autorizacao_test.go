package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/golang-jwt/jwt/v4"
)

type mockPermissionChecker struct {
	permitido bool
	err       error
}

func (m *mockPermissionChecker) VerificarPermissao(ctx context.Context, usuarioID, perfil, modulo, operacao string) (bool, error) {
	return m.permitido, m.err
}

func injetarClaims(req *http.Request, claims *casdoorsdk.Claims) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), UsuarioContextKey, claims))
}

func TestAutorizacaoMiddleware_SemContexto(t *testing.T) {
	checker := &mockPermissionChecker{permitido: true, err: nil}
	middleware := AutorizacaoMiddleware(checker, "editais", "visualizar")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/editais", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status esperado 401, got %d", rec.Code)
	}
}

func TestAutorizacaoMiddleware_PermissaoNegada(t *testing.T) {
	checker := &mockPermissionChecker{permitido: false, err: nil}
	middleware := AutorizacaoMiddleware(checker, "financeiro", "excluir")
	var capturedCode int
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCode = http.StatusOK
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/financeiro", nil)
	req = injetarClaims(req, &casdoorsdk.Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: "usuario1"},
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status esperado 403, got %d", rec.Code)
	}
	if capturedCode == http.StatusOK {
		t.Error("handler nao deve ser chamado quando permissao negada")
	}
}

func TestAutorizacaoMiddleware_PermissaoConcedida(t *testing.T) {
	checker := &mockPermissionChecker{permitido: true, err: nil}
	middleware := AutorizacaoMiddleware(checker, "editais", "visualizar")
	var capturedCode int
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCode = http.StatusOK
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/editais", nil)
	req = injetarClaims(req, &casdoorsdk.Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: "usuario1"},
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status esperado 200, got %d", rec.Code)
	}
	if capturedCode != http.StatusOK {
		t.Error("handler deve ser chamado quando permissao concedida")
	}
}
