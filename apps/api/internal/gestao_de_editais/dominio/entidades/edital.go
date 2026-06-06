package entidades

import (
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/objetos_de_valor"
)

type Edital struct {
	ID                        int64
	Nome                      string
	Descricao                 string
	DataInicio                time.Time
	DataFim                   time.Time
	Status                    objetos_de_valor.StatusEdital
	TipoChamada               objetos_de_valor.TipoChamada
	NotaDeCorte               int
	ValorGlobal               int64
	ModeloFormulario          objetos_de_valor.ModeloFormulario
	RelatoriosExigidos        []string
	TituloMinimoElegibilidade objetos_de_valor.TituloMinimoElegibilidade
	ExigeEmpresa              bool
	PorteEmpresa              []string
	EnquadramentoEmpresa      []string
	DocumentosObrigatorios    []string
	CriadoEm                  time.Time
}

type NovoEditalParams struct {
	Nome                      string
	Descricao                 string
	DataInicio                time.Time
	DataFim                   time.Time
	TipoChamada               string
	ModeloFormulario          int
	RelatoriosExigidos        []string
	TituloMinimoElegibilidade string
	ExigeEmpresa              bool
	PorteEmpresa              []string
	EnquadramentoEmpresa      []string
	DocumentosObrigatorios    []string
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

	tipoChamada, err := objetos_de_valor.NovoTipoChamada(params.TipoChamada)
	if err != nil {
		return nil, err
	}

	modeloFormulario, err := objetos_de_valor.NovoModeloFormulario(params.ModeloFormulario)
	if err != nil {
		return nil, err
	}

	tituloMinimo, err := objetos_de_valor.NovoTituloMinimoElegibilidade(params.TituloMinimoElegibilidade)
	if err != nil {
		return nil, err
	}

	if params.RelatoriosExigidos == nil {
		params.RelatoriosExigidos = []string{}
	}
	if params.PorteEmpresa == nil {
		params.PorteEmpresa = []string{}
	}
	if params.EnquadramentoEmpresa == nil {
		params.EnquadramentoEmpresa = []string{}
	}
	if params.DocumentosObrigatorios == nil {
		params.DocumentosObrigatorios = []string{}
	}

	for _, p := range params.PorteEmpresa {
		if !objetos_de_valor.PorteEmpresaValido(p) {
			return nil, fmt.Errorf("porte de empresa invalido: %s", p)
		}
	}
	for _, e := range params.EnquadramentoEmpresa {
		if !objetos_de_valor.EnquadramentoEmpresaValido(e) {
			return nil, fmt.Errorf("enquadramento de empresa invalido: %s", e)
		}
	}
	for _, d := range params.DocumentosObrigatorios {
		if !objetos_de_valor.DocumentoObrigatorioValido(d) {
			return nil, fmt.Errorf("documento obrigatorio invalido: %s", d)
		}
	}

	return &Edital{
		Nome:                      params.Nome,
		Descricao:                 params.Descricao,
		DataInicio:                params.DataInicio,
		DataFim:                   params.DataFim,
		Status:                    status,
		TipoChamada:               tipoChamada,
		ModeloFormulario:          modeloFormulario,
		RelatoriosExigidos:        params.RelatoriosExigidos,
		TituloMinimoElegibilidade: tituloMinimo,
		ExigeEmpresa:              params.ExigeEmpresa,
		PorteEmpresa:              params.PorteEmpresa,
		EnquadramentoEmpresa:      params.EnquadramentoEmpresa,
		DocumentosObrigatorios:    params.DocumentosObrigatorios,
		CriadoEm:                  time.Now(),
	}, nil
}
