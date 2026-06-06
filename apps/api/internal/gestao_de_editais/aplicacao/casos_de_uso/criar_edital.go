package casos_de_uso

import (
	"context"
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
)

type CriarEdital struct {
	repo repositorios.RepositorioDeEdital
}

func NovoEdital(repo repositorios.RepositorioDeEdital) *CriarEdital {
	return &CriarEdital{repo: repo}
}

func (c *CriarEdital) Executar(ctx context.Context, entrada dto.CriarEditalEntrada) (*dto.EditalSaida, error) {
	dataInicio, err := time.Parse("2006-01-02", entrada.DataInicio)
	if err != nil {
		return nil, fmt.Errorf("data de inicio invalida: %w", err)
	}

	dataFim, err := time.Parse("2006-01-02", entrada.DataFim)
	if err != nil {
		return nil, fmt.Errorf("data de fim invalida: %w", err)
	}

	if entrada.RelatoriosExigidos == nil {
		entrada.RelatoriosExigidos = []string{}
	}
	if entrada.PorteEmpresa == nil {
		entrada.PorteEmpresa = []string{}
	}
	if entrada.EnquadramentoEmpresa == nil {
		entrada.EnquadramentoEmpresa = []string{}
	}
	if entrada.DocumentosObrigatorios == nil {
		entrada.DocumentosObrigatorios = []string{}
	}

	params := entidades.NovoEditalParams{
		Nome:                      entrada.Nome,
		Descricao:                 entrada.Descricao,
		DataInicio:                dataInicio,
		DataFim:                   dataFim,
		TipoChamada:               entrada.TipoChamada,
		ModeloFormulario:          entrada.ModeloFormulario,
		RelatoriosExigidos:        entrada.RelatoriosExigidos,
		TituloMinimoElegibilidade: entrada.TituloMinimoElegibilidade,
		ExigeEmpresa:              entrada.ExigeEmpresa,
		PorteEmpresa:              entrada.PorteEmpresa,
		EnquadramentoEmpresa:      entrada.EnquadramentoEmpresa,
		DocumentosObrigatorios:    entrada.DocumentosObrigatorios,
	}

	edital, err := entidades.NovoEdital(params)
	if err != nil {
		return nil, err
	}

	if err := c.repo.Criar(ctx, edital); err != nil {
		return nil, fmt.Errorf("erro ao criar edital: %w", err)
	}

	return paraEditalSaida(edital), nil
}

func paraEditalSaida(e *entidades.Edital) *dto.EditalSaida {
	return &dto.EditalSaida{
		ID:                        e.ID,
		Nome:                      e.Nome,
		Descricao:                 e.Descricao,
		DataInicio:                e.DataInicio.Format("2006-01-02"),
		DataFim:                   e.DataFim.Format("2006-01-02"),
		Status:                    e.Status.String(),
		TipoChamada:               e.TipoChamada.String(),
		NotaDeCorte:               e.NotaDeCorte,
		ValorGlobal:               e.ValorGlobal,
		ModeloFormulario:          int(e.ModeloFormulario),
		RelatoriosExigidos:        e.RelatoriosExigidos,
		TituloMinimoElegibilidade: e.TituloMinimoElegibilidade.String(),
		ExigeEmpresa:              e.ExigeEmpresa,
		PorteEmpresa:              e.PorteEmpresa,
		EnquadramentoEmpresa:      e.EnquadramentoEmpresa,
		DocumentosObrigatorios:    e.DocumentosObrigatorios,
		CriadoEm:                  e.CriadoEm.Format(time.RFC3339),
	}
}
