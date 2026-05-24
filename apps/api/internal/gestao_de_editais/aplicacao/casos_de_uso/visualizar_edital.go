package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
)

type VisualizarEdital struct {
	repo repositorios.RepositorioDeEdital
}

func NovoVisualizarEdital(repo repositorios.RepositorioDeEdital) *VisualizarEdital {
	return &VisualizarEdital{repo: repo}
}

func (v *VisualizarEdital) Executar(ctx context.Context, id int64) (*dto.EditalSaida, error) {
	edital, err := v.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar edital: %w", err)
	}
	if edital == nil {
		return nil, fmt.Errorf("edital nao encontrado")
	}

	return &dto.EditalSaida{
		ID:          edital.ID,
		Nome:        edital.Nome,
		Descricao:   edital.Descricao,
		DataInicio:  edital.DataInicio.Format("2006-01-02"),
		DataFim:     edital.DataFim.Format("2006-01-02"),
		Status:      edital.Status.String(),
		TipoChamada: edital.TipoChamada,
		CriadoEm:    edital.CriadoEm.Format("2006-01-02T15:04:05-07:00"),
	}, nil
}
