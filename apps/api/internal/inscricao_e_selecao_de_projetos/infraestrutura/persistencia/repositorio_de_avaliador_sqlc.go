package persistencia

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
	sqlcpersistencia "github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia/sqlc"
)

type RepositorioDeAvaliadorSQLC struct {
	queries *sqlcpersistencia.Queries
}

func NovoRepositorioDeAvaliadorSQLC(queries *sqlcpersistencia.Queries) *RepositorioDeAvaliadorSQLC {
	return &RepositorioDeAvaliadorSQLC{queries: queries}
}

func (r *RepositorioDeAvaliadorSQLC) Criar(ctx context.Context, avaliador *entidades.Avaliador) error {
	params := sqlcpersistencia.InserirAvaliadorParams{
		UsuarioID:         avaliador.UsuarioID,
		Nome:              avaliador.Nome,
		Cpf:               avaliador.CPF,
		Email:             avaliador.Email,
		TitulacaoMaxima:   avaliador.TitulacaoMaxima,
		AreaConhecimento:  avaliador.AreaConhecimento,
		Instituicao:       avaliador.Instituicao,
		CurriculoResumido: avaliador.CurriculoResumido,
		Estado:            avaliador.Estado.String(),
		DataCadastro:      pgtype.Timestamptz{Time: avaliador.DataCadastro, Valid: true},
		DataAtualizacao:   pgtype.Timestamptz{Time: avaliador.DataAtualizacao, Valid: true},
	}

	id, err := r.queries.InserirAvaliador(ctx, params)
	if err != nil {
		return err
	}
	avaliador.ID = id
	return nil
}

func (r *RepositorioDeAvaliadorSQLC) BuscarPorID(ctx context.Context, id int64) (*entidades.Avaliador, error) {
	result, err := r.queries.BuscarAvaliadorPorID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return sqlcAvaliadorParaDominio(result), nil
}

func (r *RepositorioDeAvaliadorSQLC) BuscarPorCPF(ctx context.Context, cpf string) (*entidades.Avaliador, error) {
	result, err := r.queries.BuscarAvaliadorPorCPF(ctx, cpf)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return sqlcAvaliadorParaDominio(result), nil
}

func (r *RepositorioDeAvaliadorSQLC) BuscarPorUsuarioID(ctx context.Context, usuarioID int64) (*entidades.Avaliador, error) {
	result, err := r.queries.BuscarAvaliadorPorUsuarioID(ctx, usuarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return sqlcAvaliadorParaDominio(result), nil
}

func (r *RepositorioDeAvaliadorSQLC) Listar(ctx context.Context, filtros repositorios.FiltrosListarAvaliadores) ([]*entidades.Avaliador, error) {
	results, err := r.queries.ListarAvaliadores(ctx, sqlcpersistencia.ListarAvaliadoresParams{
		Column1: filtros.Nome,
		Column2: filtros.CPF,
		Column3: filtros.AreaConhecimento,
		Column4: filtros.Estado,
	})
	if err != nil {
		return nil, err
	}

	avaliadores := make([]*entidades.Avaliador, 0, len(results))
	for _, r := range results {
		avaliadores = append(avaliadores, sqlcAvaliadorParaDominio(r))
	}
	return avaliadores, nil
}

func (r *RepositorioDeAvaliadorSQLC) Atualizar(ctx context.Context, avaliador *entidades.Avaliador) error {
	return r.queries.AtualizarAvaliador(ctx, sqlcpersistencia.AtualizarAvaliadorParams{
		ID:                avaliador.ID,
		Nome:              avaliador.Nome,
		Cpf:               avaliador.CPF,
		Email:             avaliador.Email,
		TitulacaoMaxima:   avaliador.TitulacaoMaxima,
		AreaConhecimento:  avaliador.AreaConhecimento,
		Instituicao:       avaliador.Instituicao,
		CurriculoResumido: avaliador.CurriculoResumido,
		Estado:            avaliador.Estado.String(),
	})
}

func (r *RepositorioDeAvaliadorSQLC) ContarPropostasPorAvaliador(ctx context.Context, avaliadorID int64) (int64, error) {
	return r.queries.ContarPropostasPorAvaliador(ctx, avaliadorID)
}

func sqlcAvaliadorParaDominio(s sqlcpersistencia.Avaliadore) *entidades.Avaliador {
	estado, _ := objetos_de_valor.NovoEstadoAvaliador(s.Estado)
	return &entidades.Avaliador{
		ID:                s.ID,
		UsuarioID:         s.UsuarioID,
		Nome:              s.Nome,
		CPF:               s.Cpf,
		Email:             s.Email,
		TitulacaoMaxima:   s.TitulacaoMaxima,
		AreaConhecimento:  s.AreaConhecimento,
		Instituicao:       s.Instituicao,
		CurriculoResumido: s.CurriculoResumido,
		Estado:            estado,
		DataCadastro:      s.DataCadastro.Time,
		DataAtualizacao:   s.DataAtualizacao.Time,
	}
}
