package casos_de_uso

import (
	"context"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func strPtr(s string) *string { return &s }

func TestEditarAvaliador(t *testing.T) {
	repo := persistencia.NovoRepositorioDeAvaliadorMemoria()
	cadastrar := NovoCadastrarAvaliador(repo)
	editar := NovoEditarAvaliador(repo)
	ctx := context.Background()

	cadastrar.Executar(ctx, dto.CadastrarAvaliadorEntrada{
		UsuarioID: 1,
		Nome:      "João Pesquisador",
		CPF:       "123.456.789-00",
		Email:     "joao@teste.com",
	})

	t.Run("deve editar nome do avaliador", func(t *testing.T) {
		saida, err := editar.Executar(ctx, 1, dto.AtualizarAvaliadorEntrada{
			Nome: strPtr("João Atualizado"),
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.Nome != "João Atualizado" {
			t.Errorf("nome esperado Joao Atualizado, obteve %s", saida.Nome)
		}
	})

	t.Run("deve editar multiplos campos", func(t *testing.T) {
		saida, err := editar.Executar(ctx, 1, dto.AtualizarAvaliadorEntrada{
			TitulacaoMaxima:  strPtr("Pós-Doutor"),
			AreaConhecimento: strPtr("Bioquímica"),
			Instituicao:      strPtr("Universidade Estadual"),
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.TitulacaoMaxima != "Pós-Doutor" {
			t.Errorf("titulacao esperada Pos-Doutor, obteve %s", saida.TitulacaoMaxima)
		}
		if saida.AreaConhecimento != "Bioquímica" {
			t.Errorf("area esperada Bioquimica, obteve %s", saida.AreaConhecimento)
		}
	})

	t.Run("deve inativar avaliador", func(t *testing.T) {
		saida, err := editar.Executar(ctx, 1, dto.AtualizarAvaliadorEntrada{
			Estado: strPtr("inativo"),
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.Estado != "inativo" {
			t.Errorf("estado esperado inativo, obteve %s", saida.Estado)
		}
	})

	t.Run("deve reativar avaliador", func(t *testing.T) {
		saida, err := editar.Executar(ctx, 1, dto.AtualizarAvaliadorEntrada{
			Estado: strPtr("ativo"),
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.Estado != "ativo" {
			t.Errorf("estado esperado ativo, obteve %s", saida.Estado)
		}
	})

	t.Run("deve rejeitar estado invalido", func(t *testing.T) {
		_, err := editar.Executar(ctx, 1, dto.AtualizarAvaliadorEntrada{
			Estado: strPtr("suspenso"),
		})
		if err == nil {
			t.Error("esperava erro para estado invalido")
		}
	})

	t.Run("deve rejeitar avaliador inexistente", func(t *testing.T) {
		_, err := editar.Executar(ctx, 999, dto.AtualizarAvaliadorEntrada{
			Nome: strPtr("Inexistente"),
		})
		if err == nil {
			t.Error("esperava erro para avaliador inexistente")
		}
	})
}
