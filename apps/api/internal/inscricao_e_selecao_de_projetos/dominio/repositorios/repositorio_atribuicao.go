package repositorios

import (
	"context"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
)

type FiltrosListarAtribuicoes struct {
	AvaliadorID int64
	EditalID    int64
	Status      string
}

type RepositorioDeAtribuicao interface {
	Criar(ctx context.Context, atribuicao *entidades.AtribuicaoEdital) error
	BuscarPorID(ctx context.Context, id int64) (*entidades.AtribuicaoEdital, error)
	ListarPorAvaliador(ctx context.Context, avaliadorID int64) ([]*entidades.AtribuicaoEdital, error)
	ListarPorEdital(ctx context.Context, editalID int64) ([]*entidades.AtribuicaoEdital, error)
	Listar(ctx context.Context, filtros FiltrosListarAtribuicoes) ([]*entidades.AtribuicaoEdital, error)
	Atualizar(ctx context.Context, atribuicao *entidades.AtribuicaoEdital) error
	ContarAtivasPorAvaliador(ctx context.Context, avaliadorID int64) (int64, error)
	BuscarHashAnonimizacao(ctx context.Context, avaliadorID int64, editalID int64) (string, error)
}
