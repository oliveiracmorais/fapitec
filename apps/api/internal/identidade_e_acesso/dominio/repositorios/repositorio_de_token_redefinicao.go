package repositorios

import (
	"context"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
)

type RepositorioDeTokenRedefinicao interface {
	Inserir(ctx context.Context, token objetos_de_valor.TokenRedefinicao) error
	BuscarPorToken(ctx context.Context, token string) (*objetos_de_valor.TokenRedefinicao, error)
	Remover(ctx context.Context, token string) error
}
