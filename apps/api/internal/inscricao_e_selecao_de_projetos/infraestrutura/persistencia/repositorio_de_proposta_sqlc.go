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

type RepositorioDePropostaSQLC struct {
	queries *sqlcpersistencia.Queries
}

func NovoRepositorioDePropostaSQLC(queries *sqlcpersistencia.Queries) *RepositorioDePropostaSQLC {
	return &RepositorioDePropostaSQLC{queries: queries}
}

func (r *RepositorioDePropostaSQLC) Criar(ctx context.Context, proposta *entidades.Proposta) error {
	params := sqlcpersistencia.InserirPropostaParams{
		EditalID:                  proposta.EditalID,
		ProponenteID:              proposta.ProponenteID,
		Protocolo:                 proposta.Protocolo.String(),
		Status:                    proposta.Status.String(),
		Versao:                    int32(proposta.Versao),
		ProponenteNome:            proposta.DadosProponente.Nome,
		ProponenteCpf:             proposta.DadosProponente.CPF,
		ProponenteRg:              proposta.DadosProponente.RG,
		ProponenteGenero:          proposta.DadosProponente.Genero,
		ProponenteEtnia:           proposta.DadosProponente.Etnia,
		ProponenteDataNascimento:  proposta.DadosProponente.DataNascimento,
		ProponenteEndereco:        proposta.DadosProponente.Endereco,
		ProponenteTelefone:        proposta.DadosProponente.Telefone,
		ProponenteEmail:           proposta.DadosProponente.Email,
		AcademicoMaiorTitulacao:   proposta.DadosAcademicos.MaiorTitulacao,
		AcademicoCurso:            proposta.DadosAcademicos.Curso,
		AcademicoInstituicao:      proposta.DadosAcademicos.Instituicao,
		AcademicoAnoConclusao:     int32(proposta.DadosAcademicos.AnoConclusao),
		AcademicoAreaConhecimento: proposta.DadosAcademicos.AreaConhecimento,
		EmpresaVinculada:          proposta.EmpresaVinculada,
		ValorTotalSolicitado:      proposta.ValorTotalSolicitado,
	}

	result, err := r.queries.InserirProposta(ctx, params)
	if err != nil {
		return err
	}
	proposta.ID = result.ID
	proposta.CriadoEm = result.CriadoEm.Time
	return nil
}

func (r *RepositorioDePropostaSQLC) BuscarPorID(ctx context.Context, id int64) (*entidades.Proposta, error) {
	result, err := r.queries.BuscarPropostaPorID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.sqlcParaDominio(result), nil
}

