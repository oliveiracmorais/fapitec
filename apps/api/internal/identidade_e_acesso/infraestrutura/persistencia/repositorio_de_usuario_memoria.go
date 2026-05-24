package persistencia

import (
	"context"
	"sync"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
)

type RepositorioDeUsuarioMemoria struct {
	mu      sync.RWMutex
	usuarios []entidades.Usuario
	proximoID int64
}

func NovoRepositorioDeUsuarioMemoria() *RepositorioDeUsuarioMemoria {
	return &RepositorioDeUsuarioMemoria{
		proximoID: 1,
	}
}

func (r *RepositorioDeUsuarioMemoria) Inserir(ctx context.Context, usuario *entidades.Usuario) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	usuario.ID = r.proximoID
	usuario.CriadoEm = time.Now()
	r.proximoID++
	r.usuarios = append(r.usuarios, *usuario)
	return nil
}

func (r *RepositorioDeUsuarioMemoria) BuscarPorCPF(ctx context.Context, cpf string) (*entidades.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.usuarios {
		if u.CPF == cpf {
			u2 := u
			return &u2, nil
		}
	}
	return nil, nil
}

func (r *RepositorioDeUsuarioMemoria) BuscarPorEmail(ctx context.Context, email string) (*entidades.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.usuarios {
		if u.Email.String() == email {
			u2 := u
			return &u2, nil
		}
	}
	return nil, nil
}

func (r *RepositorioDeUsuarioMemoria) AtualizarTentativas(ctx context.Context, id int64, tentativas int, bloqueadoAte *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.usuarios {
		if r.usuarios[i].ID == id {
			r.usuarios[i].Tentativas = tentativas
			if bloqueadoAte != nil {
				t, err := time.Parse(time.RFC3339, *bloqueadoAte)
				if err == nil {
					r.usuarios[i].BloqueadoAte = &t
				}
			} else {
				r.usuarios[i].BloqueadoAte = nil
			}
			return nil
		}
	}
	return nil
}

func (r *RepositorioDeUsuarioMemoria) AtualizarSenha(ctx context.Context, id int64, senhaHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.usuarios {
		if r.usuarios[i].ID == id {
			r.usuarios[i].SenhaHash = objetos_de_valor.NovaSenhaHash(senhaHash)
			return nil
		}
	}
	return nil
}


