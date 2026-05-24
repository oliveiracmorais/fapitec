package persistencia

import (
	"context"
	"sync"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
)

type RepositorioDeTokenRedefinicaoMemoria struct {
	mu     sync.RWMutex
	tokens map[string]objetos_de_valor.TokenRedefinicao
}

func NovoRepositorioDeTokenRedefinicaoMemoria() *RepositorioDeTokenRedefinicaoMemoria {
	return &RepositorioDeTokenRedefinicaoMemoria{
		tokens: make(map[string]objetos_de_valor.TokenRedefinicao),
	}
}

func (r *RepositorioDeTokenRedefinicaoMemoria) Inserir(ctx context.Context, token objetos_de_valor.TokenRedefinicao) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token.Token] = token
	return nil
}

func (r *RepositorioDeTokenRedefinicaoMemoria) BuscarPorToken(ctx context.Context, token string) (*objetos_de_valor.TokenRedefinicao, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tokens[token]
	if !ok {
		return nil, nil
	}
	return &t, nil
}

func (r *RepositorioDeTokenRedefinicaoMemoria) Remover(ctx context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokens, token)
	return nil
}
