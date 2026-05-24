package casos_de_uso

import "context"

type RegistrarEventoInput struct {
	Acao      string
	AtorID    int64
	AtorNome  string
	AtorCPF   string
	Resultado string
	Contexto  map[string]string
}

type RegistradorAuditoria interface {
	Registrar(ctx context.Context, input RegistrarEventoInput) error
}
