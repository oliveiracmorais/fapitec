package entidades

import (
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
)

type Proposta struct {
	ID                   int64
	EditalID             int64
	ProponenteID         int64
	Protocolo            objetos_de_valor.Protocolo
	Status               objetos_de_valor.StatusProposta
	Versao               int
	DadosProponente      ProponenteInfo
	DadosAcademicos      DadosAcademicos
	EmpresaVinculada     string
	ItensOrcamentarios   []ItemOrcamentario
	Documentos           []DocumentoProposta
	ValorTotalSolicitado int64
	Pareceres            []Parecer
	DataSubmissao        time.Time
	DataAtualizacao      time.Time
	CriadoEm             time.Time
}

type NovaPropostaParams struct {
	EditalID         int64
	ProponenteID     int64
	DadosProponente  ProponenteInfo
	DadosAcademicos  DadosAcademicos
	EmpresaVinculada string
	ItensOrcamentarios []ItemOrcamentario
	Documentos       []DocumentoProposta
}

func NovaProposta(params NovaPropostaParams) (*Proposta, error) {
	if params.EditalID == 0 {
		return nil, fmt.Errorf("edital é obrigatorio")
	}
	if params.ProponenteID == 0 {
		return nil, fmt.Errorf("proponente é obrigatorio")
	}
	if params.DadosProponente.Nome == "" {
		return nil, fmt.Errorf("nome do proponente é obrigatorio")
	}
	if params.DadosProponente.CPF == "" {
		return nil, fmt.Errorf("CPF do proponente é obrigatorio")
	}

	status, err := objetos_de_valor.NovoStatusProposta("rascunho")
	if err != nil {
		return nil, err
	}

	if params.ItensOrcamentarios == nil {
		params.ItensOrcamentarios = []ItemOrcamentario{}
	}
	if params.Documentos == nil {
		params.Documentos = []DocumentoProposta{}
	}

	valorTotal := int64(0)
	for _, item := range params.ItensOrcamentarios {
		valorTotal += item.ValorTotal
	}

	agora := time.Now()

	return &Proposta{
		EditalID:             params.EditalID,
		ProponenteID:         params.ProponenteID,
		Status:               status,
		Versao:               1,
		DadosProponente:      params.DadosProponente,
		DadosAcademicos:      params.DadosAcademicos,
		EmpresaVinculada:     params.EmpresaVinculada,
		ItensOrcamentarios:   params.ItensOrcamentarios,
		Documentos:           params.Documentos,
		ValorTotalSolicitado: valorTotal,
		Pareceres:            []Parecer{},
		CriadoEm:             agora,
		DataAtualizacao:      agora,
	}, nil
}

func (p *Proposta) Submeter() error {
	if p.Status != objetos_de_valor.StatusPropostaRascunho {
		return fmt.Errorf("apenas propostas em rascunho podem ser submetidas")
	}
	status, err := objetos_de_valor.NovoStatusProposta("submetida")
	if err != nil {
		return err
	}
	p.Status = status
	p.DataSubmissao = time.Now()
	p.DataAtualizacao = time.Now()
	return nil
}

func (p *Proposta) CalcularValorTotal() {
	total := int64(0)
	for _, item := range p.ItensOrcamentarios {
		total += item.ValorTotal
	}
	p.ValorTotalSolicitado = total
}

func (p *Proposta) AdicionarParecer(parecer *Parecer) error {
	if p.Status != objetos_de_valor.StatusPropostaSubmetida &&
		p.Status != objetos_de_valor.StatusPropostaEmAvaliacao {
		return fmt.Errorf("apenas propostas submetidas ou em avaliacao podem receber pareceres")
	}

	p.Pareceres = append(p.Pareceres, *parecer)
	p.DataAtualizacao = time.Now()

	if p.Status == objetos_de_valor.StatusPropostaSubmetida {
		status, err := objetos_de_valor.NovoStatusProposta("em_avaliacao")
		if err != nil {
			return err
		}
		p.Status = status
	}

	return nil
}

func (p *Proposta) FinalizarAvaliacao(notaDeCorte int) error {
	if p.Status != objetos_de_valor.StatusPropostaEmAvaliacao {
		return fmt.Errorf("apenas propostas em avaliacao podem ser finalizadas")
	}

	notaFinal := p.CalcularNotaFinal()

	var novoStatus objetos_de_valor.StatusProposta
	var err error
	if notaFinal >= notaDeCorte {
		novoStatus, err = objetos_de_valor.NovoStatusProposta("aprovada")
	} else {
		novoStatus, err = objetos_de_valor.NovoStatusProposta("reprovada")
	}
	if err != nil {
		return err
	}

	p.Status = novoStatus
	p.DataAtualizacao = time.Now()
	return nil
}

func (p *Proposta) CalcularNotaFinal() int {
	if len(p.Pareceres) == 0 {
		return 0
	}
	soma := 0
	for _, parecer := range p.Pareceres {
		soma += parecer.Nota
	}
	return soma / len(p.Pareceres)
}
