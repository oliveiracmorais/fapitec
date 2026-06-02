package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

type mockValidator struct {
	claims *casdoorsdk.Claims
	err    error
}

func (m *mockValidator) ValidarJWT(token string) (*casdoorsdk.Claims, error) {
	return m.claims, m.err
}

func TestAutenticacaoMiddleware_TokenAusente(t *testing.T) {
	validator := &mockValidator{claims: nil, err: nil}
	middleware := AutenticacaoMiddleware(validator)
	var capturedReq *http.Request
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status esperado 401, got %d", rec.Code)
	}
	if capturedReq != nil {
		t.Error("handler nao deve ser chamado quando token ausente")
	}
}

func TestAutenticacaoMiddleware_FormatoInvalido(t *testing.T) {
	validator := &mockValidator{claims: nil, err: nil}
	middleware := AutenticacaoMiddleware(validator)
	var capturedReq *http.Request
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status esperado 401, got %d", rec.Code)
	}
	if capturedReq != nil {
		t.Error("handler nao deve ser chamado quando formato invalido")
	}
}

func TestAutenticacaoMiddleware_TokenInvalido(t *testing.T) {
	validator := &mockValidator{claims: nil, err: http.ErrAbortHandler}
	middleware := AutenticacaoMiddleware(validator)
	var capturedReq *http.Request
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
	req.Header.Set("Authorization", "Bearer token.invalido.xyz")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status esperado 401, got %d", rec.Code)
	}
	if capturedReq != nil {
		t.Error("handler nao deve ser chamado quando token invalido")
	}
}

func TestAutenticacaoMiddleware_TokenValido(t *testing.T) {
	fakeClaims := &casdoorsdk.Claims{}
	validator := &mockValidator{claims: fakeClaims, err: nil}
	middleware := AutenticacaoMiddleware(validator)
	var capturedReq *http.Request
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
	req.Header.Set("Authorization", "Bearer token.valido.xyz")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status esperado 200, got %d", rec.Code)
	}

	if capturedReq == nil {
		t.Fatal("handler deve ser chamado quando token valido")
	}
	claims := capturedReq.Context().Value(UsuarioContextKey)
	if claims == nil {
		t.Error("claims devem estar presentes no contexto")
	}
	if claims != fakeClaims {
		t.Errorf("claims no contexto incorretos: esperado %v, got %v", fakeClaims, claims)
	}
}
