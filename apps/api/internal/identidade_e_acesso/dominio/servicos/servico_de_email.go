package servicos

import "context"

type ServicoDeEmail interface {
	EnviarRedefinicaoSenha(ctx context.Context, email, token string) error
}
