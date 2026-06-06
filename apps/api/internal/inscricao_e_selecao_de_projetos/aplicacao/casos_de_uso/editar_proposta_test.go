package casos_de_uso

import (
	"context"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func TestEditarProposta(t *testing.T) {
	t.Run("deve editar proposta em rascunho", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		editar := NovoEditarProposta(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criada, _ := criar.Executar(context.Background(), entrada)

		nome := "Maria Souza"
		atualizacao := dto.AtualizarPropostaEntrada{
			DadosProponente: &dto.ProponenteInfoDTO{
				Nome: nome,
				CPF:  "123.456.789-00",
			},
		}

		saida, err := editar.Executar(context.Background(), criada.ID, atualizacao)
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if saida.DadosProponente.Nome != nome {
			t.Errorf("nome deve ser %s, obteve %s", nome, saida.DadosProponente.Nome)
		}
	})

	t.Run("deve rejeitar editar proposta submetida", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		submeter := NovoSubmeterProposta(repo, &editalVerificadorMock{info: editalAtivo()})
		editar := NovoEditarProposta(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criada, _ := criar.Executar(context.Background(), entrada)
		submeter.Executar(context.Background(), criada.ID)

		_, err := editar.Executar(context.Background(), criada.ID, dto.AtualizarPropostaEntrada{})
		if err == nil {
			t.Error("esperava erro ao editar proposta submetida")
		}
	})
}
