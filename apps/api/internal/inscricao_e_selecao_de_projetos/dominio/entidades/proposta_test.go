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

func TestAdicionarParecer(t *testing.T) {
	t.Run("deve adicionar parecer em proposta submetida e transicionar para em_avaliacao", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})
		proposta.ID = 1
		proposta.Submeter()

		parecer, err := NovoParecer(NovoParecerParams{
			PropostaID:   proposta.ID,
			Etapa:        "unica",
			AvaliadorID:  1,
			Nota:         85,
			ParecerTexto: "Excelente projeto",
		})
		if err != nil {
			t.Fatalf("erro ao criar parecer: %v", err)
		}

		err = proposta.AdicionarParecer(parecer)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}

		if proposta.Status.String() != "em_avaliacao" {
			t.Errorf("status deve ser em_avaliacao, obteve %s", proposta.Status.String())
		}
		if len(proposta.Pareceres) != 1 {
			t.Errorf("deve ter 1 parecer, obteve %d", len(proposta.Pareceres))
		}
	})

	t.Run("deve rejeitar adicionar parecer em proposta rascunho", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})
		proposta.ID = 1

		parecer, err := NovoParecer(NovoParecerParams{
			PropostaID:   proposta.ID,
			Etapa:        "unica",
			AvaliadorID:  1,
			Nota:         85,
			ParecerTexto: "Bom projeto",
		})
		if err != nil {
			t.Fatalf("erro ao criar parecer: %v", err)
		}

		err = proposta.AdicionarParecer(parecer)
		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
	})
}

func TestCalcularNotaFinal(t *testing.T) {
	t.Run("deve calcular media dos pareceres", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})
		proposta.ID = 1
		proposta.Submeter()

		p1, err := NovoParecer(NovoParecerParams{
			PropostaID: proposta.ID, Etapa: "unica", AvaliadorID: 1,
			Nota: 80, ParecerTexto: "Bom",
		})
		if err != nil {
			t.Fatalf("erro ao criar parecer: %v", err)
		}
		p2, err := NovoParecer(NovoParecerParams{
			PropostaID: proposta.ID, Etapa: "unica", AvaliadorID: 2,
			Nota: 90, ParecerTexto: "Excelente",
		})
		if err != nil {
			t.Fatalf("erro ao criar parecer: %v", err)
		}
		proposta.AdicionarParecer(p1)
		proposta.AdicionarParecer(p2)

		nota := proposta.CalcularNotaFinal()
		if nota != 85 {
			t.Errorf("nota final deve ser 85, obteve %d", nota)
		}
	})

	t.Run("deve retornar 0 quando nao ha pareceres", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})

		nota := proposta.CalcularNotaFinal()
		if nota != 0 {
			t.Errorf("nota final deve ser 0, obteve %d", nota)
		}
	})
}

func TestFinalizarAvaliacao(t *testing.T) {
	t.Run("deve aprovar proposta com nota acima do corte", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})
		proposta.ID = 1
		proposta.Submeter()

		parecer, err := NovoParecer(NovoParecerParams{
			PropostaID: proposta.ID, Etapa: "unica", AvaliadorID: 1,
			Nota: 80, ParecerTexto: "Bom projeto",
		})
		if err != nil {
			t.Fatalf("erro ao criar parecer: %v", err)
		}
		proposta.AdicionarParecer(parecer)

		err = proposta.FinalizarAvaliacao(70)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if proposta.Status.String() != "aprovada" {
			t.Errorf("status deve ser aprovada, obteve %s", proposta.Status.String())
		}
	})

	t.Run("deve reprovar proposta com nota abaixo do corte", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})
		proposta.ID = 1
		proposta.Submeter()

		parecer, err := NovoParecer(NovoParecerParams{
			PropostaID: proposta.ID, Etapa: "unica", AvaliadorID: 1,
			Nota: 50, ParecerTexto: "Projeto fraco",
		})
		if err != nil {
			t.Fatalf("erro ao criar parecer: %v", err)
		}
		proposta.AdicionarParecer(parecer)

		err = proposta.FinalizarAvaliacao(70)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if proposta.Status.String() != "reprovada" {
			t.Errorf("status deve ser reprovada, obteve %s", proposta.Status.String())
		}
	})

	t.Run("deve rejeitar finalizar proposta sem pareceres", func(t *testing.T) {
		proposta, _ := NovaProposta(NovaPropostaParams{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: ProponenteInfo{
				Nome: "João Silva",
				CPF:  "123.456.789-00",
			},
		})

		err := proposta.FinalizarAvaliacao(70)
		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
	})
}
