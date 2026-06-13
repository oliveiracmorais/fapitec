package objetos_de_valor

import (
	"crypto/sha256"
	"fmt"
	"os"
)

type HashAnonimizacao struct {
	valor string
}

func salt() string {
	s := os.Getenv("AVALIADOR_HASH_SALT")
	if s == "" {
		s = "fapitec-default-salt"
	}
	return s
}

func NovoHashAnonimizacao(avaliadorID, editalID int64) HashAnonimizacao {
	dados := fmt.Sprintf("%d:%d:%s", avaliadorID, editalID, salt())
	hash := sha256.Sum256([]byte(dados))
	return HashAnonimizacao{valor: fmt.Sprintf("%x", hash)}
}

func HashAnonimizacaoExistente(valor string) (HashAnonimizacao, error) {
	if valor == "" {
		return HashAnonimizacao{}, fmt.Errorf("hash de anonimizacao nao pode ser vazio")
	}
	return HashAnonimizacao{valor: valor}, nil
}

func (h HashAnonimizacao) String() string {
	return h.valor
}

func (h HashAnonimizacao) IsZero() bool {
	return h.valor == ""
}
