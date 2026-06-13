package casos_de_uso

import (
	"context"
	"fmt"
	"sort"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type FinalizarAvaliacaoDoEdital struct {
	propostaRepo repositorios.RepositorioDeProposta
}

func NovoFinalizarAvaliacaoDoEdital(propostaRepo repositorios.RepositorioDeProposta) *FinalizarAvaliacaoDoEdital {
	return &FinalizarAvaliacaoDoEdital{propostaRepo: propostaRepo}
}

func (uc *FinalizarAvaliacaoDoEdital) Executar(ctx context.Context, editalID int64, entrada dto.FinalizarAvaliacaoEntrada) ([]*dto.ClassificacaoSaida, error) {
	propostas, err := uc.propostaRepo.ListarPropostasPorEdital(ctx, editalID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar propostas: %w", err)
	}

	var classificacao []*dto.ClassificacaoSaida

	for _, proposta := range propostas {
		if proposta.Status.String() != "em_avaliacao" {
			continue
		}

		if err := proposta.FinalizarAvaliacao(entrada.NotaDeCorte); err != nil {
			continue
		}

		if err := uc.propostaRepo.Atualizar(ctx, proposta); err != nil {
			return nil, fmt.Errorf("erro ao atualizar proposta %d: %w", proposta.ID, err)
		}

		if err := uc.propostaRepo.AtualizarStatus(ctx, proposta.ID, proposta.Status.String()); err != nil {
			return nil, fmt.Errorf("erro ao atualizar status da proposta %d: %w", proposta.ID, err)
		}

		classificacao = append(classificacao, &dto.ClassificacaoSaida{
			PropostaID: proposta.ID,
			Protocolo:  proposta.Protocolo.String(),
			NotaFinal:  proposta.CalcularNotaFinal(),
			Status:     proposta.Status.String(),
		})
	}

	sort.Slice(classificacao, func(i, j int) bool {
		return classificacao[i].NotaFinal > classificacao[j].NotaFinal
	})

	if classificacao == nil {
		classificacao = []*dto.ClassificacaoSaida{}
	}

	return classificacao, nil
}
