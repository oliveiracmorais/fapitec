package entidades

import (
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
)

type AtribuicaoEdital struct {
	ID               int64
	AvaliadorID      int64
	EditalID         int64
	DataInicio       time.Time
	DataFim          time.Time
	StatusConvite    objetos_de_valor.StatusConvite
	HashAnonimizacao objetos_de_valor.HashAnonimizacao
	CriadoEm         time.Time
}

type NovaAtribuicaoParams struct {
	AvaliadorID int64
	EditalID    int64
	DataInicio  time.Time
	DataFim     time.Time
}

func NovaAtribuicao(params NovaAtribuicaoParams) (*AtribuicaoEdital, error) {
	if params.AvaliadorID == 0 {
		return nil, fmt.Errorf("avaliador e obrigatorio")
	}
	if params.EditalID == 0 {
		return nil, fmt.Errorf("edital e obrigatorio")
	}
	if params.DataFim.Before(params.DataInicio) || params.DataFim.Equal(params.DataInicio) {
		return nil, fmt.Errorf("data fim deve ser posterior a data inicio")
	}

	status, err := objetos_de_valor.NovoStatusConvite("pendente")
	if err != nil {
		return nil, err
	}

	hash := objetos_de_valor.NovoHashAnonimizacao(params.AvaliadorID, params.EditalID)

	return &AtribuicaoEdital{
		AvaliadorID:      params.AvaliadorID,
		EditalID:         params.EditalID,
		DataInicio:       params.DataInicio,
		DataFim:          params.DataFim,
		StatusConvite:    status,
		HashAnonimizacao: hash,
		CriadoEm:         time.Now(),
	}, nil
}

func (a *AtribuicaoEdital) AceitarConvite() error {
	if a.StatusConvite != objetos_de_valor.StatusConvitePendente {
		return fmt.Errorf("convite nao esta pendente")
	}
	a.StatusConvite = objetos_de_valor.StatusConviteAceito
	return nil
}

func (a *AtribuicaoEdital) RecusarConvite() error {
	if a.StatusConvite != objetos_de_valor.StatusConvitePendente {
		return fmt.Errorf("convite nao esta pendente")
	}
	a.StatusConvite = objetos_de_valor.StatusConviteRecusado
	return nil
}

func (a *AtribuicaoEdital) ConviteFoiAceito() bool {
	return a.StatusConvite == objetos_de_valor.StatusConviteAceito
}
