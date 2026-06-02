package servicos

import "context"

type AutorizacaoChecker interface {
	VerificarPermissao(ctx context.Context, usuarioID, perfil, modulo, operacao string) (bool, error)
}
