package entidades

import (
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
)

func TestNovoAvaliador(t *testing.T) {
	t.Run("deve criar avaliador com dados validos", func(t *testing.T) {
		a, err := NovoAvaliador(NovoAvaliadorParams{
			UsuarioID:         1,
			Nome:              "Maria Santos",
			CPF:               "123.456.789-00",
			Email:             "maria@teste.com",
			TitulacaoMaxima:   "Doutor",
			AreaConhecimento:  "Ciência da Computação",
			Instituicao:       "Universidade Federal",
			CurriculoResumido: "Experiência em pesquisa",
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if a.Nome != "Maria Santos" {
			t.Errorf("nome esperado Maria Santos, obteve %s", a.Nome)
		}
		if a.UsuarioID != 1 {
			t.Errorf("usuario_id esperado 1, obteve %d", a.UsuarioID)
		}
		if a.Estado != objetos_de_valor.EstadoAvaliadorAtivo {
			t.Errorf("estado inicial deve ser ativo, obteve %s", a.Estado)
		}
		if a.DataCadastro.IsZero() {
			t.Error("data de cadastro nao deve ser zero")
		}
	})

	t.Run("deve rejeitar nome vazio", func(t *testing.T) {
		_, err := NovoAvaliador(NovoAvaliadorParams{
			UsuarioID: 1,
			CPF:       "123.456.789-00",
		})
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar CPF vazio", func(t *testing.T) {
		_, err := NovoAvaliador(NovoAvaliadorParams{
			UsuarioID: 1,
			Nome:      "Maria Santos",
		})
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar usuario_id zero", func(t *testing.T) {
		_, err := NovoAvaliador(NovoAvaliadorParams{
			Nome: "Maria Santos",
			CPF:  "123.456.789-00",
		})
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})
}

func TestAvaliadorAtualizarDados(t *testing.T) {
	t.Run("deve atualizar todos os dados", func(t *testing.T) {
		a, _ := NovoAvaliador(NovoAvaliadorParams{
			UsuarioID: 1,
			Nome:      "Maria Santos",
			CPF:       "123.456.789-00",
		})
		dataOriginal := a.DataAtualizacao

		a.AtualizarDados(NovoAvaliadorParams{
			Nome:              "Maria Santos Atualizada",
			CPF:               "987.654.321-00",
			Email:             "maria.novo@teste.com",
			TitulacaoMaxima:   "Pós-Doutor",
			AreaConhecimento:  "Matemática",
			Instituicao:       "Universidade Estadual",
			CurriculoResumido: "Novo currículo",
		})

		if a.Nome != "Maria Santos Atualizada" {
			t.Errorf("nome nao atualizado: %s", a.Nome)
		}
		if a.CPF != "987.654.321-00" {
			t.Errorf("CPF nao atualizado: %s", a.CPF)
		}
		if a.DataAtualizacao.Equal(dataOriginal) {
			t.Error("data de atualizacao deveria ter mudado")
		}
	})
}

func TestAvaliadorAtivarInativar(t *testing.T) {
	t.Run("deve inativar avaliador", func(t *testing.T) {
		a, _ := NovoAvaliador(NovoAvaliadorParams{
			UsuarioID: 1,
			Nome:      "Maria Santos",
			CPF:       "123.456.789-00",
		})
		a.Inativar()
		if a.Estado != objetos_de_valor.EstadoAvaliadorInativo {
			t.Errorf("estado deve ser inativo, obteve %s", a.Estado)
		}
	})

	t.Run("deve ativar avaliador inativo", func(t *testing.T) {
		a, _ := NovoAvaliador(NovoAvaliadorParams{
			UsuarioID: 1,
			Nome:      "Maria Santos",
			CPF:       "123.456.789-00",
		})
		a.Inativar()
		a.Ativar()
		if a.Estado != objetos_de_valor.EstadoAvaliadorAtivo {
			t.Errorf("estado deve ser ativo, obteve %s", a.Estado)
		}
	})
}
