package entidades

import (
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
)

type Avaliador struct {
	ID              int64
	UsuarioID       int64
	Nome            string
	CPF             string
	Email           string
	TitulacaoMaxima string
	AreaConhecimento string
	Instituicao     string
	CurriculoResumido string
	Estado          objetos_de_valor.EstadoAvaliador
	DataCadastro    time.Time
	DataAtualizacao time.Time
}

type NovoAvaliadorParams struct {
	UsuarioID         int64
	Nome              string
	CPF               string
	Email             string
	TitulacaoMaxima   string
	AreaConhecimento  string
	Instituicao       string
	CurriculoResumido string
}

func NovoAvaliador(params NovoAvaliadorParams) (*Avaliador, error) {
	if params.Nome == "" {
		return nil, fmt.Errorf("nome do avaliador e obrigatorio")
	}
	if params.CPF == "" {
		return nil, fmt.Errorf("CPF do avaliador e obrigatorio")
	}
	if params.UsuarioID == 0 {
		return nil, fmt.Errorf("usuario do avaliador e obrigatorio")
	}

	estado, err := objetos_de_valor.NovoEstadoAvaliador("ativo")
	if err != nil {
		return nil, err
	}

	agora := time.Now()
	return &Avaliador{
		UsuarioID:         params.UsuarioID,
		Nome:              params.Nome,
		CPF:               params.CPF,
		Email:             params.Email,
		TitulacaoMaxima:   params.TitulacaoMaxima,
		AreaConhecimento:  params.AreaConhecimento,
		Instituicao:       params.Instituicao,
		CurriculoResumido: params.CurriculoResumido,
		Estado:            estado,
		DataCadastro:      agora,
		DataAtualizacao:   agora,
	}, nil
}

func (a *Avaliador) AtualizarDados(params NovoAvaliadorParams) {
	a.Nome = params.Nome
	a.CPF = params.CPF
	a.Email = params.Email
	a.TitulacaoMaxima = params.TitulacaoMaxima
	a.AreaConhecimento = params.AreaConhecimento
	a.Instituicao = params.Instituicao
	a.CurriculoResumido = params.CurriculoResumido
	a.DataAtualizacao = time.Now()
}

func (a *Avaliador) Ativar() {
	a.Estado = objetos_de_valor.EstadoAvaliadorAtivo
	a.DataAtualizacao = time.Now()
}

func (a *Avaliador) Inativar() {
	a.Estado = objetos_de_valor.EstadoAvaliadorInativo
	a.DataAtualizacao = time.Now()
}

func (a *Avaliador) EstaInativo() bool {
	return a.Estado == objetos_de_valor.EstadoAvaliadorInativo
}
