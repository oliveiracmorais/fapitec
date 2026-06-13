package casos_de_uso

import (
	"context"
	"testing"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func TestFinalizarAvaliacaoDoEdital(t *testing.T) {
	t.Run("deve finalizar avaliacao e gerar classificacao", func(t *testing.T) {
		propostaRepo := persistencia.NovoRepositorioDePropostaMemoria()
		avaliadorRepo := persistencia.NovoRepositorioDeAvaliadorMemoria()
		atribuicaoRepo := persistencia.NovoRepositorioDeAtribuicaoMemoria()

		avaliador, _ := entidades.NovoAvaliador(entidades.NovoAvaliadorParams{
			UsuarioID: 1, Nome: "Avaliador",
			CPF: "000.000.000-01", Email: "a@teste.com",
		})
		avaliadorRepo.Criar(context.Background(), avaliador)

		proposta, _ := entidades.NovaProposta(entidades.NovaPropostaParams{
			EditalID: 1, ProponenteID: 1,
			DadosProponente: entidades.ProponenteInfo{Nome: "Proponente", CPF: "111.111.111-11"},
		})
		propostaRepo.Criar(context.Background(), proposta)
		proposta.Submeter()
		propostaRepo.Atualizar(context.Background(), proposta)

		atr, _ := entidades.NovaAtribuicao(entidades.NovaAtribuicaoParams{
			AvaliadorID: avaliador.ID, EditalID: 1,
			DataInicio: time.Now().Add(-24 * time.Hour),
			DataFim:    time.Now().Add(24 * time.Hour),
		})
		atribuicaoRepo.Criar(context.Background(), atr)
		atr.AceitarConvite()
		atribuicaoRepo.Atualizar(context.Background(), atr)

		emitirUC := NovoEmitirParecer(propostaRepo, avaliadorRepo, atribuicaoRepo)
		_, err := emitirUC.Executar(context.Background(), dto.EmitirParecerEntrada{
			PropostaID: proposta.ID, Etapa: "unica",
			AvaliadorID: avaliador.ID, Nota: 85,
			ParecerTexto: "Excelente",
		})
		if err != nil {
			t.Fatalf("erro ao emitir parecer: %v", err)
		}

		finalizarUC := NovoFinalizarAvaliacaoDoEdital(propostaRepo)
		classificacao, err := finalizarUC.Executar(context.Background(), 1, dto.FinalizarAvaliacaoEntrada{NotaDeCorte: 70})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}

		if len(classificacao) != 1 {
			t.Fatalf("esperava 1 item, obteve %d", len(classificacao))
		}
		if classificacao[0].NotaFinal != 85 {
			t.Errorf("nota final deve ser 85, obteve %d", classificacao[0].NotaFinal)
		}
		if classificacao[0].Status != "aprovada" {
			t.Errorf("status deve ser aprovada, obteve %s", classificacao[0].Status)
		}
	})

	t.Run("deve retornar lista vazia quando nao ha propostas em avaliacao", func(t *testing.T) {
		propostaRepo := persistencia.NovoRepositorioDePropostaMemoria()
		uc := NovoFinalizarAvaliacaoDoEdital(propostaRepo)

		classificacao, err := uc.Executar(context.Background(), 99, dto.FinalizarAvaliacaoEntrada{NotaDeCorte: 70})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(classificacao) != 0 {
			t.Errorf("esperava 0 itens, obteve %d", len(classificacao))
		}
	})
}
