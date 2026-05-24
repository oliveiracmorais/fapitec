package casos_de_uso

import (
	"context"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/infraestrutura/persistencia"
)

func TestCriarEdital(t *testing.T) {
	t.Run("deve criar edital com dados validos", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDeEditalMemoria()
		caso := NovoEdital(repo)

		entrada := dto.CriarEditalEntrada{
			Nome:        "Edital APQ 2026",
			Descricao:   "Edital de apoio a pesquisa",
			DataInicio:  "2026-06-01",
			DataFim:     "2026-12-31",
			TipoChamada: "APQ",
		}

		saida, err := caso.Executar(context.Background(), entrada)
		if err != nil {
			t.Fatalf("esperava nil, got %v", err)
		}
		if saida.Nome != "Edital APQ 2026" {
			t.Errorf("esperava 'Edital APQ 2026', got %s", saida.Nome)
		}
		if saida.Status != "ativo" {
			t.Errorf("esperava 'ativo', got %s", saida.Status)
		}
		if saida.ID == 0 {
			t.Error("esperava ID > 0")
		}
	})

	t.Run("deve rejeitar data invalida", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDeEditalMemoria()
		caso := NovoEdital(repo)

		entrada := dto.CriarEditalEntrada{
			Nome:        "Edital Teste",
			Descricao:   "Descricao",
			DataInicio:  "data-invalida",
			DataFim:     "2026-12-31",
			TipoChamada: "APQ",
		}

		_, err := caso.Executar(context.Background(), entrada)
		if err == nil {
			t.Fatal("esperava erro, got nil")
		}
	})

	t.Run("deve rejeitar nome vazio", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDeEditalMemoria()
		caso := NovoEdital(repo)

		entrada := dto.CriarEditalEntrada{
			Nome:        "",
			Descricao:   "Descricao",
			DataInicio:  "2026-06-01",
			DataFim:     "2026-12-31",
			TipoChamada: "APQ",
		}

		_, err := caso.Executar(context.Background(), entrada)
		if err == nil {
			t.Fatal("esperava erro, got nil")
		}
	})
}
