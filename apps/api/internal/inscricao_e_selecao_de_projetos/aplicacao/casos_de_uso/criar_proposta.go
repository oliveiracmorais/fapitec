package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type CriarProposta struct {
	repo repositorios.RepositorioDeProposta
}

func NovoCriarProposta(repo repositorios.RepositorioDeProposta) *CriarProposta {
	return &CriarProposta{repo: repo}
}

func (uc *CriarProposta) Executar(ctx context.Context, entrada dto.CriarPropostaEntrada) (*dto.PropostaSaida, error) {
	itens := make([]entidades.ItemOrcamentario, len(entrada.ItensOrcamentarios))
	for i, item := range entrada.ItensOrcamentarios {
		tipo, err := objetos_de_valor.NovoTipoItemOrcamentario(item.Tipo)
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		itens[i] = entidades.ItemOrcamentario{
			Descricao:     item.Descricao,
			Tipo:          tipo,
			Quantidade:    item.Quantidade,
			ValorUnitario: item.ValorUnitario,
			ValorTotal:    item.ValorTotal,
		}
	}

	documentos := make([]entidades.DocumentoProposta, len(entrada.Documentos))
	for i, doc := range entrada.Documentos {
		tipo, err := objetos_de_valor.NovoTipoDocumento(doc.Tipo)
		if err != nil {
			return nil, fmt.Errorf("documento %d: %w", i, err)
		}
		documentos[i] = entidades.DocumentoProposta{
			Tipo:        tipo,
			NomeArquivo: doc.NomeArquivo,
			Caminho:     doc.Caminho,
		}
	}

	params := entidades.NovaPropostaParams{
		EditalID:         entrada.EditalID,
		ProponenteID:     entrada.ProponenteID,
		DadosProponente:  toProponenteInfo(entrada.DadosProponente),
		DadosAcademicos:  toDadosAcademicos(entrada.DadosAcademicos),
		EmpresaVinculada: entrada.EmpresaVinculada,
		ItensOrcamentarios: itens,
		Documentos:       documentos,
	}

	proposta, err := entidades.NovaProposta(params)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Criar(ctx, proposta); err != nil {
		return nil, err
	}

	return paraPropostaSaida(proposta), nil
}

func toProponenteInfo(d dto.ProponenteInfoDTO) entidades.ProponenteInfo {
	return entidades.ProponenteInfo{
		Nome:          d.Nome,
		CPF:           d.CPF,
		RG:            d.RG,
		Genero:        d.Genero,
		Etnia:         d.Etnia,
		DataNascimento: d.DataNascimento,
		Endereco:      d.Endereco,
		Telefone:      d.Telefone,
		Email:         d.Email,
	}
}

func toDadosAcademicos(d dto.DadosAcademicosDTO) entidades.DadosAcademicos {
	return entidades.DadosAcademicos{
		MaiorTitulacao:   d.MaiorTitulacao,
		Curso:            d.Curso,
		Instituicao:      d.Instituicao,
		AnoConclusao:     d.AnoConclusao,
		AreaConhecimento: d.AreaConhecimento,
	}
}

func paraPropostaSaida(p *entidades.Proposta) *dto.PropostaSaida {
	itens := make([]dto.ItemOrcamentarioDTO, len(p.ItensOrcamentarios))
	for i, item := range p.ItensOrcamentarios {
		itens[i] = dto.ItemOrcamentarioDTO{
			Descricao:     item.Descricao,
			Tipo:          item.Tipo.String(),
			Quantidade:    item.Quantidade,
			ValorUnitario: item.ValorUnitario,
			ValorTotal:    item.ValorTotal,
		}
	}

	documentos := make([]dto.DocumentoPropostaDTO, len(p.Documentos))
	for i, doc := range p.Documentos {
		documentos[i] = dto.DocumentoPropostaDTO{
			Tipo:        doc.Tipo.String(),
			NomeArquivo: doc.NomeArquivo,
			Caminho:     doc.Caminho,
		}
	}

	pareceres := make([]dto.ParecerDTO, len(p.Pareceres))
	for i, par := range p.Pareceres {
		pareceres[i] = dto.ParecerDTO{
			Etapa:        par.Etapa,
			AvaliadorID:  par.AvaliadorID,
			Nota:         par.Nota,
			ParecerTexto: par.ParecerTexto,
			Data:         par.Data.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	dataSubmissao := ""
	if !p.DataSubmissao.IsZero() {
		dataSubmissao = p.DataSubmissao.Format("2006-01-02T15:04:05Z07:00")
	}

	return &dto.PropostaSaida{
		ID:                 p.ID,
		EditalID:           p.EditalID,
		ProponenteID:       p.ProponenteID,
		Protocolo:          p.Protocolo.String(),
		Status:             p.Status.String(),
		Versao:             p.Versao,
		DadosProponente:    toProponenteInfoDTO(p.DadosProponente),
		DadosAcademicos:    toDadosAcademicosDTO(p.DadosAcademicos),
		EmpresaVinculada:   p.EmpresaVinculada,
		ItensOrcamentarios: itens,
		Documentos:         documentos,
		ValorTotal:         p.ValorTotalSolicitado,
		Pareceres:          pareceres,
		DataSubmissao:      dataSubmissao,
		DataAtualizacao:    p.DataAtualizacao.Format("2006-01-02T15:04:05Z07:00"),
		CriadoEm:           p.CriadoEm.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toProponenteInfoDTO(p entidades.ProponenteInfo) dto.ProponenteInfoDTO {
	return dto.ProponenteInfoDTO{
		Nome:          p.Nome,
		CPF:           p.CPF,
		RG:            p.RG,
		Genero:        p.Genero,
		Etnia:         p.Etnia,
		DataNascimento: p.DataNascimento,
		Endereco:      p.Endereco,
		Telefone:      p.Telefone,
		Email:         p.Email,
	}
}

func toDadosAcademicosDTO(d entidades.DadosAcademicos) dto.DadosAcademicosDTO {
	return dto.DadosAcademicosDTO{
		MaiorTitulacao:   d.MaiorTitulacao,
		Curso:            d.Curso,
		Instituicao:      d.Instituicao,
		AnoConclusao:     d.AnoConclusao,
		AreaConhecimento: d.AreaConhecimento,
	}
}
