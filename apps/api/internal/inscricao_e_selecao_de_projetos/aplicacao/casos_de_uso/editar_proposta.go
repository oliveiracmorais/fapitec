package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type EditarProposta struct {
	repo repositorios.RepositorioDeProposta
}

func NovoEditarProposta(repo repositorios.RepositorioDeProposta) *EditarProposta {
	return &EditarProposta{repo: repo}
}

func (uc *EditarProposta) Executar(ctx context.Context, id int64, entrada dto.AtualizarPropostaEntrada) (*dto.PropostaSaida, error) {
	proposta, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar proposta: %w", err)
	}
	if proposta == nil {
		return nil, fmt.Errorf("proposta nao encontrada")
	}

	if proposta.Status != objetos_de_valor.StatusPropostaRascunho {
		return nil, fmt.Errorf("apenas propostas em rascunho podem ser editadas")
	}

	if entrada.DadosProponente != nil {
		proposta.DadosProponente = toProponenteInfo(*entrada.DadosProponente)
	}
	if entrada.DadosAcademicos != nil {
		proposta.DadosAcademicos = toDadosAcademicos(*entrada.DadosAcademicos)
	}
	if entrada.EmpresaVinculada != nil {
		proposta.EmpresaVinculada = *entrada.EmpresaVinculada
	}
	if entrada.ItensOrcamentarios != nil {
		itens := make([]entidades.ItemOrcamentario, len(*entrada.ItensOrcamentarios))
		for i, item := range *entrada.ItensOrcamentarios {
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
		proposta.ItensOrcamentarios = itens
	}
	if entrada.Documentos != nil {
		documentos := make([]entidades.DocumentoProposta, len(*entrada.Documentos))
		for i, doc := range *entrada.Documentos {
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
		proposta.Documentos = documentos
	}

	proposta.CalcularValorTotal()

	if err := uc.repo.Atualizar(ctx, proposta); err != nil {
		return nil, fmt.Errorf("erro ao atualizar proposta: %w", err)
	}

	return paraPropostaSaida(proposta), nil
}
