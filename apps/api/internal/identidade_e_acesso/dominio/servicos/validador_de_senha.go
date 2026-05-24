package servicos

import "errors"

type ValidadorDeSenha struct{}

func NovoValidadorDeSenha() ValidadorDeSenha {
	return ValidadorDeSenha{}
}

func (v ValidadorDeSenha) Validar(senha string) error {
	if len(senha) < 8 {
		return errors.New("senha deve ter no minimo 8 caracteres")
	}
	temMaiuscula := false
	temMinuscula := false
	temNumero := false
	temEspecial := false

	for _, c := range senha {
		switch {
		case c >= 'A' && c <= 'Z':
			temMaiuscula = true
		case c >= 'a' && c <= 'z':
			temMinuscula = true
		case c >= '0' && c <= '9':
			temNumero = true
		default:
			temEspecial = true
		}
	}

	if !temMaiuscula {
		return errors.New("senha deve conter pelo menos uma letra maiuscula")
	}
	if !temMinuscula {
		return errors.New("senha deve conter pelo menos uma letra minuscula")
	}
	if !temNumero {
		return errors.New("senha deve conter pelo menos um numero")
	}
	if !temEspecial {
		return errors.New("senha deve conter pelo menos um caractere especial")
	}
	return nil
}
