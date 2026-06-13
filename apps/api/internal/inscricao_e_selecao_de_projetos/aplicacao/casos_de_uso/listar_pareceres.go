package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type ListarPareceres struct {
	propostaRepo   repositorios.RepositorioDeProposta
	atribuicaoRepo repositorios.RepositorioDeAtribuicao
}

func NovoListarPareceres(propostaRepo repositorios.RepositorioDeProposta, atribuicaoRepo repositorios.RepositorioDeAtribuicao) *ListarPareceres {
	return &ListarPareceres{
		propostaRepo:   propostaRepo,
		atribuicaoRepo: atribuicaoRepo,
	}
}

func (uc *ListarPareceres) ExecutarPorProposta(ctx context.Context, propostaID int64) ([]*dto.ParecerAnonimizadoSaida, error) {
	proposta, err := uc.propostaRepo.BuscarPorID(ctx, propostaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar proposta: %w", err)
	}
	if proposta == nil {
		return nil, fmt.Errorf("proposta nao encontrada")
	}

	pareceres, err := uc.propostaRepo.ListarPareceresPorProposta(ctx, propostaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pareceres: %w", err)
	}

	saida := make([]*dto.ParecerAnonimizadoSaida, len(pareceres))
	for i, p := range pareceres {
		hash := fmt.Sprintf("avaliador-%d", p.AvaliadorID)
		if h, err := uc.atribuicaoRepo.BuscarHashAnonimizacao(ctx, p.AvaliadorID, proposta.EditalID); err == nil && h != "" {
			hash = h
		}

		saida[i] = &dto.ParecerAnonimizadoSaida{
			ID:            p.ID,
			PropostaID:    p.PropostaID,
			Etapa:         p.Etapa,
			HashAvaliador: hash,
			Nota:          p.Nota,
			ParecerTexto:  p.ParecerTexto,
			Data:          p.Data.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return saida, nil
}
