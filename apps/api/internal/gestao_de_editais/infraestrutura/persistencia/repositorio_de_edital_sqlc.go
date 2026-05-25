package persistencia

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
	sqlcpersistencia "github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/infraestrutura/persistencia/sqlc"
)

type RepositorioDeEditalSQLC struct {
	queries *sqlcpersistencia.Queries
}

func NovoRepositorioDeEditalSQLC(queries *sqlcpersistencia.Queries) *RepositorioDeEditalSQLC {
	return &RepositorioDeEditalSQLC{queries: queries}
}

func (r *RepositorioDeEditalSQLC) Criar(ctx context.Context, edital *entidades.Edital) error {
	dataInicio := pgtype.Date{}
	if err := dataInicio.Scan(edital.DataInicio); err != nil {
		return err
	}
	dataFim := pgtype.Date{}
	if err := dataFim.Scan(edital.DataFim); err != nil {
		return err
	}

	params := sqlcpersistencia.InserirEditalParams{
		Nome:        edital.Nome,
		Descricao:   edital.Descricao,
		DataInicio:  dataInicio,
		DataFim:     dataFim,
		Status:      edital.Status.String(),
		TipoChamada: edital.TipoChamada,
		NotaDeCorte: int32(edital.NotaDeCorte),
		ValorGlobal: edital.ValorGlobal,
	}
	result, err := r.queries.InserirEdital(ctx, params)
	if err != nil {
		return err
	}
	edital.ID = result.ID
	edital.CriadoEm = result.CriadoEm.Time
	return nil
}

func (r *RepositorioDeEditalSQLC) BuscarPorID(ctx context.Context, id int64) (*entidades.Edital, error) {
	result, err := r.queries.BuscarEditalPorID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.sqlcParaDominio(result), nil
}

func (r *RepositorioDeEditalSQLC) Listar(ctx context.Context, filtros repositorios.FiltrosListarEditais) ([]*entidades.Edital, error) {
	params := sqlcpersistencia.ListarEditaisParams{
		Column1: filtros.Titulo,
		Column2: filtros.Status,
		Column3: filtros.TipoChamada,
	}
	results, err := r.queries.ListarEditais(ctx, params)
	if err != nil {
		return nil, err
	}
	editais := make([]*entidades.Edital, 0, len(results))
	for _, result := range results {
		editais = append(editais, r.sqlcParaDominio(result))
	}
	return editais, nil
}

func (r *RepositorioDeEditalSQLC) Atualizar(ctx context.Context, edital *entidades.Edital) error {
	dataInicio := pgtype.Date{}
	if err := dataInicio.Scan(edital.DataInicio); err != nil {
		return err
	}
	dataFim := pgtype.Date{}
	if err := dataFim.Scan(edital.DataFim); err != nil {
		return err
	}

	params := sqlcpersistencia.AtualizarEditalParams{
		ID:          edital.ID,
		Nome:        edital.Nome,
		Descricao:   edital.Descricao,
		DataInicio:  dataInicio,
		DataFim:     dataFim,
		Status:      edital.Status.String(),
		TipoChamada: edital.TipoChamada,
		NotaDeCorte: int32(edital.NotaDeCorte),
		ValorGlobal: edital.ValorGlobal,
	}
	return r.queries.AtualizarEdital(ctx, params)
}

func (r *RepositorioDeEditalSQLC) Deletar(ctx context.Context, id int64) error {
	return r.queries.DeletarEdital(ctx, id)
}

func (r *RepositorioDeEditalSQLC) sqlcParaDominio(s sqlcpersistencia.Editai) *entidades.Edital {
	status, err := objetos_de_valor.NovoStatusEdital(s.Status)
	if err != nil {
		status = objetos_de_valor.StatusEditalAtivo
	}

	var dataInicio time.Time
	if s.DataInicio.Valid {
		dataInicio = s.DataInicio.Time
	}
	var dataFim time.Time
	if s.DataFim.Valid {
		dataFim = s.DataFim.Time
	}

	return &entidades.Edital{
		ID:          s.ID,
		Nome:        s.Nome,
		Descricao:   s.Descricao,
		DataInicio:  dataInicio,
		DataFim:     dataFim,
		Status:      status,
		TipoChamada: s.TipoChamada,
		NotaDeCorte: int(s.NotaDeCorte),
		ValorGlobal: s.ValorGlobal,
		CriadoEm:    s.CriadoEm.Time,
	}
}
