package casos_de_uso

import (
	"context"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func setupAvaliadoresTest(t *testing.T) (repositorios.RepositorioDeAvaliador, repositorios.RepositorioDeAtribuicao, context.Context) {
	t.Helper()
	repo := persistencia.NovoRepositorioDeAvaliadorMemoria()
	repoAtrib := persistencia.NovoRepositorioDeAtribuicaoMemoria()
	ctx := context.Background()
	cadastrar := NovoCadastrarAvaliador(repo)

	entradas := []dto.CadastrarAvaliadorEntrada{
		{UsuarioID: 1, Nome: "Ana Pesquisadora", CPF: "111.111.111-11", Email: "ana@teste.com"},
		{UsuarioID: 2, Nome: "Bruno Cientista", CPF: "222.222.222-22", Email: "bruno@teste.com"},
		{UsuarioID: 3, Nome: "Carla Doutora", CPF: "333.333.333-33", Email: "carla@teste.com"},
	}
	for _, e := range entradas {
		cadastrar.Executar(ctx, e)
	}

	return repo, repoAtrib, ctx
}

func TestListarAvaliadores(t *testing.T) {
	repo, _, ctx := setupAvaliadoresTest(t)
	uc := NovoListarAvaliadores(repo)

	t.Run("deve listar todos os avaliadores", func(t *testing.T) {
		resultados, err := uc.Executar(ctx, repositorios.FiltrosListarAvaliadores{})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(resultados) != 3 {
			t.Errorf("esperava 3 avaliadores, obteve %d", len(resultados))
		}
	})

	t.Run("deve filtrar por nome", func(t *testing.T) {
		resultados, err := uc.Executar(ctx, repositorios.FiltrosListarAvaliadores{
			Nome: "Ana",
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(resultados) != 1 {
			t.Errorf("esperava 1 avaliador, obteve %d", len(resultados))
		}
		if resultados[0].Nome != "Ana Pesquisadora" {
			t.Errorf("esperava Ana Pesquisadora, obteve %s", resultados[0].Nome)
		}
	})

	t.Run("deve retornar lista vazia para filtro sem resultados", func(t *testing.T) {
		resultados, err := uc.Executar(ctx, repositorios.FiltrosListarAvaliadores{
			Nome: "Inexistente",
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(resultados) != 0 {
			t.Errorf("esperava 0 avaliadores, obteve %d", len(resultados))
		}
	})
}

func TestVisualizarAvaliador(t *testing.T) {
	repo, repoAtrib, ctx := setupAvaliadoresTest(t)
	uc := NovoVisualizarAvaliador(repo, repoAtrib)

	t.Run("deve visualizar avaliador existente", func(t *testing.T) {
		saida, err := uc.Executar(ctx, 1)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.ID != 1 {
			t.Errorf("ID esperado 1, obteve %d", saida.ID)
		}
		if saida.Nome != "Ana Pesquisadora" {
			t.Errorf("nome esperado Ana Pesquisadora, obteve %s", saida.Nome)
		}
	})

	t.Run("deve retornar erro para avaliador inexistente", func(t *testing.T) {
		_, err := uc.Executar(ctx, 999)
		if err == nil {
			t.Error("esperava erro para avaliador inexistente")
		}
	})
}
