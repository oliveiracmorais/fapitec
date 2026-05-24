package casos_de_uso

import (
	"context"
	"testing"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia"
	"golang.org/x/crypto/bcrypt"
)

type hashRedefinicaoTeste struct{}

func (h hashRedefinicaoTeste) Hash(senha string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.MinCost)
	return string(bytes), err
}

type emailServiceTeste struct {
	ultimoEmail string
	ultimoToken string
}

func (e *emailServiceTeste) EnviarRedefinicaoSenha(ctx context.Context, email, token string) error {
	e.ultimoEmail = email
	e.ultimoToken = token
	return nil
}

func TestSolicitarRedefinicaoSenhaParaEmailExistente(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	hash := hashTeste{}
	emailSvc := &emailServiceTeste{}

	cadastrar := NovoCadastrarUsuario(repo, hash)
	cadastrar.Executar(context.Background(), dto.CadastrarUsuarioEntrada{
		Nome: "João Silva", CPF: "529.982.247-25",
		Email: "joao@example.com", ConfirmacaoEmail: "joao@example.com",
		Senha: "Senha@123", ConfirmacaoSenha: "Senha@123",
	})

	uc := NovoSolicitarRedefinicaoSenha(repo, tokenRepo, emailSvc)
	err := uc.Executar(context.Background(), dto.SolicitarRedefinicaoSenhaEntrada{
		Email: "joao@example.com",
	})
	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if emailSvc.ultimoEmail != "joao@example.com" {
		t.Fatalf("email nao enviado para o destinatario correto")
	}
	if emailSvc.ultimoToken == "" {
		t.Fatalf("token nao foi gerado")
	}
}

func TestSolicitarRedefinicaoSenhaParaEmailInexistenteNaoRetornaErro(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	emailSvc := &emailServiceTeste{}

	uc := NovoSolicitarRedefinicaoSenha(repo, tokenRepo, emailSvc)
	err := uc.Executar(context.Background(), dto.SolicitarRedefinicaoSenhaEntrada{
		Email: "inexistente@example.com",
	})
	if err != nil {
		t.Fatalf("esperava nil (protecao contra enumeracao), got %v", err)
	}
	if emailSvc.ultimoToken != "" {
		t.Fatalf("token nao deveria ter sido gerado para email inexistente")
	}
}

func TestSolicitarRedefinicaoSenhaComEmailInvalidoNaoRetornaErro(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	emailSvc := &emailServiceTeste{}

	uc := NovoSolicitarRedefinicaoSenha(repo, tokenRepo, emailSvc)
	err := uc.Executar(context.Background(), dto.SolicitarRedefinicaoSenhaEntrada{
		Email: "email-invalido",
	})
	if err != nil {
		t.Fatalf("esperava nil (protecao contra enumeracao), got %v", err)
	}
}

func TestRedefinirSenhaComSucesso(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	hash := hashTeste{}
	hashRedef := hashRedefinicaoTeste{}

	cadastrar := NovoCadastrarUsuario(repo, hash)
	saida, _ := cadastrar.Executar(context.Background(), dto.CadastrarUsuarioEntrada{
		Nome: "João Silva", CPF: "529.982.247-25",
		Email: "joao@example.com", ConfirmacaoEmail: "joao@example.com",
		Senha: "Senha@123", ConfirmacaoSenha: "Senha@123",
	})

	token := objetos_de_valor.NovoTokeRedefinicao(saida.ID, 1*time.Hour)
	tokenRepo.Inserir(context.Background(), token)

	uc := NovoRedefinirSenha(repo, tokenRepo, hashRedef)
	err := uc.Executar(context.Background(), dto.RedefinirSenhaEntrada{
		Token:            token.Token,
		Senha:            "NovaSenha@456",
		ConfirmacaoSenha: "NovaSenha@456",
	})
	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}

	usuario, _ := repo.BuscarPorCPF(context.Background(), "52998224725")
	if usuario == nil {
		t.Fatal("usuario nao encontrado")
	}
	if usuario.SenhaHash.String() == "" {
		t.Fatal("senha nao foi atualizada")
	}
}

