package casos_de_uso

import (
	"context"
	"strings"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia"
	"golang.org/x/crypto/bcrypt"
)

type hashTeste struct{}

func (h hashTeste) Hash(senha string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.MinCost)
	return string(bytes), err
}

func TestCadastrarUsuarioComSucesso(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	hash := hashTeste{}
	uc := NovoCadastrarUsuario(repo, hash)

	entrada := dto.CadastrarUsuarioEntrada{
		Nome:             "João Silva",
		CPF:              "529.982.247-25",
		Email:            "joao@example.com",
		ConfirmacaoEmail: "joao@example.com",
		Senha:            "Senha@123",
		ConfirmacaoSenha: "Senha@123",
	}

	saida, err := uc.Executar(context.Background(), entrada)
	if err != nil {
		t.Fatalf("Cadastro retornou erro: %v", err)
	}
	if saida.ID == 0 {
		t.Fatal("ID do usuario nao foi preenchido")
	}
	if saida.Nome != "João Silva" {
		t.Fatalf("Nome esperado João Silva, got %s", saida.Nome)
	}
}

func TestCadastrarUsuarioCPFInvalido(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	uc := NovoCadastrarUsuario(repo, hashTeste{})

	entrada := dto.CadastrarUsuarioEntrada{
		Nome:             "João Silva",
		CPF:              "123.456.789-00",
		Email:            "joao@example.com",
		ConfirmacaoEmail: "joao@example.com",
		Senha:            "Senha@123",
		ConfirmacaoSenha: "Senha@123",
	}

	_, err := uc.Executar(context.Background(), entrada)
	if err == nil {
		t.Fatal("CPF invalido deveria retornar erro")
	}
	if !strings.Contains(err.Error(), "CPF") {
		t.Fatalf("Erro deveria mencionar CPF, got: %v", err)
	}
}

func TestCadastrarUsuarioSenhaNaoConfere(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	uc := NovoCadastrarUsuario(repo, hashTeste{})

	entrada := dto.CadastrarUsuarioEntrada{
		Nome:             "João Silva",
		CPF:              "529.982.247-25",
		Email:            "joao@example.com",
		ConfirmacaoEmail: "joao@example.com",
		Senha:            "Senha@123",
		ConfirmacaoSenha: "Senha@456",
	}

	_, err := uc.Executar(context.Background(), entrada)
	if err == nil {
		t.Fatal("senhas diferentes deveriam retornar erro")
	}
}

func TestCadastrarUsuarioCPFDuplicado(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	uc := NovoCadastrarUsuario(repo, hashTeste{})

	entrada := dto.CadastrarUsuarioEntrada{
		Nome:             "João Silva",
		CPF:              "529.982.247-25",
		Email:            "joao@example.com",
		ConfirmacaoEmail: "joao@example.com",
		Senha:            "Senha@123",
		ConfirmacaoSenha: "Senha@123",
	}

	uc.Executar(context.Background(), entrada)

	_, err := uc.Executar(context.Background(), entrada)
	if err == nil {
		t.Fatal("CPF duplicado deveria retornar erro")
	}
}

func TestCadastrarUsuarioComPassaporte(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	uc := NovoCadastrarUsuario(repo, hashTeste{})

	entrada := dto.CadastrarUsuarioEntrada{
		Nome:             "John Doe",
		CPF:              "AB123456",
		Email:            "john@example.com",
		ConfirmacaoEmail: "john@example.com",
		Senha:            "Senha@123",
		ConfirmacaoSenha: "Senha@123",
		Estrangeiro:      true,
	}

	saida, err := uc.Executar(context.Background(), entrada)
	if err != nil {
		t.Fatalf("Cadastro com passaporte retornou erro: %v", err)
	}
	if saida.Documento != "AB123456" {
		t.Fatalf("Passaporte esperado AB123456, got %s", saida.Documento)
	}
	if !saida.Estrangeiro {
		t.Fatal("Esperava estrangeiro=true")
	}
}


