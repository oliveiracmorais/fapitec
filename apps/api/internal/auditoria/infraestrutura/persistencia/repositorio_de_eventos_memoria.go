package persistencia

import (
	"context"
	"sync"

	"github.com/oliveiracmorais/fapitec/api/internal/auditoria/dominio"
)

type RepositorioDeEventosMemoria struct {
	mu     sync.RWMutex
	eventos []dominio.EventoAuditoria
	seq    int64
}

func NovoRepositorioDeEventosMemoria() *RepositorioDeEventosMemoria {
	return &RepositorioDeEventosMemoria{
		eventos: make([]dominio.EventoAuditoria, 0),
	}
}

func (r *RepositorioDeEventosMemoria) Inserir(ctx context.Context, evento *dominio.EventoAuditoria) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	evento.ID = r.seq
	r.eventos = append(r.eventos, *evento)
	return nil
}

func (r *RepositorioDeEventosMemoria) Listar(ctx context.Context) ([]dominio.EventoAuditoria, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]dominio.EventoAuditoria, len(r.eventos))
	copy(result, r.eventos)
	return result, nil
}
