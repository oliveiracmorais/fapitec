package repositorios

import (
	"context"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/entidades"
)

type FiltrosListarEditais struct {
	Titulo      string
	Status      string
	TipoChamada string
}

type RepositorioDeEdital interface {
	Criar(ctx context.Context, edital *entidades.Edital) error
	BuscarPorID(ctx context.Context, id int64) (*entidades.Edital, error)
	Listar(ctx context.Context, filtros FiltrosListarEditais) ([]*entidades.Edital, error)
	Atualizar(ctx context.Context, edital *entidades.Edital) error
}
