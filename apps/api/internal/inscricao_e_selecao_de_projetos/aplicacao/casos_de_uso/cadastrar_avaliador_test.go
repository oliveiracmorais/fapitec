package casos_de_uso

import (
	"context"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func TestCadastrarAvaliador(t *testing.T) {
	repo := persistencia.NovoRepositorioDeAvaliadorMemoria()
	uc := NovoCadastrarAvaliador(repo)
	ctx := context.Background()

	t.Run("deve cadastrar avaliador com dados validos", func(t *testing.T) {
		saida, err := uc.Executar(ctx, dto.CadastrarAvaliadorEntrada{
			UsuarioID:         1,
			Nome:              "João Pesquisador",
			CPF:               "123.456.789-00",
			Email:             "joao@teste.com",
			TitulacaoMaxima:   "Doutor",
			AreaConhecimento:  "Biologia",
			Instituicao:       "Universidade Federal",
			CurriculoResumido: "Pesquisador sênior",
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.ID == 0 {
			t.Error("ID do avaliador nao deve ser zero")
		}
		if saida.Nome != "João Pesquisador" {
			t.Errorf("nome esperado Joao Pesquisador, obteve %s", saida.Nome)
		}
		if saida.Estado != "ativo" {
			t.Errorf("estado inicial deve ser ativo, obteve %s", saida.Estado)
		}
	})

	t.Run("deve rejeitar CPF duplicado", func(t *testing.T) {
		_, err := uc.Executar(ctx, dto.CadastrarAvaliadorEntrada{
			UsuarioID: 2,
			Nome:      "Maria Pesquisadora",
			CPF:       "123.456.789-00",
			Email:     "maria@teste.com",
		})
		if err == nil {
			t.Error("esperava erro para CPF duplicado")
		}
	})

	t.Run("deve rejeitar nome vazio", func(t *testing.T) {
		_, err := uc.Executar(ctx, dto.CadastrarAvaliadorEntrada{
			UsuarioID: 3,
			CPF:       "987.654.321-00",
			Email:     "teste@teste.com",
		})
		if err == nil {
			t.Error("esperava erro para nome vazio")
		}
	})

	t.Run("deve rejeitar usuario_id zero", func(t *testing.T) {
		_, err := uc.Executar(ctx, dto.CadastrarAvaliadorEntrada{
			Nome:  "Teste",
			CPF:   "111.222.333-44",
			Email: "teste@teste.com",
		})
		if err == nil {
			t.Error("esperava erro para usuario_id zero")
		}
	})
}
