package casos_de_uso

import (
	"context"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
)

type ListarEditais struct {
	repo repositorios.RepositorioDeEdital
}

func NovoListarEditais(repo repositorios.RepositorioDeEdital) *ListarEditais {
	return &ListarEditais{repo: repo}
}

type FiltrosListarEditais struct {
	Titulo      string
	Status      string
	TipoChamada string
}

func (l *ListarEditais) Executar(ctx context.Context, filtros FiltrosListarEditais) (*dto.ListarEditaisSaida, error) {
	repoFiltros := repositorios.FiltrosListarEditais{
		Titulo:      filtros.Titulo,
		Status:      filtros.Status,
		TipoChamada: filtros.TipoChamada,
	}

	editais, err := l.repo.Listar(ctx, repoFiltros)
	if err != nil {
		return nil, err
	}

	saida := make([]dto.EditalSaida, 0, len(editais))
	for _, e := range editais {
		saida = append(saida, *paraEditalSaida(e))
	}

	return &dto.ListarEditaisSaida{
		Editais: saida,
		Total:   len(saida),
	}, nil
}
