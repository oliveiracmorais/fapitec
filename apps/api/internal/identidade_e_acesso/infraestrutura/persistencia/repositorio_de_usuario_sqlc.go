package persistencia

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia/sqlc"
)

type RepositorioDeUsuarioSQLC struct {
	queries *sqlcpersistencia.Queries
}

func NovoRepositorioDeUsuarioSQLC(queries *sqlcpersistencia.Queries) *RepositorioDeUsuarioSQLC {
	return &RepositorioDeUsuarioSQLC{queries: queries}
}

func (r *RepositorioDeUsuarioSQLC) Inserir(ctx context.Context, usuario *entidades.Usuario) error {
	params := sqlcpersistencia.InserirUsuarioParams{
		Nome:          usuario.Nome,
		Cpf:           usuario.CPF,
		Email:         usuario.Email.String(),
		SenhaHash:     usuario.SenhaHash.String(),
		EhEstrangeiro: usuario.Estrangeiro,
	}
	result, err := r.queries.InserirUsuario(ctx, params)
	if err != nil {
		return err
	}
	usuario.ID = result.ID
	usuario.CriadoEm = result.CriadoEm.Time
	return nil
}

func (r *RepositorioDeUsuarioSQLC) BuscarPorCPF(ctx context.Context, cpf string) (*entidades.Usuario, error) {
	result, err := r.queries.BuscarUsuarioPorCPF(ctx, cpf)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.sqlcParaDominio(result), nil
}

func (r *RepositorioDeUsuarioSQLC) BuscarPorEmail(ctx context.Context, email string) (*entidades.Usuario, error) {
	result, err := r.queries.BuscarUsuarioPorEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.sqlcParaDominio(result), nil
}

func (r *RepositorioDeUsuarioSQLC) AtualizarTentativas(ctx context.Context, id int64, tentativas int, bloqueadoAte *string) error {
	params := sqlcpersistencia.AtualizarTentativasUsuarioParams{
		ID:         id,
		Tentativas: int32(tentativas),
	}
	if bloqueadoAte != nil {
		t, err := time.Parse(time.RFC3339, *bloqueadoAte)
		if err == nil {
			params.BloqueadoAte = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}
	return r.queries.AtualizarTentativasUsuario(ctx, params)
}

func (r *RepositorioDeUsuarioSQLC) AtualizarSenha(ctx context.Context, id int64, senhaHash string) error {
	params := sqlcpersistencia.AtualizarSenhaUsuarioParams{
		ID:        id,
		SenhaHash: senhaHash,
	}
	return r.queries.AtualizarSenhaUsuario(ctx, params)
}

func (r *RepositorioDeUsuarioSQLC) sqlcParaDominio(s sqlcpersistencia.Usuario) *entidades.Usuario {
	email, _ := objetos_de_valor.NovoEmail(s.Email)

	var bloqueadoAte *time.Time
	if s.BloqueadoAte.Valid {
		bloqueadoAte = &s.BloqueadoAte.Time
	}

	return &entidades.Usuario{
		ID:           s.ID,
		Nome:         s.Nome,
		CPF:          s.Cpf,
		Estrangeiro:  s.EhEstrangeiro,
		Email:        email,
		SenhaHash:    objetos_de_valor.NovaSenhaHash(s.SenhaHash),
		Tentativas:   int(s.Tentativas),
		BloqueadoAte: bloqueadoAte,
		CriadoEm:     s.CriadoEm.Time,
	}
}
