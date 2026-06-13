package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type ListarAtribuicoes struct {
	repo repositorios.RepositorioDeAtribuicao
}

func NovoListarAtribuicoes(repo repositorios.RepositorioDeAtribuicao) *ListarAtribuicoes {
	return &ListarAtribuicoes{repo: repo}
}

func (uc *ListarAtribuicoes) ExecutarPorAvaliador(ctx context.Context, avaliadorID int64) ([]dto.AtribuicaoSaida, error) {
	atribuicoes, err := uc.repo.ListarPorAvaliador(ctx, avaliadorID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar atribuicoes: %w", err)
	}

	saida := make([]dto.AtribuicaoSaida, 0, len(atribuicoes))
	for _, a := range atribuicoes {
		saida = append(saida, *paraAtribuicaoSaida(a))
	}
	return saida, nil
}

func (uc *ListarAtribuicoes) ExecutarPorEdital(ctx context.Context, editalID int64) ([]dto.AtribuicaoSaida, error) {
	atribuicoes, err := uc.repo.ListarPorEdital(ctx, editalID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar atribuicoes: %w", err)
	}

	saida := make([]dto.AtribuicaoSaida, 0, len(atribuicoes))
	for _, a := range atribuicoes {
		saida = append(saida, *paraAtribuicaoSaida(a))
	}
	return saida, nil
}
