package entidades

import (
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/objetos_de_valor"
)

type Edital struct {
	ID             int64
	Nome           string
	Descricao      string
	DataInicio     time.Time
	DataFim        time.Time
	Status         objetos_de_valor.StatusEdital
	TipoChamada    string
	NotaDeCorte    int
	ValorGlobal    int64
	CriadoEm       time.Time
}

type NovoEditalParams struct {
	Nome          string
	Descricao     string
	DataInicio    time.Time
	DataFim       time.Time
	TipoChamada   string
}

func NovoEdital(params NovoEditalParams) (*Edital, error) {
	if params.Nome == "" {
		return nil, fmt.Errorf("nome do edital é obrigatorio")
	}
	if params.Descricao == "" {
		return nil, fmt.Errorf("descricao do edital é obrigatoria")
	}
	if params.DataInicio.IsZero() {
		return nil, fmt.Errorf("data de inicio é obrigatoria")
	}
	if params.DataFim.IsZero() {
		return nil, fmt.Errorf("data de fim é obrigatoria")
	}
	if params.DataInicio.After(params.DataFim) {
		return nil, fmt.Errorf("data de inicio nao pode ser posterior a data de fim")
	}

	status, err := objetos_de_valor.NovoStatusEdital("ativo")
	if err != nil {
		return nil, err
	}

	return &Edital{
		Nome:          params.Nome,
		Descricao:     params.Descricao,
		DataInicio:    params.DataInicio,
		DataFim:       params.DataFim,
		Status:        status,
		TipoChamada:   params.TipoChamada,
		CriadoEm:      time.Now(),
	}, nil
}
