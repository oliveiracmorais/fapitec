package casos_de_uso

import (
	"context"
	"sync"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/hash"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia"
)

type registradorAuditoriaTeste struct {
	mu     sync.Mutex
	eventos []struct {
		Acao      string
		AtorID    int64
		AtorNome  string
		AtorCPF   string
		Resultado string
	}
}

func (r *registradorAuditoriaTeste) Registrar(_ context.Context, input RegistrarEventoInput) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.eventos = append(r.eventos, struct {
		Acao      string
		AtorID    int64
		AtorNome  string
		AtorCPF   string
		Resultado string
	}{Acao: input.Acao, AtorID: input.AtorID, AtorNome: input.AtorNome, AtorCPF: input.AtorCPF, Resultado: input.Resultado})
	return nil
}

func (r *registradorAuditoriaTeste) ultimoEvento() *struct {
	Acao      string
	AtorID    int64
	AtorNome  string
	AtorCPF   string
	Resultado string
} {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.eventos) == 0 {
		return nil
	}
	return &r.eventos[len(r.eventos)-1]
}

func (r *registradorAuditoriaTeste) totalEventos() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.eventos)
}

func TestCadastrarUsuarioGeraEventoAuditoria(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	audit := &registradorAuditoriaTeste{}
	uc := NovoCadastrarUsuarioComAuditoria(repo, hashTeste{}, audit)

	entrada := dto.CadastrarUsuarioEntrada{
		Nome:             "João Silva",
		CPF:              "529.982.247-25",
		Email:            "joao@example.com",
		ConfirmacaoEmail: "joao@example.com",
		Senha:            "Senha@123",
		ConfirmacaoSenha: "Senha@123",
	}

	_, err := uc.Executar(context.Background(), entrada)
	if err != nil {
		t.Fatalf("Cadastro retornou erro: %v", err)
	}

	if audit.totalEventos() != 1 {
		t.Fatalf("Esperado 1 evento de auditoria, got %d", audit.totalEventos())
	}
	ev := audit.ultimoEvento()
	if ev.Acao != "cadastro_de_usuario" {
		t.Fatalf("Acao esperada 'cadastro_de_usuario', got '%s'", ev.Acao)
	}
	if ev.Resultado != "sucesso" {
		t.Fatalf("Resultado esperado 'sucesso', got '%s'", ev.Resultado)
	}
}

func TestAutenticarSucessoGeraEventoAuditoria(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	audit := &registradorAuditoriaTeste{}

	hashService := hash.NovoServicoDeHashBcrypt()
	cadastro := NovoCadastrarUsuarioComAuditoria(repo, hashService, audit)
	cadastro.Executar(context.Background(), dto.CadastrarUsuarioEntrada{
		Nome: "Maria Souza", CPF: "529.982.247-25", Email: "maria@example.com",
		ConfirmacaoEmail: "maria@example.com", Senha: "Senha@123", ConfirmacaoSenha: "Senha@123",
	})
	audit.eventos = nil

	uc := NovoAutenticarUsuarioComAuditoria(repo, hashService, audit)
	_, err := uc.Executar(context.Background(), dto.AutenticarUsuarioEntrada{
		CPF: "529.982.247-25", Senha: "Senha@123",
	})
	if err != nil {
		t.Fatalf("Autenticacao retornou erro: %v", err)
	}

	ev := audit.ultimoEvento()
	if ev == nil {
		t.Fatal("Nenhum evento de auditoria registrado")
	}
	if ev.Acao != "login" {
		t.Fatalf("Acao esperada 'login', got '%s'", ev.Acao)
	}
	if ev.Resultado != "sucesso" {
		t.Fatalf("Resultado esperado 'sucesso', got '%s'", ev.Resultado)
	}
}

