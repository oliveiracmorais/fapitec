package persistencia

import (
	"context"
	"strings"
	"sync"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
)

type RepositorioDeEditalMemoria struct {
	mu       sync.RWMutex
	editais  map[int64]*entidades.Edital
	proximoID int64
}

func NovoRepositorioDeEditalMemoria() *RepositorioDeEditalMemoria {
	return &RepositorioDeEditalMemoria{
		editais:  make(map[int64]*entidades.Edital),
		proximoID: 1,
	}
}

func (r *RepositorioDeEditalMemoria) Criar(ctx context.Context, edital *entidades.Edital) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	edital.ID = r.proximoID
	r.proximoID++
	r.editais[edital.ID] = edital
	return nil
}

func (r *RepositorioDeEditalMemoria) BuscarPorID(ctx context.Context, id int64) (*entidades.Edital, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	edital, ok := r.editais[id]
	if !ok {
		return nil, nil
	}
	return edital, nil
}

func (r *RepositorioDeEditalMemoria) Listar(ctx context.Context, filtros repositorios.FiltrosListarEditais) ([]*entidades.Edital, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resultado := make([]*entidades.Edital, 0, len(r.editais))
	for _, e := range r.editais {
		if filtros.Titulo != "" && !strings.Contains(strings.ToLower(e.Nome), strings.ToLower(filtros.Titulo)) {
			continue
		}
		if filtros.Status != "" && string(e.Status) != filtros.Status {
			continue
		}
		if filtros.TipoChamada != "" && e.TipoChamada != filtros.TipoChamada {
			continue
		}
		resultado = append(resultado, e)
	}

	return resultado, nil
}

func (r *RepositorioDeEditalMemoria) Atualizar(ctx context.Context, edital *entidades.Edital) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.editais[edital.ID] = edital
	return nil
}

func (r *RepositorioDeEditalMemoria) Deletar(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.editais, id)
	return nil
}
