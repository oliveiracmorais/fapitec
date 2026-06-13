package persistencia

import (
	"context"
	"sync"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type RepositorioDeAtribuicaoMemoria struct {
	mu        sync.RWMutex
	dados     map[int64]*entidades.AtribuicaoEdital
	hashIdx   map[string]int64
	proximoID int64
}

func NovoRepositorioDeAtribuicaoMemoria() *RepositorioDeAtribuicaoMemoria {
	return &RepositorioDeAtribuicaoMemoria{
		dados:     make(map[int64]*entidades.AtribuicaoEdital),
		hashIdx:   make(map[string]int64),
		proximoID: 1,
	}
}

func (r *RepositorioDeAtribuicaoMemoria) Criar(_ context.Context, atribuicao *entidades.AtribuicaoEdital) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.proximoID
	r.proximoID++
	atribuicao.ID = id
	c := *atribuicao
	r.dados[id] = &c
	r.hashIdx[atribuicao.HashAnonimizacao.String()] = id
	return nil
}

func (r *RepositorioDeAtribuicaoMemoria) BuscarPorID(_ context.Context, id int64) (*entidades.AtribuicaoEdital, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, ok := r.dados[id]
	if !ok {
		return nil, nil
	}
	c := *a
	return &c, nil
}

func (r *RepositorioDeAtribuicaoMemoria) ListarPorAvaliador(_ context.Context, avaliadorID int64) ([]*entidades.AtribuicaoEdital, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resultado []*entidades.AtribuicaoEdital
	for _, a := range r.dados {
		if a.AvaliadorID == avaliadorID {
			c := *a
			resultado = append(resultado, &c)
		}
	}
	return resultado, nil
}

func (r *RepositorioDeAtribuicaoMemoria) ListarPorEdital(_ context.Context, editalID int64) ([]*entidades.AtribuicaoEdital, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resultado []*entidades.AtribuicaoEdital
	for _, a := range r.dados {
		if a.EditalID == editalID {
			c := *a
			resultado = append(resultado, &c)
		}
	}
	return resultado, nil
}

func (r *RepositorioDeAtribuicaoMemoria) Listar(_ context.Context, filtros repositorios.FiltrosListarAtribuicoes) ([]*entidades.AtribuicaoEdital, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resultado []*entidades.AtribuicaoEdital
	for _, a := range r.dados {
		if filtros.AvaliadorID != 0 && a.AvaliadorID != filtros.AvaliadorID {
			continue
		}
		if filtros.EditalID != 0 && a.EditalID != filtros.EditalID {
			continue
		}
		if filtros.Status != "" && a.StatusConvite.String() != filtros.Status {
			continue
		}
		c := *a
		resultado = append(resultado, &c)
	}
	return resultado, nil
}

func (r *RepositorioDeAtribuicaoMemoria) Atualizar(_ context.Context, atribuicao *entidades.AtribuicaoEdital) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.dados[atribuicao.ID]; !ok {
		return nil
	}
	c := *atribuicao
	r.dados[atribuicao.ID] = &c
	return nil
}

func (r *RepositorioDeAtribuicaoMemoria) ContarAtivasPorAvaliador(_ context.Context, _ int64) (int64, error) {
	return 0, nil
}

func (r *RepositorioDeAtribuicaoMemoria) BuscarHashAnonimizacao(_ context.Context, avaliadorID int64, editalID int64) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, a := range r.dados {
		if a.AvaliadorID == avaliadorID && a.EditalID == editalID {
			return a.HashAnonimizacao.String(), nil
		}
	}
	return "", nil
}
