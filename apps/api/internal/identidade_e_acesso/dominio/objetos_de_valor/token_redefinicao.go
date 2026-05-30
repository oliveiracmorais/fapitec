package objetos_de_valor

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

type TokenRedefinicao struct {
	Token     string
	TokenHash string
	UsuarioID int64
	ExpiradoEm time.Time
}

func NovoTokeRedefinicao(usuarioID int64, duracao time.Duration) TokenRedefinicao {
	b := make([]byte, 16)
	rand.Read(b)
	token := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
	return TokenRedefinicao{
		Token:      token,
		TokenHash:  hash,
		UsuarioID:  usuarioID,
		ExpiradoEm: time.Now().Add(duracao),
	}
}

func HashToken(token string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
}

func (t TokenRedefinicao) Expirado() bool {
	return time.Now().After(t.ExpiradoEm)
}
