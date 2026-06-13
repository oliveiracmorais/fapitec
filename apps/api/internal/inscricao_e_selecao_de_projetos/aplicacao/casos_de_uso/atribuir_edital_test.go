package casos_de_uso

import (
	"context"
	"testing"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func TestAtribuirEdital(t *testing.T) {
	repoAvaliador := persistencia.NovoRepositorioDeAvaliadorMemoria()
	repoAtribuicao := persistencia.NovoRepositorioDeAtribuicaoMemoria()
	ctx := context.Background()

	cadastrar := NovoCadastrarAvaliador(repoAvaliador)
	atribuir := NovoAtribuirEdital(repoAvaliador, repoAtribuicao)

	cadastrar.Executar(ctx, dto.CadastrarAvaliadorEntrada{
		UsuarioID: 1,
		Nome:      "Ana Pesquisadora",
		CPF:       "111.111.111-11",
		Email:     "ana@teste.com",
	})

	t.Run("deve atribuir avaliador a edital", func(t *testing.T) {
		saida, err := atribuir.Executar(ctx, 1, dto.AtribuirEditalEntrada{
			EditalID:  10,
			DataInicio: time.Now().Format("2006-01-02"),
			DataFim:    time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.EditalID != 10 {
			t.Errorf("edital_id esperado 10, obteve %d", saida.EditalID)
		}
		if saida.AvaliadorID != 1 {
			t.Errorf("avaliador_id esperado 1, obteve %d", saida.AvaliadorID)
		}
		if saida.StatusConvite != "pendente" {
			t.Errorf("status convite deve ser pendente, obteve %s", saida.StatusConvite)
		}
		if saida.HashAnonimizacao == "" {
			t.Error("hash de anonimizacao nao deve ser vazio")
		}
	})

	t.Run("deve rejeitar atribuir a avaliador inexistente", func(t *testing.T) {
		_, err := atribuir.Executar(ctx, 999, dto.AtribuirEditalEntrada{
			EditalID:   10,
			DataInicio: time.Now().Format("2006-01-02"),
			DataFim:    time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
		})
		if err == nil {
			t.Error("esperava erro para avaliador inexistente")
		}
	})

	t.Run("deve rejeitar data inicio invalida", func(t *testing.T) {
		_, err := atribuir.Executar(ctx, 1, dto.AtribuirEditalEntrada{
			EditalID:   10,
			DataInicio: "data-invalida",
			DataFim:    time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
		})
		if err == nil {
			t.Error("esperava erro para data invalida")
		}
	})

	t.Run("deve rejeitar data fim anterior a data inicio", func(t *testing.T) {
		_, err := atribuir.Executar(ctx, 1, dto.AtribuirEditalEntrada{
			EditalID:   10,
			DataInicio: time.Now().Format("2006-01-02"),
			DataFim:    time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
		})
		if err == nil {
			t.Error("esperava erro para data fim anterior a data inicio")
		}
	})
}

func TestListarAtribuicoes(t *testing.T) {
	repoAvaliador := persistencia.NovoRepositorioDeAvaliadorMemoria()
	repoAtribuicao := persistencia.NovoRepositorioDeAtribuicaoMemoria()
	ctx := context.Background()

	cadastrar := NovoCadastrarAvaliador(repoAvaliador)
	atribuir := NovoAtribuirEdital(repoAvaliador, repoAtribuicao)
	listar := NovoListarAtribuicoes(repoAtribuicao)

	cadastrar.Executar(ctx, dto.CadastrarAvaliadorEntrada{
		UsuarioID: 1, Nome: "Ana", CPF: "111.111.111-11", Email: "ana@teste.com",
	})

	atribuir.Executar(ctx, 1, dto.AtribuirEditalEntrada{
		EditalID: 10, DataInicio: "2026-01-01", DataFim: "2026-06-30",
	})
	atribuir.Executar(ctx, 1, dto.AtribuirEditalEntrada{
		EditalID: 11, DataInicio: "2026-07-01", DataFim: "2026-12-31",
	})

	t.Run("deve listar atribuicoes por avaliador", func(t *testing.T) {
		resultados, err := listar.ExecutarPorAvaliador(ctx, 1)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(resultados) != 2 {
			t.Errorf("esperava 2 atribuicoes, obteve %d", len(resultados))
		}
	})

	t.Run("deve listar atribuicoes por edital", func(t *testing.T) {
		resultados, err := listar.ExecutarPorEdital(ctx, 10)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(resultados) != 1 {
			t.Errorf("esperava 1 atribuicao, obteve %d", len(resultados))
		}
	})

	t.Run("deve retornar lista vazia para avaliador sem atribuicoes", func(t *testing.T) {
		resultados, err := listar.ExecutarPorAvaliador(ctx, 999)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(resultados) != 0 {
			t.Errorf("esperava 0 atribuicoes, obteve %d", len(resultados))
		}
	})
}

func TestGerenciarConvite(t *testing.T) {
	repoAvaliador := persistencia.NovoRepositorioDeAvaliadorMemoria()
	repoAtribuicao := persistencia.NovoRepositorioDeAtribuicaoMemoria()
	ctx := context.Background()

	cadastrar := NovoCadastrarAvaliador(repoAvaliador)
	atribuir := NovoAtribuirEdital(repoAvaliador, repoAtribuicao)
	gerenciar := NovoGerenciarConvite(repoAtribuicao)

	cadastrar.Executar(ctx, dto.CadastrarAvaliadorEntrada{
		UsuarioID: 1, Nome: "Ana", CPF: "111.111.111-11", Email: "ana@teste.com",
	})
	atribuir.Executar(ctx, 1, dto.AtribuirEditalEntrada{
		EditalID: 10, DataInicio: "2026-01-01", DataFim: "2026-06-30",
	})

	t.Run("deve aceitar convite", func(t *testing.T) {
		saida, err := gerenciar.Executar(ctx, 1, dto.GerenciarConviteEntrada{Acao: "aceitar"})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.StatusConvite != "aceito" {
			t.Errorf("status deve ser aceito, obteve %s", saida.StatusConvite)
		}
	})

	t.Run("deve rejeitar aceitar convite ja processado", func(t *testing.T) {
		_, err := gerenciar.Executar(ctx, 1, dto.GerenciarConviteEntrada{Acao: "aceitar"})
		if err == nil {
			t.Error("esperava erro para convite ja aceito")
		}
	})

	t.Run("deve recusar convite pendente", func(t *testing.T) {
		repoAvaliador2 := persistencia.NovoRepositorioDeAvaliadorMemoria()
		repoAtribuicao2 := persistencia.NovoRepositorioDeAtribuicaoMemoria()
		ctx := context.Background()

		cadastrar2 := NovoCadastrarAvaliador(repoAvaliador2)
		atribuir2 := NovoAtribuirEdital(repoAvaliador2, repoAtribuicao2)
		gerenciar2 := NovoGerenciarConvite(repoAtribuicao2)

		cadastrar2.Executar(ctx, dto.CadastrarAvaliadorEntrada{
			UsuarioID: 2, Nome: "Bruno", CPF: "222.222.222-22", Email: "bruno@teste.com",
		})
		atribuir2.Executar(ctx, 1, dto.AtribuirEditalEntrada{
			EditalID: 20, DataInicio: "2026-01-01", DataFim: "2026-06-30",
		})

		saida, err := gerenciar2.Executar(ctx, 1, dto.GerenciarConviteEntrada{Acao: "recusar"})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.StatusConvite != "recusado" {
			t.Errorf("status deve ser recusado, obteve %s", saida.StatusConvite)
		}
	})

	t.Run("deve rejeitar acao invalida", func(t *testing.T) {
		_, err := gerenciar.Executar(ctx, 1, dto.GerenciarConviteEntrada{Acao: "cancelar"})
		if err == nil {
			t.Error("esperava erro para acao invalida")
		}
	})

	t.Run("deve rejeitar atribuicao inexistente", func(t *testing.T) {
		_, err := gerenciar.Executar(ctx, 999, dto.GerenciarConviteEntrada{Acao: "aceitar"})
		if err == nil {
			t.Error("esperava erro para atribuicao inexistente")
		}
	})
}
