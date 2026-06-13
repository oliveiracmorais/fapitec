package repositorios

import (
	"context"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
)

type FiltrosListarAvaliadores struct {
	Nome            string
	CPF             string
	AreaConhecimento string
	Estado          string
}

type RepositorioDeAvaliador interface {
	Criar(ctx context.Context, avaliador *entidades.Avaliador) error
	BuscarPorID(ctx context.Context, id int64) (*entidades.Avaliador, error)
	BuscarPorCPF(ctx context.Context, cpf string) (*entidades.Avaliador, error)
	BuscarPorUsuarioID(ctx context.Context, usuarioID int64) (*entidades.Avaliador, error)
	Listar(ctx context.Context, filtros FiltrosListarAvaliadores) ([]*entidades.Avaliador, error)
	Atualizar(ctx context.Context, avaliador *entidades.Avaliador) error
	ContarPropostasPorAvaliador(ctx context.Context, avaliadorID int64) (int64, error)
}