func (r *RepositorioDePropostaSQLC) BuscarPorProtocolo(ctx context.Context, protocolo string) (*entidades.Proposta, error) {
	result, err := r.queries.BuscarPropostaPorProtocolo(ctx, protocolo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.sqlcParaDominio(result), nil
}

func (r *RepositorioDePropostaSQLC) Listar(ctx context.Context, filtros repositorios.FiltrosListarPropostas) ([]*entidades.Proposta, error) {
	params := sqlcpersistencia.ListarPropostasParams{
		EditalID:     filtros.EditalID,
		ProponenteID: filtros.ProponenteID,
		Status:       filtros.Status,
	}
	results, err := r.queries.ListarPropostas(ctx, params)
	if err != nil {
		return nil, err
	}
	propostas := make([]*entidades.Proposta, 0, len(results))
	for _, result := range results {
		propostas = append(propostas, r.sqlcParaDominio(result))
	}
	return propostas, nil
}

func (r *RepositorioDePropostaSQLC) Atualizar(ctx context.Context, proposta *entidades.Proposta) error {
	existente, err := r.queries.BuscarPropostaPorID(ctx, proposta.ID)
	if err != nil {
		return err
	}

	versaoAntiga := int32(existente.Versao)
	novaVersao := int(existente.Versao) + 1

	dataSubmissao := pgtype.Timestamptz{}
	if !existente.DataSubmissao.Time.IsZero() {
		dataSubmissao = existente.DataSubmissao
	}
	if !proposta.DataSubmissao.IsZero() {
		dataSubmissao.Scan(proposta.DataSubmissao)
	}

	versaoParams := sqlcpersistencia.InserirVersaoPropostaParams{
		PropostaID:                proposta.ID,
		Versao:                    versaoAntiga,
		Status:                    existente.Status,
		ProponenteNome:            existente.ProponenteNome,
		ProponenteCpf:             existente.ProponenteCpf,
		ProponenteRg:              existente.ProponenteRg,
		ProponenteGenero:          existente.ProponenteGenero,
		ProponenteEtnia:           existente.ProponenteEtnia,
		ProponenteDataNascimento:  existente.ProponenteDataNascimento,
		ProponenteEndereco:        existente.ProponenteEndereco,
		ProponenteTelefone:        existente.ProponenteTelefone,
		ProponenteEmail:           existente.ProponenteEmail,
		AcademicoMaiorTitulacao:   existente.AcademicoMaiorTitulacao,
		AcademicoCurso:            existente.AcademicoCurso,
		AcademicoInstituicao:      existente.AcademicoInstituicao,
		AcademicoAnoConclusao:     existente.AcademicoAnoConclusao,
		AcademicoAreaConhecimento: existente.AcademicoAreaConhecimento,
		EmpresaVinculada:          existente.EmpresaVinculada,
		ValorTotalSolicitado:      existente.ValorTotalSolicitado,
		Protocolo:                 existente.Protocolo,
		DataSubmissao:             dataSubmissao,
	}
	if err := r.queries.InserirVersaoProposta(ctx, versaoParams); err != nil {
		return err
	}

	dataSubmissaoAtual := pgtype.Timestamptz{}
	if !proposta.DataSubmissao.IsZero() {
		dataSubmissaoAtual.Scan(proposta.DataSubmissao)
	}

	params := sqlcpersistencia.AtualizarPropostaParams{
		ID:                        proposta.ID,
		Versao:                    int32(novaVersao),
		Status:                    proposta.Status.String(),
		ProponenteNome:            proposta.DadosProponente.Nome,
		ProponenteCpf:             proposta.DadosProponente.CPF,
		ProponenteRg:              proposta.DadosProponente.RG,
		ProponenteGenero:          proposta.DadosProponente.Genero,
		ProponenteEtnia:           proposta.DadosProponente.Etnia,
		ProponenteDataNascimento:  proposta.DadosProponente.DataNascimento,
		ProponenteEndereco:        proposta.DadosProponente.Endereco,
		ProponenteTelefone:        proposta.DadosProponente.Telefone,
		ProponenteEmail:           proposta.DadosProponente.Email,
		AcademicoMaiorTitulacao:   proposta.DadosAcademicos.MaiorTitulacao,
		AcademicoCurso:            proposta.DadosAcademicos.Curso,
		AcademicoInstituicao:      proposta.DadosAcademicos.Instituicao,
		AcademicoAnoConclusao:     int32(proposta.DadosAcademicos.AnoConclusao),
		AcademicoAreaConhecimento: proposta.DadosAcademicos.AreaConhecimento,
		EmpresaVinculada:          proposta.EmpresaVinculada,
		ValorTotalSolicitado:      proposta.ValorTotalSolicitado,
		Protocolo:                 proposta.Protocolo.String(),
		DataSubmissao:             dataSubmissaoAtual,
	}
	if err := r.queries.AtualizarProposta(ctx, params); err != nil {
		return err
	}

	proposta.Versao = novaVersao
	return nil
}

func (r *RepositorioDePropostaSQLC) Deletar(ctx context.Context, id int64) error {
	return r.queries.DeletarProposta(ctx, id)
}

func (r *RepositorioDePropostaSQLC) ContarPorEdital(ctx context.Context, editalID int64) (int64, error) {
	return r.queries.ContarPropostasPorEdital(ctx, editalID)
}

func (r *RepositorioDePropostaSQLC) sqlcParaDominio(s sqlcpersistencia.Proposta) *entidades.Proposta {
	status, _ := objetos_de_valor.NovoStatusProposta(s.Status)
	protocolo, _ := objetos_de_valor.ProtocoloExistente(s.Protocolo)

	dataSubmissao := s.DataSubmissao.Time
	dataAtualizacao := s.DataAtualizacao.Time

	return &entidades.Proposta{
		ID:           s.ID,
		EditalID:     s.EditalID,
		ProponenteID: s.ProponenteID,
		Protocolo:    protocolo,
		Status:       status,
		Versao:       int(s.Versao),
		DadosProponente: entidades.ProponenteInfo{
			Nome:          s.ProponenteNome,
			CPF:           s.ProponenteCpf,
			RG:            s.ProponenteRg,
			Genero:        s.ProponenteGenero,
			Etnia:         s.ProponenteEtnia,
			DataNascimento: s.ProponenteDataNascimento,
			Endereco:      s.ProponenteEndereco,
			Telefone:      s.ProponenteTelefone,
			Email:         s.ProponenteEmail,
		},
		DadosAcademicos: entidades.DadosAcademicos{
			MaiorTitulacao:   s.AcademicoMaiorTitulacao,
			Curso:            s.AcademicoCurso,
			Instituicao:      s.AcademicoInstituicao,
			AnoConclusao:     int(s.AcademicoAnoConclusao),
			AreaConhecimento: s.AcademicoAreaConhecimento,
		},
		EmpresaVinculada:     s.EmpresaVinculada,
		ItensOrcamentarios:   []entidades.ItemOrcamentario{},
		Documentos:           []entidades.DocumentoProposta{},
		ValorTotalSolicitado: s.ValorTotalSolicitado,
		Pareceres:            []entidades.Parecer{},
		DataSubmissao:        dataSubmissao,
		DataAtualizacao:      dataAtualizacao,
		CriadoEm:             s.CriadoEm.Time,
	}
}