func TestAutenticarFalhaGeraEventoAuditoria(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	audit := &registradorAuditoriaTeste{}

	hashService := hash.NovoServicoDeHashBcrypt()
	cadastro := NovoCadastrarUsuarioComAuditoria(repo, hashService, audit)
	cadastro.Executar(context.Background(), dto.CadastrarUsuarioEntrada{
		Nome: "Maria Souza", CPF: "529.982.247-25", Email: "maria@example.com",
		ConfirmacaoEmail: "maria@example.com", Senha: "Senha@123", ConfirmacaoSenha: "Senha@123",
	})
	audit.eventos = nil

	uc := NovoAutenticarUsuarioComAuditoria(repo, hashService, audit)
	uc.Executar(context.Background(), dto.AutenticarUsuarioEntrada{
		CPF: "529.982.247-25", Senha: "SenhaErrada@123",
	})

	ev := audit.ultimoEvento()
	if ev == nil {
		t.Fatal("Nenhum evento de auditoria registrado")
	}
	if ev.Acao != "falha_de_login" {
		t.Fatalf("Acao esperada 'falha_de_login', got '%s'", ev.Acao)
	}
	if ev.Resultado != "falha" {
		t.Fatalf("Resultado esperado 'falha', got '%s'", ev.Resultado)
	}
}

func TestAutenticarBloqueioGeraEventoAuditoria(t *testing.T) {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	audit := &registradorAuditoriaTeste{}

	hashService := hash.NovoServicoDeHashBcrypt()
	cadastro := NovoCadastrarUsuarioComAuditoria(repo, hashService, audit)
	cadastro.Executar(context.Background(), dto.CadastrarUsuarioEntrada{
		Nome: "Maria Souza", CPF: "529.982.247-25", Email: "maria@example.com",
		ConfirmacaoEmail: "maria@example.com", Senha: "Senha@123", ConfirmacaoSenha: "Senha@123",
	})
	audit.eventos = nil

	uc := NovoAutenticarUsuarioComAuditoria(repo, hashService, audit)
	entrada := dto.AutenticarUsuarioEntrada{CPF: "529.982.247-25", Senha: "SenhaErrada@123"}

	for i := 0; i < 4; i++ {
		uc.Executar(context.Background(), entrada)
	}

	uc.Executar(context.Background(), entrada)

	eventos := audit.eventos
	quantidade := len(eventos)
	if quantidade < 2 {
		t.Fatal("Esperava ao menos 2 eventos na 5a tentativa")
	}
	penultimo := eventos[quantidade-2]
	if penultimo.Acao != "falha_de_login" {
		t.Fatalf("Penultimo evento deveria ser 'falha_de_login', got '%s'", penultimo.Acao)
	}
	if penultimo.Resultado != "falha" {
		t.Fatalf("Resultado do penultimo evento esperado 'falha', got '%s'", penultimo.Resultado)
	}
	ultimo := eventos[quantidade-1]
	if ultimo.Acao != "bloqueio_de_conta" {
		t.Fatalf("Ultimo evento deveria ser 'bloqueio_de_conta', got '%s'", ultimo.Acao)
	}
	if ultimo.Resultado != "negado" {
		t.Fatalf("Resultado esperado 'negado', got '%s'", ultimo.Resultado)
	}

	_, err := uc.Executar(context.Background(), entrada)
	if err == nil {
		t.Fatal("Conta deveria estar bloqueada")
	}
}

func TestAuditoriaNaoRegistraSenha(t *testing.T) {
	audit := &registradorAuditoriaTeste{}
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	hashService := hash.NovoServicoDeHashBcrypt()

	uc := NovoAutenticarUsuarioComAuditoria(repo, hashService, audit)
	uc.Executar(context.Background(), dto.AutenticarUsuarioEntrada{
		CPF: "529.982.247-25", Senha: "SenhaSecreta@123",
	})

	for _, ev := range audit.eventos {
		if ev.AtorCPF == "SenhaSecreta@123" {
			t.Fatal("Senha nao deve aparecer nos eventos de auditoria")
		}
	}
}
