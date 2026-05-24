package entidades

import (
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
)

type Usuario struct {
	ID           int64
	Nome         string
	CPF          string
	Estrangeiro  bool
	Email        objetos_de_valor.Email
	SenhaHash    objetos_de_valor.SenhaHash
	CriadoEm     time.Time
	Tentativas   int
	BloqueadoAte *time.Time
}

func (u *Usuario) EstaBloqueado() bool {
	if u.BloqueadoAte == nil {
		return false
	}
	return time.Now().Before(*u.BloqueadoAte)
}

func (u *Usuario) RegistrarTentativaFalha(tempoBloqueio time.Duration) {
	u.Tentativas++
	if u.Tentativas >= 5 {
		agora := time.Now()
		bloqueio := agora.Add(tempoBloqueio)
		u.BloqueadoAte = &bloqueio
	}
}

func (u *Usuario) ResetarTentativas() {
	u.Tentativas = 0
	u.BloqueadoAte = nil
}
