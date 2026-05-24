package adaptadores

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/auditoria/dominio"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/casos_de_uso"
)

type RepositorioDeEventos interface {
	Inserir(ctx context.Context, evento *dominio.EventoAuditoria) error
}

type RegistradorAuditoria struct {
	repo RepositorioDeEventos
}

func NovoRegistradorAuditoria(repo RepositorioDeEventos) *RegistradorAuditoria {
	return &RegistradorAuditoria{repo: repo}
}

var ResultadoEventoValido = map[dominio.ResultadoEvento]bool{
	dominio.ResultadoSucesso: true,
	dominio.ResultadoFalha:   true,
	dominio.ResultadoNegado:  true,
	dominio.ResultadoErro:    true,
}

func (r *RegistradorAuditoria) Registrar(ctx context.Context, input casos_de_uso.RegistrarEventoInput) error {
	resultado := dominio.ResultadoEvento(input.Resultado)
	if !ResultadoEventoValido[resultado] {
		return fmt.Errorf("resultado de evento invalido: %s", input.Resultado)
	}
	evento := dominio.NovoEventoAuditoria(input.Acao, input.AtorID, input.AtorNome, input.AtorCPF, "", resultado)
	if input.Contexto != nil {
		evento.Contexto = input.Contexto
	}
	return r.repo.Inserir(ctx, &evento)
}

var _ casos_de_uso.RegistradorAuditoria = (*RegistradorAuditoria)(nil)
