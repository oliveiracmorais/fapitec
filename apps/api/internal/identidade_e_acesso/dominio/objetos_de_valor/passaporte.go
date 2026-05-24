package objetos_de_valor

import (
	"errors"
	"regexp"
)

var passaporteRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3,30}$`)

type Passaporte struct {
	valor string
}

func NovoPassaporte(valor string) (Passaporte, error) {
	if !passaporteRegex.MatchString(valor) {
		return Passaporte{}, errors.New("passaporte deve ter entre 3 e 30 caracteres alfanumericos")
	}
	return Passaporte{valor: valor}, nil
}

func (p Passaporte) String() string {
	return p.valor
}
