package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type ListarAvaliadores struct {
	repo repositorios.RepositorioDeAvaliador
}

func NovoListarAvaliadores(repo repositorios.RepositorioDeAvaliador) *ListarAvaliadores {
	return &ListarAvaliadores{repo: repo}
}

func (uc *ListarAvaliadores) Executar(ctx context.Context, filtros repositorios.FiltrosListarAvaliadores) ([]dto.AvaliadorSaida, error) {
	avaliadores, err := uc.repo.Listar(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar avaliadores: %w", err)
	}

	saida := make([]dto.AvaliadorSaida, 0, len(avaliadores))
	for _, a := range avaliadores {
		totalPropostas, _ := uc.repo.ContarPropostasPorAvaliador(ctx, a.ID)
		saida = append(saida, *paraAvaliadorSaida(a, totalPropostas, 0))
	}
	return saida, nil
}

type VisualizarAvaliador struct {
	repo              repositorios.RepositorioDeAvaliador
	repoAtribuicao    repositorios.RepositorioDeAtribuicao
}

func NovoVisualizarAvaliador(repo repositorios.RepositorioDeAvaliador, repoAtribuicao repositorios.RepositorioDeAtribuicao) *VisualizarAvaliador {
	return &VisualizarAvaliador{repo: repo, repoAtribuicao: repoAtribuicao}
}

func (uc *VisualizarAvaliador) Executar(ctx context.Context, id int64) (*dto.AvaliadorSaida, error) {
	avaliador, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar avaliador: %w", err)
	}
	if avaliador == nil {
		return nil, fmt.Errorf("avaliador nao encontrado")
	}

	totalPropostas, _ := uc.repo.ContarPropostasPorAvaliador(ctx, id)
	atribuicoesAtivas, _ := uc.repoAtribuicao.ContarAtivasPorAvaliador(ctx, id)

	return paraAvaliadorSaida(avaliador, totalPropostas, atribuicoesAtivas), nil
}
