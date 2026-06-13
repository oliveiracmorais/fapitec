package casos_de_uso

import (
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
)

func paraAvaliadorSaida(a *entidades.Avaliador, totalPropostas, atribuicoesAtivas int64) *dto.AvaliadorSaida {
	return &dto.AvaliadorSaida{
		ID:                a.ID,
		UsuarioID:         a.UsuarioID,
		Nome:              a.Nome,
		CPF:               a.CPF,
		Email:             a.Email,
		TitulacaoMaxima:   a.TitulacaoMaxima,
		AreaConhecimento:  a.AreaConhecimento,
		Instituicao:       a.Instituicao,
		CurriculoResumido: a.CurriculoResumido,
		Estado:            a.Estado.String(),
		DataCadastro:      a.DataCadastro.Format("2006-01-02T15:04:05Z07:00"),
		DataAtualizacao:   a.DataAtualizacao.Format("2006-01-02T15:04:05Z07:00"),
		TotalPropostas:    totalPropostas,
		AtribuicoesAtivas: atribuicoesAtivas,
	}
}

func paraAtribuicaoSaida(a *entidades.AtribuicaoEdital) *dto.AtribuicaoSaida {
	return &dto.AtribuicaoSaida{
		ID:               a.ID,
		AvaliadorID:      a.AvaliadorID,
		EditalID:         a.EditalID,
		DataInicio:       a.DataInicio.Format("2006-01-02T15:04:05Z07:00"),
		DataFim:          a.DataFim.Format("2006-01-02T15:04:05Z07:00"),
		StatusConvite:    a.StatusConvite.String(),
		HashAnonimizacao: a.HashAnonimizacao.String(),
		CriadoEm:         a.CriadoEm.Format("2006-01-02T15:04:05Z07:00"),
	}
}
