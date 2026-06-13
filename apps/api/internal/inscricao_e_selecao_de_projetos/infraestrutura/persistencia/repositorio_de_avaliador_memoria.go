package persistencia

import (
	"context"
	"strings"
	"sync"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type RepositorioDeAvaliadorMemoria struct {
	mu       sync.RWMutex
	dados    map[int64]*entidades.Avaliador
	cpfIdx   map[string]int64
	proximoID int64
}

func NovoRepositorioDeAvaliadorMemoria() *RepositorioDeAvaliadorMemoria {
	return &RepositorioDeAvaliadorMemoria{
		dados:    make(map[int64]*entidades.Avaliador),
		cpfIdx:   make(map[string]int64),
		proximoID: 1,
	}
}

func (r *RepositorioDeAvaliadorMemoria) Criar(_ context.Context, avaliador *entidades.Avaliador) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.proximoID
	r.proximoID++
	avaliador.ID = id
	r.dados[id] = copiarAvaliador(avaliador)
	r.cpfIdx[avaliador.CPF] = id
	return nil
}

func (r *RepositorioDeAvaliadorMemoria) BuscarPorID(_ context.Context, id int64) (*entidades.Avaliador, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, ok := r.dados[id]
	if !ok {
		return nil, nil
	}
	return copiarAvaliador(a), nil
}

func (r *RepositorioDeAvaliadorMemoria) BuscarPorCPF(_ context.Context, cpf string) (*entidades.Avaliador, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.cpfIdx[cpf]
	if !ok {
		return nil, nil
	}
	a, ok := r.dados[id]
	if !ok {
		return nil, nil
	}
	return copiarAvaliador(a), nil
}

func (r *RepositorioDeAvaliadorMemoria) BuscarPorUsuarioID(_ context.Context, usuarioID int64) (*entidades.Avaliador, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, a := range r.dados {
		if a.UsuarioID == usuarioID {
			return copiarAvaliador(a), nil
		}
	}
	return nil, nil
}

func (r *RepositorioDeAvaliadorMemoria) Listar(_ context.Context, filtros repositorios.FiltrosListarAvaliadores) ([]*entidades.Avaliador, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resultado []*entidades.Avaliador
	for _, a := range r.dados {
		if filtros.Nome != "" && !strings.Contains(strings.ToLower(a.Nome), strings.ToLower(filtros.Nome)) {
			continue
		}
		if filtros.CPF != "" && !strings.Contains(a.CPF, filtros.CPF) {
			continue
		}
		if filtros.AreaConhecimento != "" && !strings.Contains(strings.ToLower(a.AreaConhecimento), strings.ToLower(filtros.AreaConhecimento)) {
			continue
		}
		if filtros.Estado != "" && a.Estado.String() != filtros.Estado {
			continue
		}
		resultado = append(resultado, copiarAvaliador(a))
	}
	return resultado, nil
}

func (r *RepositorioDeAvaliadorMemoria) Atualizar(_ context.Context, avaliador *entidades.Avaliador) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.dados[avaliador.ID]; !ok {
		return nil
	}

	antigo := r.dados[avaliador.ID]
	if antigo.CPF != avaliador.CPF {
		delete(r.cpfIdx, antigo.CPF)
		r.cpfIdx[avaliador.CPF] = avaliador.ID
	}

	r.dados[avaliador.ID] = copiarAvaliador(avaliador)
	return nil
}

func (r *RepositorioDeAvaliadorMemoria) ContarPropostasPorAvaliador(_ context.Context, _ int64) (int64, error) {
	return 0, nil
}

func copiarAvaliador(a *entidades.Avaliador) *entidades.Avaliador {
	if a == nil {
		return nil
	}
	c := *a
	return &c
}
