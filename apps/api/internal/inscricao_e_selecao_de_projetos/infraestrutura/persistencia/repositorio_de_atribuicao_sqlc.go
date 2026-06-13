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

type RepositorioDeAtribuicaoSQLC struct {
	queries *sqlcpersistencia.Queries
}

func NovoRepositorioDeAtribuicaoSQLC(queries *sqlcpersistencia.Queries) *RepositorioDeAtribuicaoSQLC {
	return &RepositorioDeAtribuicaoSQLC{queries: queries}
}

func (r *RepositorioDeAtribuicaoSQLC) Criar(ctx context.Context, atribuicao *entidades.AtribuicaoEdital) error {
	params := sqlcpersistencia.InserirAtribuicaoParams{
		AvaliadorID:      atribuicao.AvaliadorID,
		EditalID:         atribuicao.EditalID,
		DataInicio:       pgtype.Timestamptz{Time: atribuicao.DataInicio, Valid: true},
		DataFim:          pgtype.Timestamptz{Time: atribuicao.DataFim, Valid: true},
		StatusConvite:    atribuicao.StatusConvite.String(),
		HashAnonimizacao: atribuicao.HashAnonimizacao.String(),
		CriadoEm:         pgtype.Timestamptz{Time: atribuicao.CriadoEm, Valid: true},
	}

	id, err := r.queries.InserirAtribuicao(ctx, params)
	if err != nil {
		return err
	}
	atribuicao.ID = id
	return nil
}

func (r *RepositorioDeAtribuicaoSQLC) BuscarPorID(ctx context.Context, id int64) (*entidades.AtribuicaoEdital, error) {
	result, err := r.queries.BuscarAtribuicaoPorID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return sqlcAtribuicaoParaDominio(result), nil
}

func (r *RepositorioDeAtribuicaoSQLC) ListarPorAvaliador(ctx context.Context, avaliadorID int64) ([]*entidades.AtribuicaoEdital, error) {
	results, err := r.queries.ListarAtribuicoesPorAvaliador(ctx, avaliadorID)
	if err != nil {
		return nil, err
	}

	atribuicoes := make([]*entidades.AtribuicaoEdital, 0, len(results))
	for _, r := range results {
		atribuicoes = append(atribuicoes, sqlcAtribuicaoParaDominio(r))
	}
	return atribuicoes, nil
}

func (r *RepositorioDeAtribuicaoSQLC) ListarPorEdital(ctx context.Context, editalID int64) ([]*entidades.AtribuicaoEdital, error) {
	results, err := r.queries.ListarAtribuicoesPorEdital(ctx, editalID)
	if err != nil {
		return nil, err
	}

	atribuicoes := make([]*entidades.AtribuicaoEdital, 0, len(results))
	for _, r := range results {
		atribuicoes = append(atribuicoes, sqlcAtribuicaoParaDominio(r))
	}
	return atribuicoes, nil
}

func (r *RepositorioDeAtribuicaoSQLC) Listar(ctx context.Context, filtros repositorios.FiltrosListarAtribuicoes) ([]*entidades.AtribuicaoEdital, error) {
	results, err := r.queries.ListarAtribuicoes(ctx, sqlcpersistencia.ListarAtribuicoesParams{
		Column1: filtros.AvaliadorID,
		Column2: filtros.EditalID,
		Column3: filtros.Status,
	})
	if err != nil {
		return nil, err
	}

	atribuicoes := make([]*entidades.AtribuicaoEdital, 0, len(results))
	for _, r := range results {
		atribuicoes = append(atribuicoes, sqlcAtribuicaoParaDominio(r))
	}
	return atribuicoes, nil
}

func (r *RepositorioDeAtribuicaoSQLC) Atualizar(ctx context.Context, atribuicao *entidades.AtribuicaoEdital) error {
	return r.queries.AtualizarAtribuicaoStatus(ctx, sqlcpersistencia.AtualizarAtribuicaoStatusParams{
		ID:            atribuicao.ID,
		StatusConvite: atribuicao.StatusConvite.String(),
	})
}

func (r *RepositorioDeAtribuicaoSQLC) ContarAtivasPorAvaliador(ctx context.Context, avaliadorID int64) (int64, error) {
	return r.queries.ContarAtribuicoesAtivasPorAvaliador(ctx, avaliadorID)
}

func (r *RepositorioDeAtribuicaoSQLC) BuscarHashAnonimizacao(_ context.Context, _ int64, _ int64) (string, error) {
	return "", errors.New("BuscarHashAnonimizacao: implementacao SQL pendente - usar repositorio em memoria")
}

func sqlcAtribuicaoParaDominio(s sqlcpersistencia.AtribuicoesEditai) *entidades.AtribuicaoEdital {
	status, _ := objetos_de_valor.NovoStatusConvite(s.StatusConvite)
	hash, _ := objetos_de_valor.HashAnonimizacaoExistente(s.HashAnonimizacao)
	return &entidades.AtribuicaoEdital{
		ID:               s.ID,
		AvaliadorID:      s.AvaliadorID,
		EditalID:         s.EditalID,
		DataInicio:       s.DataInicio.Time,
		DataFim:          s.DataFim.Time,
		StatusConvite:    status,
		HashAnonimizacao: hash,
		CriadoEm:         s.CriadoEm.Time,
	}
}
