package casos_de_uso

import (
	"context"
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
)

type AtualizarEdital struct {
	repo repositorios.RepositorioDeEdital
}

func NovoAtualizarEdital(repo repositorios.RepositorioDeEdital) *AtualizarEdital {
	return &AtualizarEdital{repo: repo}
}

func (a *AtualizarEdital) Executar(ctx context.Context, id int64, entrada dto.AtualizarEditalEntrada) (*dto.EditalSaida, error) {
	edital, err := a.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar edital: %w", err)
	}
	if edital == nil {
		return nil, fmt.Errorf("edital nao encontrado")
	}

	if entrada.Nome != "" {
		edital.Nome = entrada.Nome
	}
	if entrada.Descricao != "" {
		edital.Descricao = entrada.Descricao
	}
	if entrada.DataInicio != "" {
		dataInicio, err := time.Parse("2006-01-02", entrada.DataInicio)
		if err != nil {
			return nil, fmt.Errorf("data de inicio invalida: %w", err)
		}
		edital.DataInicio = dataInicio
	}
	if entrada.DataFim != "" {
		dataFim, err := time.Parse("2006-01-02", entrada.DataFim)
		if err != nil {
			return nil, fmt.Errorf("data de fim invalida: %w", err)
		}
		edital.DataFim = dataFim
	}
	if entrada.Status != "" {
		status, err := objetos_de_valor.NovoStatusEdital(entrada.Status)
		if err != nil {
			return nil, err
		}
		edital.Status = status
	}
	if entrada.TipoChamada != "" {
		tipoChamada, err := objetos_de_valor.NovoTipoChamada(entrada.TipoChamada)
		if err != nil {
			return nil, err
		}
		edital.TipoChamada = tipoChamada
	}
	if entrada.NotaDeCorte != nil {
		edital.NotaDeCorte = *entrada.NotaDeCorte
	}
	if entrada.ValorGlobal != nil {
		edital.ValorGlobal = *entrada.ValorGlobal
	}
	if entrada.ModeloFormulario != nil {
		modelo, err := objetos_de_valor.NovoModeloFormulario(*entrada.ModeloFormulario)
		if err != nil {
			return nil, err
		}
		edital.ModeloFormulario = modelo
	}
	if entrada.RelatoriosExigidos != nil {
		edital.RelatoriosExigidos = entrada.RelatoriosExigidos
	}
	if entrada.TituloMinimoElegibilidade != nil {
		titulo, err := objetos_de_valor.NovoTituloMinimoElegibilidade(*entrada.TituloMinimoElegibilidade)
		if err != nil {
			return nil, err
		}
		edital.TituloMinimoElegibilidade = titulo
	}
	if entrada.ExigeEmpresa != nil {
		edital.ExigeEmpresa = *entrada.ExigeEmpresa
	}
	if entrada.PorteEmpresa != nil {
		edital.PorteEmpresa = entrada.PorteEmpresa
	}
	if entrada.EnquadramentoEmpresa != nil {
		edital.EnquadramentoEmpresa = entrada.EnquadramentoEmpresa
	}
	if entrada.DocumentosObrigatorios != nil {
		edital.DocumentosObrigatorios = entrada.DocumentosObrigatorios
	}

	if !edital.DataFim.IsZero() && !edital.DataInicio.IsZero() && edital.DataInicio.After(edital.DataFim) {
		return nil, fmt.Errorf("data de inicio nao pode ser posterior a data de fim")
	}

	if err := a.repo.Atualizar(ctx, edital); err != nil {
		return nil, fmt.Errorf("erro ao atualizar edital: %w", err)
	}

	return paraEditalSaida(edital), nil
}
