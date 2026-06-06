package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type ListarPropostas struct {
	repo repositorios.RepositorioDeProposta
}

func NovoListarPropostas(repo repositorios.RepositorioDeProposta) *ListarPropostas {
	return &ListarPropostas{repo: repo}
}

func (uc *ListarPropostas) Executar(ctx context.Context, filtros repositorios.FiltrosListarPropostas) ([]dto.PropostaSaida, error) {
	propostas, err := uc.repo.Listar(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar propostas: %w", err)
	}

	saida := make([]dto.PropostaSaida, len(propostas))
	for i, p := range propostas {
		saida[i] = *paraPropostaSaida(p)
	}
	return saida, nil
}
