package persistencia

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia/sqlc"
)

type RepositorioDeTokenRedefinicaoSQLC struct {
	queries *sqlcpersistencia.Queries
}

func NovoRepositorioDeTokenRedefinicaoSQLC(queries *sqlcpersistencia.Queries) *RepositorioDeTokenRedefinicaoSQLC {
	return &RepositorioDeTokenRedefinicaoSQLC{queries: queries}
}

func (r *RepositorioDeTokenRedefinicaoSQLC) Inserir(ctx context.Context, token objetos_de_valor.TokenRedefinicao) error {
	params := sqlcpersistencia.InserirTokenRedefinicaoParams{
		UsuarioID:  token.UsuarioID,
		Token:      token.TokenHash,
		ExpiradoEm: pgtype.Timestamptz{Time: token.ExpiradoEm, Valid: true},
	}
	return r.queries.InserirTokenRedefinicao(ctx, params)
}

func (r *RepositorioDeTokenRedefinicaoSQLC) BuscarPorToken(ctx context.Context, token string) (*objetos_de_valor.TokenRedefinicao, error) {
	hash := objetos_de_valor.HashToken(token)
	result, err := r.queries.BuscarTokenRedefinicao(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	t := &objetos_de_valor.TokenRedefinicao{
		TokenHash:  result.Token,
		UsuarioID:  result.UsuarioID,
		ExpiradoEm: result.ExpiradoEm.Time,
	}
	return t, nil
}

func (r *RepositorioDeTokenRedefinicaoSQLC) Remover(ctx context.Context, token string) error {
	hash := objetos_de_valor.HashToken(token)
	return r.queries.ConsumirTokenRedefinicao(ctx, hash)
}
