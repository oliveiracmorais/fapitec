package entidades

import (
	"testing"
	"time"
)

func TestNovoEdital(t *testing.T) {
	t.Run("deve criar edital com dados validos", func(t *testing.T) {
		edital, err := NovoEdital(NovoEditalParams{
			Nome:        "Edital de Teste",
			Descricao:   "Descricao do edital de teste",
			DataInicio:  time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
			DataFim:     time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
			TipoChamada: "APQ",
		})
		if err != nil {
			t.Fatalf("esperava nil, got %v", err)
		}
		if edital.Nome != "Edital de Teste" {
			t.Errorf("esperava 'Edital de Teste', got %s", edital.Nome)
		}
		if edital.Status.String() != "ativo" {
			t.Errorf("esperava 'ativo', got %s", edital.Status.String())
		}
	})

	t.Run("deve rejeitar nome vazio", func(t *testing.T) {
		_, err := NovoEdital(NovoEditalParams{
			Nome:       "",
			Descricao:  "Descricao",
			DataInicio: time.Now(),
			DataFim:    time.Now().AddDate(0, 6, 0),
		})
		if err == nil {
			t.Fatal("esperava erro, got nil")
		}
	})

	t.Run("deve rejeitar descricao vazia", func(t *testing.T) {
		_, err := NovoEdital(NovoEditalParams{
			Nome:       "Edital Teste",
			Descricao:  "",
			DataInicio: time.Now(),
			DataFim:    time.Now().AddDate(0, 6, 0),
		})
		if err == nil {
			t.Fatal("esperava erro, got nil")
		}
	})

	t.Run("deve rejeitar data de inicio maior que data fim", func(t *testing.T) {
		_, err := NovoEdital(NovoEditalParams{
			Nome:       "Edital Teste",
			Descricao:  "Descricao",
			DataInicio: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
			DataFim:    time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		})
		if err == nil {
			t.Fatal("esperava erro, got nil")
		}
	})
}