func TestRedefinirSenhaComTokenExpirado(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	hash := hashTeste{}
	hashRedef := hashRedefinicaoTeste{}

	cadastrar := NovoCadastrarUsuario(repo, hash)
	saida, _ := cadastrar.Executar(context.Background(), dto.CadastrarUsuarioEntrada{
		Nome: "João Silva", CPF: "529.982.247-25",
		Email: "joao@example.com", ConfirmacaoEmail: "joao@example.com",
		Senha: "Senha@123", ConfirmacaoSenha: "Senha@123",
	})

	token := objetos_de_valor.NovoTokeRedefinicao(saida.ID, 1*time.Nanosecond)
	tokenRepo.Inserir(context.Background(), token)
	time.Sleep(2 * time.Nanosecond)

	uc := NovoRedefinirSenha(repo, tokenRepo, hashRedef)
	err := uc.Executar(context.Background(), dto.RedefinirSenhaEntrada{
		Token:            token.Token,
		Senha:            "NovaSenha@456",
		ConfirmacaoSenha: "NovaSenha@456",
	})
	if err == nil {
		t.Fatal("esperava erro de token expirado")
	}
}

func TestRedefinirSenhaComTokenInexistente(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	hashRedef := hashRedefinicaoTeste{}

	uc := NovoRedefinirSenha(repo, tokenRepo, hashRedef)
	err := uc.Executar(context.Background(), dto.RedefinirSenhaEntrada{
		Token:            "token-invalido",
		Senha:            "NovaSenha@456",
		ConfirmacaoSenha: "NovaSenha@456",
	})
	if err == nil {
		t.Fatal("esperava erro de token invalido")
	}
}

func TestRedefinirSenhaComSenhaFraca(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	hash := hashTeste{}
	hashRedef := hashRedefinicaoTeste{}

	cadastrar := NovoCadastrarUsuario(repo, hash)
	saida, _ := cadastrar.Executar(context.Background(), dto.CadastrarUsuarioEntrada{
		Nome: "João Silva", CPF: "529.982.247-25",
		Email: "joao@example.com", ConfirmacaoEmail: "joao@example.com",
		Senha: "Senha@123", ConfirmacaoSenha: "Senha@123",
	})

	token := objetos_de_valor.NovoTokeRedefinicao(saida.ID, 1*time.Hour)
	tokenRepo.Inserir(context.Background(), token)

	uc := NovoRedefinirSenha(repo, tokenRepo, hashRedef)
	err := uc.Executar(context.Background(), dto.RedefinirSenhaEntrada{
		Token:            token.Token,
		Senha:            "fraca",
		ConfirmacaoSenha: "fraca",
	})
	if err == nil {
		t.Fatal("esperava erro de senha fraca")
	}
}

func TestRedefinirSenhaComConfirmacaoDiferente(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	hash := hashTeste{}
	hashRedef := hashRedefinicaoTeste{}

	cadastrar := NovoCadastrarUsuario(repo, hash)
	saida, _ := cadastrar.Executar(context.Background(), dto.CadastrarUsuarioEntrada{
		Nome: "João Silva", CPF: "529.982.247-25",
		Email: "joao@example.com", ConfirmacaoEmail: "joao@example.com",
		Senha: "Senha@123", ConfirmacaoSenha: "Senha@123",
	})

	token := objetos_de_valor.NovoTokeRedefinicao(saida.ID, 1*time.Hour)
	tokenRepo.Inserir(context.Background(), token)

	uc := NovoRedefinirSenha(repo, tokenRepo, hashRedef)
	err := uc.Executar(context.Background(), dto.RedefinirSenhaEntrada{
		Token:            token.Token,
		Senha:            "NovaSenha@456",
		ConfirmacaoSenha: "OutraSenha@789",
	})
	if err == nil {
		t.Fatal("esperava erro de confirmacao diferente")
	}
}
