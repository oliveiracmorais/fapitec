package casos_de_uso

import (
	"context"
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
)

type CriarEdital struct {
	repo repositorios.RepositorioDeEdital
}

func NovoEdital(repo repositorios.RepositorioDeEdital) *CriarEdital {
	return &CriarEdital{repo: repo}
}

func (c *CriarEdital) Executar(ctx context.Context, entrada dto.CriarEditalEntrada) (*dto.EditalSaida, error) {
	dataInicio, err := time.Parse("2006-01-02", entrada.DataInicio)
	if err != nil {
		return nil, fmt.Errorf("data de inicio invalida: %w", err)
	}

	dataFim, err := time.Parse("2006-01-02", entrada.DataFim)
	if err != nil {
		return nil, fmt.Errorf("data de fim invalida: %w", err)
	}

	params := entidades.NovoEditalParams{
		Nome:        entrada.Nome,
		Descricao:   entrada.Descricao,
		DataInicio:  dataInicio,
		DataFim:     dataFim,
		TipoChamada: entrada.TipoChamada,
	}

	edital, err := entidades.NovoEdital(params)
	if err != nil {
		return nil, err
	}

	if err := c.repo.Criar(ctx, edital); err != nil {
		return nil, fmt.Errorf("erro ao criar edital: %w", err)
	}

	return &dto.EditalSaida{
		ID:          edital.ID,
		Nome:        edital.Nome,
		Descricao:   edital.Descricao,
		DataInicio:  edital.DataInicio.Format("2006-01-02"),
		DataFim:     edital.DataFim.Format("2006-01-02"),
		Status:      edital.Status.String(),
		TipoChamada: edital.TipoChamada,
		CriadoEm:    edital.CriadoEm.Format(time.RFC3339),
	}, nil
}
