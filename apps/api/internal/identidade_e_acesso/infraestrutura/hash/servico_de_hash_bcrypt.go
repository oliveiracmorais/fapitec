package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type ServicoDeHashBcrypt struct {
	custo int
}

func NovoServicoDeHashBcrypt() ServicoDeHashBcrypt {
	return ServicoDeHashBcrypt{custo: bcrypt.DefaultCost}
}

func (s ServicoDeHashBcrypt) Hash(senha string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), s.custo)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s ServicoDeHashBcrypt) Comparar(hash, senha string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
	return err == nil
}
