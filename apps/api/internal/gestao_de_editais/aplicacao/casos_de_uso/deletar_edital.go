package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
)

type DeletarEdital struct {
	repo repositorios.RepositorioDeEdital
}

func NovoDeletarEdital(repo repositorios.RepositorioDeEdital) *DeletarEdital {
	return &DeletarEdital{repo: repo}
}

func (d *DeletarEdital) Executar(ctx context.Context, id int64) error {
	edital, err := d.repo.BuscarPorID(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao buscar edital: %w", err)
	}
	if edital == nil {
		return fmt.Errorf("edital nao encontrado")
	}

	if err := d.repo.Deletar(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar edital: %w", err)
	}

	return nil
}
