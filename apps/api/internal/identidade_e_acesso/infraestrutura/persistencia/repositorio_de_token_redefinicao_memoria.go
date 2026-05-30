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
	r.tokens[token.TokenHash] = token
	return nil
}

func (r *RepositorioDeTokenRedefinicaoMemoria) BuscarPorToken(ctx context.Context, token string) (*objetos_de_valor.TokenRedefinicao, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hash := objetos_de_valor.HashToken(token)
	t, ok := r.tokens[hash]
	if !ok {
		return nil, nil
	}
	return &t, nil
}

func (r *RepositorioDeTokenRedefinicaoMemoria) Remover(ctx context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	hash := objetos_de_valor.HashToken(token)
	delete(r.tokens, hash)
	return nil
}
