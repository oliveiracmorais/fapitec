package persistencia

import (
	"context"
	"encoding/json"
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

	relatoriosJSON, err := json.Marshal(edital.RelatoriosExigidos)
	if err != nil {
		return err
	}
	porteJSON, err := json.Marshal(edital.PorteEmpresa)
	if err != nil {
		return err
	}
	enquadramentoJSON, err := json.Marshal(edital.EnquadramentoEmpresa)
	if err != nil {
		return err
	}
	documentosJSON, err := json.Marshal(edital.DocumentosObrigatorios)
	if err != nil {
		return err
	}

	params := sqlcpersistencia.InserirEditalParams{
		Nome:                      edital.Nome,
		Descricao:                 edital.Descricao,
		DataInicio:                dataInicio,
		DataFim:                   dataFim,
		Status:                    edital.Status.String(),
		TipoChamada:               edital.TipoChamada.String(),
		NotaDeCorte:               int32(edital.NotaDeCorte),
		ValorGlobal:               edital.ValorGlobal,
		ModeloFormulario:          int32(edital.ModeloFormulario),
		RelatoriosExigidos:        relatoriosJSON,
		TituloMinimoElegibilidade: edital.TituloMinimoElegibilidade.String(),
		ExigeEmpresa:              edital.ExigeEmpresa,
		PorteEmpresa:              porteJSON,
		EnquadramentoEmpresa:      enquadramentoJSON,
		DocumentosObrigatorios:    documentosJSON,
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
	return r.sqlcParaDominio(result)
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
		e, err := r.sqlcParaDominioListar(result)
		if err != nil {
			return nil, err
		}
		editais = append(editais, e)
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

	relatoriosJSON, err := json.Marshal(edital.RelatoriosExigidos)
	if err != nil {
		return err
	}
	porteJSON, err := json.Marshal(edital.PorteEmpresa)
	if err != nil {
		return err
	}
	enquadramentoJSON, err := json.Marshal(edital.EnquadramentoEmpresa)
	if err != nil {
		return err
	}
	documentosJSON, err := json.Marshal(edital.DocumentosObrigatorios)
	if err != nil {
		return err
	}

	params := sqlcpersistencia.AtualizarEditalParams{
		ID:                        edital.ID,
		Nome:                      edital.Nome,
		Descricao:                 edital.Descricao,
		DataInicio:                dataInicio,
		DataFim:                   dataFim,
		Status:                    edital.Status.String(),
		TipoChamada:               edital.TipoChamada.String(),
		NotaDeCorte:               int32(edital.NotaDeCorte),
		ValorGlobal:               edital.ValorGlobal,
		ModeloFormulario:          int32(edital.ModeloFormulario),
		RelatoriosExigidos:        relatoriosJSON,
		TituloMinimoElegibilidade: edital.TituloMinimoElegibilidade.String(),
		ExigeEmpresa:              edital.ExigeEmpresa,
		PorteEmpresa:              porteJSON,
		EnquadramentoEmpresa:      enquadramentoJSON,
		DocumentosObrigatorios:    documentosJSON,
	}
	return r.queries.AtualizarEdital(ctx, params)
}

func (r *RepositorioDeEditalSQLC) Deletar(ctx context.Context, id int64) error {
	return r.queries.DeletarEdital(ctx, id)
}

func (r *RepositorioDeEditalSQLC) sqlcParaDominioListar(s sqlcpersistencia.ListarEditaisRow) (*entidades.Edital, error) {
	return r.paraEdital(s.Status, s.CriadoEm, s.ID, s.Nome, s.Descricao, s.DataInicio, s.DataFim,
		s.TipoChamada, s.NotaDeCorte, s.ValorGlobal, s.ModeloFormulario,
		s.RelatoriosExigidos, s.TituloMinimoElegibilidade, s.ExigeEmpresa,
		s.PorteEmpresa, s.EnquadramentoEmpresa, s.DocumentosObrigatorios)
}

func (r *RepositorioDeEditalSQLC) sqlcParaDominio(s sqlcpersistencia.BuscarEditalPorIDRow) (*entidades.Edital, error) {
	return r.paraEdital(s.Status, s.CriadoEm, s.ID, s.Nome, s.Descricao, s.DataInicio, s.DataFim,
		s.TipoChamada, s.NotaDeCorte, s.ValorGlobal, s.ModeloFormulario,
		s.RelatoriosExigidos, s.TituloMinimoElegibilidade, s.ExigeEmpresa,
		s.PorteEmpresa, s.EnquadramentoEmpresa, s.DocumentosObrigatorios)
}

func (r *RepositorioDeEditalSQLC) paraEdital(
	statusStr string,
	criadoEm pgtype.Timestamptz,
	id int64, nome, descricao string,
	dataInicio, dataFim pgtype.Date,
	tipoChamada string,
	notaDeCorte int32,
	valorGlobal int64,
	modeloFormulario int32,
	relatoriosExigidos []byte,
	tituloMinimoElegibilidade string,
	exigeEmpresa bool,
	porteEmpresa, enquadramentoEmpresa, documentosObrigatorios []byte,
) (*entidades.Edital, error) {
	status, err := objetos_de_valor.NovoStatusEdital(statusStr)
	if err != nil {
		status = objetos_de_valor.StatusEditalAtivo
	}

	var inicio, fim time.Time
	if dataInicio.Valid {
		inicio = dataInicio.Time
	}
	if dataFim.Valid {
		fim = dataFim.Time
	}

	modelo, errModelo := objetos_de_valor.NovoModeloFormulario(int(modeloFormulario))
	if errModelo != nil {
		modelo = objetos_de_valor.ModeloFormularioNaoDefinido
	}

	tipoChamadaVO, errTipo := objetos_de_valor.NovoTipoChamada(tipoChamada)
	if errTipo != nil {
		tipoChamadaVO = objetos_de_valor.TipoChamadaAPQ
	}

	tituloMinimo, errTitulo := objetos_de_valor.NovoTituloMinimoElegibilidade(tituloMinimoElegibilidade)
	if errTitulo != nil {
		tituloMinimo = objetos_de_valor.TituloMinimoNaoExigido
	}

	var relatorios, porte, enquadramento, documentos []string
	if err := json.Unmarshal(relatoriosExigidos, &relatorios); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(porteEmpresa, &porte); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(enquadramentoEmpresa, &enquadramento); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(documentosObrigatorios, &documentos); err != nil {
		return nil, err
	}

	if relatorios == nil {
		relatorios = []string{}
	}
	if porte == nil {
		porte = []string{}
	}
	if enquadramento == nil {
		enquadramento = []string{}
	}
	if documentos == nil {
		documentos = []string{}
	}

	return &entidades.Edital{
		ID:                        id,
		Nome:                      nome,
		Descricao:                 descricao,
		DataInicio:                inicio,
		DataFim:                   fim,
		Status:                    status,
		TipoChamada:               tipoChamadaVO,
		NotaDeCorte:               int(notaDeCorte),
		ValorGlobal:               valorGlobal,
		ModeloFormulario:          modelo,
		RelatoriosExigidos:        relatorios,
		TituloMinimoElegibilidade: tituloMinimo,
		ExigeEmpresa:              exigeEmpresa,
		PorteEmpresa:              porte,
		EnquadramentoEmpresa:      enquadramento,
		DocumentosObrigatorios:    documentos,
		CriadoEm:                  criadoEm.Time,
	}, nil
}
