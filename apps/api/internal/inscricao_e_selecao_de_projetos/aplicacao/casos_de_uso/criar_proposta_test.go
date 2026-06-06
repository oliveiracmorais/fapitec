package casos_de_uso

import (
	"context"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func TestCriarProposta(t *testing.T) {
	t.Run("deve criar proposta com dados validos", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		uc := NovoCriarProposta(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
			DadosAcademicos: dto.DadosAcademicosDTO{
				MaiorTitulacao: "Mestre",
			},
		}

		saida, err := uc.Executar(context.Background(), entrada)
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if saida.Status != "rascunho" {
			t.Errorf("status deve ser rascunho, obteve %s", saida.Status)
		}
		if saida.DadosProponente.Nome != "João Silva" {
			t.Errorf("nome deve ser Joao Silva, obteve %s", saida.DadosProponente.Nome)
		}
	})

	t.Run("deve rejeitar proposta sem nome do proponente", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		uc := NovoCriarProposta(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				CPF: "123.456.789-00",
			},
		}

		_, err := uc.Executar(context.Background(), entrada)
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve criar proposta com itens orcamentarios", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		uc := NovoCriarProposta(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
			ItensOrcamentarios: []dto.ItemOrcamentarioDTO{
				{
					Descricao:     "Item 1",
					Tipo:          "consumo",
					Quantidade:    2,
					ValorUnitario: 500,
					ValorTotal:    1000,
				},
			},
		}

		saida, err := uc.Executar(context.Background(), entrada)
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if saida.ValorTotal != 1000 {
			t.Errorf("valor total deve ser 1000, obteve %d", saida.ValorTotal)
		}
		if len(saida.ItensOrcamentarios) != 1 {
			t.Errorf("deve ter 1 item, obteve %d", len(saida.ItensOrcamentarios))
		}
	})

	t.Run("deve rejeitar item com tipo invalido", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		uc := NovoCriarProposta(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
			ItensOrcamentarios: []dto.ItemOrcamentarioDTO{
				{
					Descricao: "Item invalido",
					Tipo:      "invalido",
				},
			},
		}

		_, err := uc.Executar(context.Background(), entrada)
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})
}
