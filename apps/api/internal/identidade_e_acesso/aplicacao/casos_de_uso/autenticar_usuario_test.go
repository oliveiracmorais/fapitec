package casos_de_uso

import (
	"context"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/hash"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia"
)

func setupUsuario(t *testing.T) *persistencia.RepositorioDeUsuarioMemoria {
	t.Helper()
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	hashService := hash.NovoServicoDeHashBcrypt()
	uc := NovoCadastrarUsuario(repo, hashService)

	entrada := dto.CadastrarUsuarioEntrada{
		Nome:             "Maria Souza",
		CPF:              "529.982.247-25",
		Email:            "maria@example.com",
		ConfirmacaoEmail: "maria@example.com",
		Senha:            "Senha@123",
		ConfirmacaoSenha: "Senha@123",
	}

	_, err := uc.Executar(context.Background(), entrada)
	if err != nil {
		t.Fatalf("Setup falhou: %v", err)
	}
	return repo
}

func TestAutenticarComSucesso(t *testing.T) {
	repo := setupUsuario(t)
	hashService := hash.NovoServicoDeHashBcrypt()
	uc := NovoAutenticarUsuario(repo, hashService)

	entrada := dto.AutenticarUsuarioEntrada{
		CPF:   "529.982.247-25",
		Senha: "Senha@123",
	}

	saida, err := uc.Executar(context.Background(), entrada)
	if err != nil {
		t.Fatalf("Autenticacao retornou erro: %v", err)
	}
	if saida.Nome != "Maria Souza" {
		t.Fatalf("Nome esperado Maria Souza, got %s", saida.Nome)
	}
}

func TestAutenticarSenhaErrada(t *testing.T) {
	repo := setupUsuario(t)
	hashService := hash.NovoServicoDeHashBcrypt()
	uc := NovoAutenticarUsuario(repo, hashService)

	entrada := dto.AutenticarUsuarioEntrada{
		CPF:   "529.982.247-25",
		Senha: "SenhaErrada@123",
	}

	_, err := uc.Executar(context.Background(), entrada)
	if err == nil {
		t.Fatal("Senha errada deveria retornar erro")
	}
}

func TestAutenticarBloqueioApos5Falhas(t *testing.T) {
	repo := setupUsuario(t)
	hashService := hash.NovoServicoDeHashBcrypt()
	uc := NovoAutenticarUsuario(repo, hashService)

	entrada := dto.AutenticarUsuarioEntrada{
		CPF:   "529.982.247-25",
		Senha: "SenhaErrada@123",
	}

	for i := 0; i < 5; i++ {
		uc.Executar(context.Background(), entrada)
	}

	_, err := uc.Executar(context.Background(), entrada)
	if err == nil {
		t.Fatal("Conta deveria estar bloqueada apos 5 tentativas falhas")
	}
}

func TestAutenticarCPFInexistente(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	hashService := hash.NovoServicoDeHashBcrypt()
	uc := NovoAutenticarUsuario(repo, hashService)

	entrada := dto.AutenticarUsuarioEntrada{
		CPF:   "529.982.247-25",
		Senha: "Senha@123",
	}

	_, err := uc.Executar(context.Background(), entrada)
	if err == nil {
		t.Fatal("CPF inexistente deveria retornar erro")
	}
}
