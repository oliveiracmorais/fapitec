package repositorios

import (
	"context"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/entidades"
)

type RepositorioDeUsuario interface {
	Inserir(ctx context.Context, usuario *entidades.Usuario) error
	BuscarPorCPF(ctx context.Context, cpf string) (*entidades.Usuario, error)
	BuscarPorEmail(ctx context.Context, email string) (*entidades.Usuario, error)
	AtualizarTentativas(ctx context.Context, id int64, tentativas int, bloqueadoAte *string) error
	AtualizarSenha(ctx context.Context, id int64, senhaHash string) error
}
