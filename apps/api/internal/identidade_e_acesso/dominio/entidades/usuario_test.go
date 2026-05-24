package entidades

import (
	"testing"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
)

func TestUsuarioBloqueioApos5Tentativas(t *testing.T) {
	u := usuarioTeste()
	for i := 0; i < 5; i++ {
		u.RegistrarTentativaFalha(15 * time.Minute)
	}
	if !u.EstaBloqueado() {
		t.Fatal("Usuario deveria estar bloqueado apos 5 tentativas")
	}
}

func TestUsuarioNaoBloqueiaApos3Tentativas(t *testing.T) {
	u := usuarioTeste()
	for i := 0; i < 3; i++ {
		u.RegistrarTentativaFalha(15 * time.Minute)
	}
	if u.EstaBloqueado() {
		t.Fatal("Usuario nao deveria estar bloqueado com 3 tentativas")
	}
}

func TestUsuarioResetarTentativas(t *testing.T) {
	u := usuarioTeste()
	for i := 0; i < 5; i++ {
		u.RegistrarTentativaFalha(15 * time.Minute)
	}
	if !u.EstaBloqueado() {
		t.Fatal("Usuario deveria estar bloqueado")
	}
	u.ResetarTentativas()
	if u.EstaBloqueado() {
		t.Fatal("Usuario nao deveria estar bloqueado apos reset")
	}
	if u.Tentativas != 0 {
		t.Fatalf("Tentativas deveria ser 0, got %d", u.Tentativas)
	}
}

func TestUsuarioBloqueioExpirado(t *testing.T) {
	u := usuarioTeste()
	for i := 0; i < 5; i++ {
		u.RegistrarTentativaFalha(15 * time.Minute)
	}
	passado := time.Now().Add(-16 * time.Minute)
	u.BloqueadoAte = &passado
	if u.EstaBloqueado() {
		t.Fatal("Usuario nao deveria estar bloqueado apos expiracao do prazo")
	}
}

func usuarioTeste() Usuario {
	return Usuario{
		Nome:      "Teste",
		CPF:       "52998224725",
		Email:     objetos_de_valor.NewEmailSeguro("teste@example.com"),
		SenhaHash: objetos_de_valor.NovaSenhaHash("hash_teste"),
	}
}
