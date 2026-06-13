package entidades

import (
	"testing"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
)

func TestNovaAtribuicao(t *testing.T) {
	t.Run("deve criar atribuicao com dados validos", func(t *testing.T) {
		dataInicio := time.Now()
		dataFim := dataInicio.AddDate(0, 6, 0)

		atr, err := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  dataInicio,
			DataFim:     dataFim,
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if atr.AvaliadorID != 1 {
			t.Errorf("avaliador_id esperado 1, obteve %d", atr.AvaliadorID)
		}
		if atr.EditalID != 10 {
			t.Errorf("edital_id esperado 10, obteve %d", atr.EditalID)
		}
		if atr.StatusConvite != objetos_de_valor.StatusConvitePendente {
			t.Errorf("status inicial deve ser pendente, obteve %s", atr.StatusConvite)
		}
		if atr.HashAnonimizacao.IsZero() {
			t.Error("hash de anonimizacao nao deve ser vazio")
		}
	})

	t.Run("deve rejeitar avaliador_id zero", func(t *testing.T) {
		_, err := NovaAtribuicao(NovaAtribuicaoParams{
			EditalID:   10,
			DataInicio: time.Now(),
			DataFim:    time.Now().AddDate(0, 6, 0),
		})
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar edital_id zero", func(t *testing.T) {
		_, err := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar data fim anterior a data inicio", func(t *testing.T) {
		_, err := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, -1, 0),
		})
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar data fim igual a data inicio", func(t *testing.T) {
		agora := time.Now()
		_, err := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  agora,
			DataFim:     agora,
		})
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})
}

func TestAtribuicaoAceitarConvite(t *testing.T) {
	t.Run("deve aceitar convite pendente", func(t *testing.T) {
		atr, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})

		err := atr.AceitarConvite()
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if atr.StatusConvite != objetos_de_valor.StatusConviteAceito {
			t.Errorf("status deve ser aceito, obteve %s", atr.StatusConvite)
		}
	})

	t.Run("deve rejeitar aceitar convite ja aceito", func(t *testing.T) {
		atr, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})
		atr.AceitarConvite()

		err := atr.AceitarConvite()
		if err == nil {
			t.Error("esperava erro ao aceitar convite ja aceito")
		}
	})

	t.Run("deve rejeitar aceitar convite recusado", func(t *testing.T) {
		atr, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})
		atr.RecusarConvite()

		err := atr.AceitarConvite()
		if err == nil {
			t.Error("esperava erro ao aceitar convite recusado")
		}
	})
}

func TestAtribuicaoRecusarConvite(t *testing.T) {
	t.Run("deve recusar convite pendente", func(t *testing.T) {
		atr, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})

		err := atr.RecusarConvite()
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if atr.StatusConvite != objetos_de_valor.StatusConviteRecusado {
			t.Errorf("status deve ser recusado, obteve %s", atr.StatusConvite)
		}
	})

	t.Run("deve rejeitar recusar convite ja aceito", func(t *testing.T) {
		atr, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})
		atr.AceitarConvite()

		err := atr.RecusarConvite()
		if err == nil {
			t.Error("esperava erro ao recusar convite ja aceito")
		}
	})
}

func TestAtribuicaoHashAnonimizacao(t *testing.T) {
	t.Run("hash deve ser deterministico para mesmo avaliador e edital", func(t *testing.T) {
		atr1, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})
		atr2, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})

		if atr1.HashAnonimizacao.String() != atr2.HashAnonimizacao.String() {
			t.Error("hash deve ser deterministico para mesmo avaliador_id + edital_id")
		}
	})

	t.Run("hash deve ser diferente para editais diferentes", func(t *testing.T) {
		atr1, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    10,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})
		atr2, _ := NovaAtribuicao(NovaAtribuicaoParams{
			AvaliadorID: 1,
			EditalID:    20,
			DataInicio:  time.Now(),
			DataFim:     time.Now().AddDate(0, 6, 0),
		})

		if atr1.HashAnonimizacao.String() == atr2.HashAnonimizacao.String() {
			t.Error("hash deve ser diferente para editais diferentes")
		}
	})
}
