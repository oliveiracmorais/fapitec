package entidades

import (
	"testing"
)

func TestNovaProposta(t *testing.T) {
	t.Run("deve criar proposta com dados validos", func(t *testing.T) {
		params := NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
			DadosAcademicos: DadosAcademicos{
				MaiorTitulacao: "Mestre",
				Curso:          "Ciência da Computação",
			},
		}

		proposta, err := NovaProposta(params)
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if proposta.Status.String() != "rascunho" {
			t.Errorf("status inicial deve ser rascunho, obteve %s", proposta.Status.String())
		}
		if proposta.ValorTotalSolicitado != 0 {
			t.Errorf("valor total inicial deve ser 0, obteve %d", proposta.ValorTotalSolicitado)
		}
	})

	t.Run("deve rejeitar edital vazio", func(t *testing.T) {
		params := NovaPropostaParams{
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		}
		_, err := NovaProposta(params)
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar proponente sem nome", func(t *testing.T) {
		params := NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				CPF: "123.456.789-00",
			},
		}
		_, err := NovaProposta(params)
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve calcular valor total dos itens", func(t *testing.T) {
		params := NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
			ItensOrcamentarios: []ItemOrcamentario{
				{Descricao: "Item 1", ValorTotal: 1000},
				{Descricao: "Item 2", ValorTotal: 2000},
			},
		}

		proposta, err := NovaProposta(params)
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if proposta.ValorTotalSolicitado != 3000 {
			t.Errorf("valor total deve ser 3000, obteve %d", proposta.ValorTotalSolicitado)
		}
	})
}

func TestSubmeter(t *testing.T) {
	t.Run("deve submeter proposta em rascunho", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})

		err := proposta.Submeter()
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if proposta.Status.String() != "submetida" {
			t.Errorf("status deve ser submetida, obteve %s", proposta.Status.String())
		}
	})

	t.Run("deve rejeitar submeter proposta ja submetida", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})
		proposta.Submeter()

		err := proposta.Submeter()
		if err == nil {
			t.Error("esperava erro ao submeter proposta ja submetida")
		}
	})
}

func TestCalcularValorTotal(t *testing.T) {
	t.Run("deve recalcular valor total", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
			ItensOrcamentarios: []ItemOrcamentario{
				{Descricao: "Item 1", ValorTotal: 5000},
			},
		})

		proposta.ItensOrcamentarios = append(proposta.ItensOrcamentarios, ItemOrcamentario{Descricao: "Item 2", ValorTotal: 3000})
		proposta.CalcularValorTotal()

		if proposta.ValorTotalSolicitado != 8000 {
			t.Errorf("valor total deve ser 8000, obteve %d", proposta.ValorTotalSolicitado)
		}
	})
}
