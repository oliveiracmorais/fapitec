package objetos_de_valor

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	valor string
}

func NovoEmail(valor string) (Email, error) {
	if !emailRegex.MatchString(valor) {
		return Email{}, errors.New("email invalido")
	}
	return Email{valor: valor}, nil
}

func (e Email) String() string {
	return e.valor
}

func NewEmailSeguro(valor string) Email {
	email, _ := NovoEmail(valor)
	return email
}
