package entidades

import (
	"testing"
)

func TestNovoParecer(t *testing.T) {
	t.Run("deve criar parecer valido", func(t *testing.T) {
		p, err := NovoParecer(NovoParecerParams{
			PropostaID:   1,
			Etapa:        "unica",
			AvaliadorID:  1,
			Nota:         85,
			ParecerTexto: "Projeto bem estruturado e viavel",
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if p.PropostaID != 1 {
			t.Errorf("esperava PropostaID 1, obteve %d", p.PropostaID)
		}
		if p.Nota != 85 {
			t.Errorf("esperava Nota 85, obteve %d", p.Nota)
		}
		if p.ParecerTexto != "Projeto bem estruturado e viavel" {
			t.Errorf("parecer textual incorreto")
		}
		if p.Data.IsZero() {
			t.Error("data nao deveria ser zero")
		}
	})

	t.Run("deve rejeitar nota abaixo de 0", func(t *testing.T) {
		_, err := NovoParecer(NovoParecerParams{
			PropostaID:   1,
			Etapa:        "unica",
			AvaliadorID:  1,
			Nota:         -1,
			ParecerTexto: "Inviavel",
		})
		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar nota acima de 100", func(t *testing.T) {
		_, err := NovoParecer(NovoParecerParams{
			PropostaID:   1,
			Etapa:        "unica",
			AvaliadorID:  1,
			Nota:         101,
			ParecerTexto: "Excede limite",
		})
		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar proposta ID zero", func(t *testing.T) {
		_, err := NovoParecer(NovoParecerParams{
			PropostaID:   0,
			Etapa:        "unica",
			AvaliadorID:  1,
			Nota:         50,
			ParecerTexto: "Teste",
		})
		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar etapa vazia", func(t *testing.T) {
		_, err := NovoParecer(NovoParecerParams{
			PropostaID:   1,
			Etapa:        "",
			AvaliadorID:  1,
			Nota:         50,
			ParecerTexto: "Teste",
		})
		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar parecer textual vazio", func(t *testing.T) {
		_, err := NovoParecer(NovoParecerParams{
			PropostaID:   1,
			Etapa:        "unica",
			AvaliadorID:  1,
			Nota:         50,
			ParecerTexto: "",
		})
		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
	})
}
