package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type DeletarProposta struct {
	repo repositorios.RepositorioDeProposta
}

func NovoDeletarProposta(repo repositorios.RepositorioDeProposta) *DeletarProposta {
	return &DeletarProposta{repo: repo}
}

func (uc *DeletarProposta) Executar(ctx context.Context, id int64) error {
	proposta, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao buscar proposta: %w", err)
	}
	if proposta == nil {
		return fmt.Errorf("proposta nao encontrada")
	}

	if proposta.Status != objetos_de_valor.StatusPropostaRascunho {
		return fmt.Errorf("apenas propostas em rascunho podem ser deletadas")
	}

	return uc.repo.Deletar(ctx, id)
}
