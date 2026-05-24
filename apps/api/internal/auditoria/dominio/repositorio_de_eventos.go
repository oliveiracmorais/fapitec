package dominio

import "context"

type RepositorioDeEventos interface {
	Inserir(ctx context.Context, evento *EventoAuditoria) error
	Listar(ctx context.Context) ([]EventoAuditoria, error)
}
