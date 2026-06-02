package autenticacao

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/golang-jwt/jwt/v4"
)

func gerarCertificadoTeste(t *testing.T) (privPEM, certPEM string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("erro ao gerar chave: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{Organization: []string{"FAPITEC Test"}},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("erro ao criar certificado: %v", err)
	}

	certPEM = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}))

	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("erro ao serializar chave privada: %v", err)
	}
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}))

	return
}

func assinarJWT(t *testing.T, privPEM string, claims jwt.MapClaims) string {
	t.Helper()

	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		t.Fatal("falha ao decodificar pem da chave privada")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("erro ao parsear chave privada: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("erro ao assinar JWT: %v", err)
	}

	return signed
}

func TestNovoAdaptadorCasdoor(t *testing.T) {
	a := NovoAdaptadorCasdoor("http://localhost:8000", "test-client", "test-secret", "test-cert", "fapitec", "fapitec")
	if a == nil {
		t.Fatal("adapter nao deve ser nil")
	}
	if a.owner != "fapitec" {
		t.Errorf("owner esperado fapitec, got %s", a.owner)
	}
}

func TestGerarURLDeAutorizacao(t *testing.T) {
	a := NovoAdaptadorCasdoor("http://localhost:8000", "client-id", "secret", "cert", "fapitec", "fapitec")
	url := a.GerarURLDeAutorizacao("http://localhost:3000/auth/callback", "test-state")

	if !strings.Contains(url, "http://localhost:8000/api/login/oauth/authorize") {
		t.Errorf("URL deve conter o endpoint de autorizacao: %s", url)
	}
	if !strings.Contains(url, "client_id=client-id") {
		t.Errorf("URL deve conter client_id: %s", url)
	}
	if !strings.Contains(url, "redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fauth%2Fcallback") {
		t.Errorf("URL deve conter redirect_uri codificado: %s", url)
	}
	if !strings.Contains(url, "state=test-state") {
		t.Errorf("URL deve conter state: %s", url)
	}
	if !strings.Contains(url, "openid") {
		t.Errorf("URL deve conter scope openid: %s", url)
	}
}

func TestValidarJWT_TokenValido(t *testing.T) {
	privPEM, certPEM := gerarCertificadoTeste(t)

	a := NovoAdaptadorCasdoor("http://localhost:8000", "client-id", "secret", certPEM, "fapitec", "fapitec")

	claims := jwt.MapClaims{
		"sub":   "123",
		"name":  "11122233344",
		"email": "teste@fapitec.se.gov.br",
		"type":  "proponente",
		"exp":   float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	tokenStr := assinarJWT(t, privPEM, claims)

	claimsRetornados, err := a.ValidarJWT(tokenStr)
	if err != nil {
		t.Fatalf("ValidarJWT retornou erro para token valido: %v", err)
	}

	if claimsRetornados == nil {
		t.Fatal("claims nao deve ser nil")
	}

	if claimsRetornados.Email != "teste@fapitec.se.gov.br" {
		t.Errorf("email esperado teste@fapitec.se.gov.br, got %s", claimsRetornados.Email)
	}
}

func TestValidarJWT_TokenInvalido(t *testing.T) {
	_, certPEM := gerarCertificadoTeste(t)

	a := NovoAdaptadorCasdoor("http://localhost:8000", "client-id", "secret", certPEM, "fapitec", "fapitec")

	_, err := a.ValidarJWT("eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjMifQ.invalido")
	if err == nil {
		t.Fatal("ValidarJWT deve retornar erro para token invalido")
	}
}

func TestValidarJWT_TokenComAlgoritmoDiferente(t *testing.T) {
	_, certPEM := gerarCertificadoTeste(t)

	a := NovoAdaptadorCasdoor("http://localhost:8000", "client-id", "secret", certPEM, "fapitec", "fapitec")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "123",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	})

	tokenStr, _ := token.SignedString([]byte("secret"))
	_, err := a.ValidarJWT(tokenStr)
	if err == nil {
		t.Fatal("ValidarJWT deve rejeitar token com algoritmo HMAC (nao suportado)")
	}
}

func TestValidarJWT_TokenExpirado(t *testing.T) {
	privPEM, certPEM := gerarCertificadoTeste(t)

	a := NovoAdaptadorCasdoor("http://localhost:8000", "client-id", "secret", certPEM, "fapitec", "fapitec")

	claims := jwt.MapClaims{
		"sub": "123",
		"exp": float64(time.Now().Add(-1 * time.Hour).Unix()),
	}

	tokenStr := assinarJWT(t, privPEM, claims)

	_, err := a.ValidarJWT(tokenStr)
	if err == nil {
		t.Fatal("ValidarJWT deve retornar erro para token expirado")
	}
}

func TestVerificarPermissao(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "/api/enforce") {
			resp := casdoorsdk.Response{
				Status: "ok",
				Data:   []interface{}{true},
				Msg:    "",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	casdoorsdk.SetHttpClient(mockServer.Client())
	t.Cleanup(func() { casdoorsdk.SetHttpClient(http.DefaultClient) })

	a := NovoAdaptadorCasdoor(mockServer.URL, "client-id", "secret", "cert", "fapitec", "fapitec")

	permitido, err := a.VerificarPermissao(nil, "usuario1", "proponente", "editais", "visualizar")
	if err != nil {
		t.Fatalf("VerificarPermissao retornou erro: %v", err)
	}
	if !permitido {
		t.Fatal("VerificarPermissao deve retornar true para permissao concedida")
	}
}

func TestVerificarPermissaoNegada(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "/api/enforce") {
			resp := casdoorsdk.Response{
				Status: "ok",
				Data:   []interface{}{false},
				Msg:    "",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	casdoorsdk.SetHttpClient(mockServer.Client())
	t.Cleanup(func() { casdoorsdk.SetHttpClient(http.DefaultClient) })

	a := NovoAdaptadorCasdoor(mockServer.URL, "client-id", "secret", "cert", "fapitec", "fapitec")

	permitido, err := a.VerificarPermissao(nil, "usuario1", "proponente", "financeiro", "excluir")
	if err != nil {
		t.Fatalf("VerificarPermissao retornou erro: %v", err)
	}
	if permitido {
		t.Fatal("VerificarPermissao deve retornar false para permissao negada")
	}
}
