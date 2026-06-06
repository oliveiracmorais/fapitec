package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type VisualizarProposta struct {
	repo repositorios.RepositorioDeProposta
}

func NovoVisualizarProposta(repo repositorios.RepositorioDeProposta) *VisualizarProposta {
	return &VisualizarProposta{repo: repo}
}

func (uc *VisualizarProposta) Executar(ctx context.Context, id int64) (*dto.PropostaSaida, error) {
	proposta, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar proposta: %w", err)
	}
	if proposta == nil {
		return nil, fmt.Errorf("proposta nao encontrada")
	}
	return paraPropostaSaida(proposta), nil
}
