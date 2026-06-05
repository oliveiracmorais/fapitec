package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"erro": msg})
}

func setupCasdoorMux() *http.ServeMux {
	mux := http.NewServeMux()
	goneHandler := func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, `{"erro":"autenticacao propria desabilitada — use /api/v1/auth/login com Casdoor"}`, http.StatusGone)
	}
	mux.HandleFunc("POST /api/v1/cadastro", goneHandler)
	mux.HandleFunc("POST /api/v1/register", goneHandler)
	mux.HandleFunc("POST /api/v1/login", goneHandler)
	mux.HandleFunc("POST /api/v1/solicitar-redefinicao-senha", goneHandler)
	mux.HandleFunc("POST /api/v1/reset-password", goneHandler)
	mux.HandleFunc("POST /api/v1/redefinir-senha", goneHandler)
	return mux
}

func TestEndpoint410Gone_ModoCasdoor(t *testing.T) {
	mux := setupCasdoorMux()
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "cadastro", method: "POST", path: "/api/v1/cadastro"},
		{name: "register", method: "POST", path: "/api/v1/register"},
		{name: "login", method: "POST", path: "/api/v1/login"},
		{name: "solicitar-redefinicao-senha", method: "POST", path: "/api/v1/solicitar-redefinicao-senha"},
		{name: "reset-password", method: "POST", path: "/api/v1/reset-password"},
		{name: "redefinir-senha", method: "POST", path: "/api/v1/redefinir-senha"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusGone {
				t.Errorf("%s %s: esperado 410, got %d", tt.method, tt.path, rec.Code)
			}

			ct := rec.Header().Get("Content-Type")
			if !strings.HasPrefix(ct, "application/json") {
				t.Errorf("Content-Type esperado application/json, got %s", ct)
			}

			var body map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("resposta deve ser JSON valido: %v", err)
			}
			if body["erro"] == "" {
				t.Error("resposta deve conter campo 'erro'")
			}
			if !strings.Contains(body["erro"], "desabilitada") {
				t.Errorf("mensagem de erro deve mencionar desabilitacao, got: %s", body["erro"])
			}
		})
	}
}

func TestEndpoint410Gone_NaoRetorna410EmInternal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/api/v1/login", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code == http.StatusGone {
		t.Error("endpoint nao deve retornar 410 em modo internal")
	}
}

func TestAutenticacaoMiddleware_TokenAusenteResiliencia(t *testing.T) {
	validator := &mockValidator{claims: nil, err: nil}
	middleware := AutenticacaoMiddleware(validator)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401 para token ausente, got %d", rec.Code)
	}
}

func TestAutenticacaoMiddleware_TokenExpirado(t *testing.T) {
	validator := &mockValidator{claims: nil, err: http.ErrAbortHandler}
	middleware := AutenticacaoMiddleware(validator)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
	req.Header.Set("Authorization", "Bearer token.expirado.xyz")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401 para token expirado/invalido, got %d", rec.Code)
	}
}

func TestAutenticacaoMiddleware_TokenAssinaturaInvalida(t *testing.T) {
	validator := &mockValidator{claims: nil, err: http.ErrAbortHandler}
	middleware := AutenticacaoMiddleware(validator)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJFUzI1NiJ9.eyJzdWIiOiIxMjMifQ.assinatura-trocada")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401 para token com assinatura invalida, got %d", rec.Code)
	}
}

func TestAutenticacaoMiddleware_SemHeaderAuthorization(t *testing.T) {
	validator := &mockValidator{claims: nil, err: nil}
	middleware := AutenticacaoMiddleware(validator)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401 sem Authorization header, got %d", rec.Code)
	}

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	if body["erro"] != "token ausente" {
		t.Errorf("mensagem esperada 'token ausente', got '%s'", body["erro"])
	}
}

func TestAutenticacaoMiddleware_FormatoBearerInvalido(t *testing.T) {
	validator := &mockValidator{claims: nil, err: nil}
	middleware := AutenticacaoMiddleware(validator)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []string{
		"Bearer",
		"Basic token123",
		"",
	}
	for _, tt := range tests {
		t.Run("formato_"+tt, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/protegido", nil)
			if tt != "" {
				req.Header.Set("Authorization", tt)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("esperado 401 para Authorization '%s', got %d", tt, rec.Code)
			}
		})
	}
}

func TestAutorizacaoMiddleware_SemClaimsNoContexto(t *testing.T) {
	checker := &mockPermissionChecker{permitido: true, err: nil}
	middleware := AutorizacaoMiddleware(checker, "editais", "visualizar")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/editais", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401 sem claims, got %d", rec.Code)
	}
}
