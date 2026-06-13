package casos_de_uso

import (
	"context"
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type AtribuirEdital struct {
	repoAvaliador  repositorios.RepositorioDeAvaliador
	repoAtribuicao repositorios.RepositorioDeAtribuicao
}

func NovoAtribuirEdital(repoAvaliador repositorios.RepositorioDeAvaliador, repoAtribuicao repositorios.RepositorioDeAtribuicao) *AtribuirEdital {
	return &AtribuirEdital{repoAvaliador: repoAvaliador, repoAtribuicao: repoAtribuicao}
}

func (uc *AtribuirEdital) Executar(ctx context.Context, avaliadorID int64, entrada dto.AtribuirEditalEntrada) (*dto.AtribuicaoSaida, error) {
	avaliador, err := uc.repoAvaliador.BuscarPorID(ctx, avaliadorID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar avaliador: %w", err)
	}
	if avaliador == nil {
		return nil, fmt.Errorf("avaliador nao encontrado")
	}

	dataInicio, err := time.Parse("2006-01-02T15:04:05Z07:00", entrada.DataInicio)
	if err != nil {
		dataInicio, err = time.Parse("2006-01-02", entrada.DataInicio)
		if err != nil {
			return nil, fmt.Errorf("data de inicio invalida")
		}
	}

	dataFim, err := time.Parse("2006-01-02T15:04:05Z07:00", entrada.DataFim)
	if err != nil {
		dataFim, err = time.Parse("2006-01-02", entrada.DataFim)
		if err != nil {
			return nil, fmt.Errorf("data de fim invalida")
		}
	}

	atribuicao, err := entidades.NovaAtribuicao(entidades.NovaAtribuicaoParams{
		AvaliadorID: avaliadorID,
		EditalID:    entrada.EditalID,
		DataInicio:  dataInicio,
		DataFim:     dataFim,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar atribuicao: %w", err)
	}

	if err := uc.repoAtribuicao.Criar(ctx, atribuicao); err != nil {
		return nil, fmt.Errorf("erro ao salvar atribuicao: %w", err)
	}

	return paraAtribuicaoSaida(atribuicao), nil
}
