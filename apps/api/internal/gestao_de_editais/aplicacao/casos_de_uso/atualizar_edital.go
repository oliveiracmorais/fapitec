package casos_de_uso

import (
	"context"
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
)

type AtualizarEdital struct {
	repo repositorios.RepositorioDeEdital
}

func NovoAtualizarEdital(repo repositorios.RepositorioDeEdital) *AtualizarEdital {
	return &AtualizarEdital{repo: repo}
}

func (a *AtualizarEdital) Executar(ctx context.Context, id int64, entrada dto.AtualizarEditalEntrada) (*dto.EditalSaida, error) {
	edital, err := a.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar edital: %w", err)
	}
	if edital == nil {
		return nil, fmt.Errorf("edital nao encontrado")
	}

	if entrada.Nome != "" {
		edital.Nome = entrada.Nome
	}
	if entrada.Descricao != "" {
		edital.Descricao = entrada.Descricao
	}
	if entrada.DataInicio != "" {
		dataInicio, err := time.Parse("2006-01-02", entrada.DataInicio)
		if err != nil {
			return nil, fmt.Errorf("data de inicio invalida: %w", err)
		}
		edital.DataInicio = dataInicio
	}
	if entrada.DataFim != "" {
		dataFim, err := time.Parse("2006-01-02", entrada.DataFim)
		if err != nil {
			return nil, fmt.Errorf("data de fim invalida: %w", err)
		}
		edital.DataFim = dataFim
	}
	if entrada.Status != "" {
		status, err := objetos_de_valor.NovoStatusEdital(entrada.Status)
		if err != nil {
			return nil, err
		}
		edital.Status = status
	}
	if entrada.TipoChamada != "" {
		edital.TipoChamada = entrada.TipoChamada
	}

	if !edital.DataFim.IsZero() && !edital.DataInicio.IsZero() && edital.DataInicio.After(edital.DataFim) {
		return nil, fmt.Errorf("data de inicio nao pode ser posterior a data de fim")
	}

	if err := a.repo.Atualizar(ctx, edital); err != nil {
		return nil, fmt.Errorf("erro ao atualizar edital: %w", err)
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
