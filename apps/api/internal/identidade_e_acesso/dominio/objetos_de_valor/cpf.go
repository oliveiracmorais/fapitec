package objetos_de_valor

import (
	"errors"
	"regexp"
	"strconv"
)

type CPF struct {
	valor string
}

func NovoCPF(valor string) (CPF, error) {
	digits := regexp.MustCompile(`\D`).ReplaceAllString(valor, "")
	if len(digits) != 11 {
		return CPF{}, errors.New("CPF deve ter 11 digitos")
	}
	if todosDigitosIguais(digits) {
		return CPF{}, errors.New("CPF invalido: todos os digitos sao iguais")
	}
	if !validarDigitosVerificadores(digits) {
		return CPF{}, errors.New("CPF invalido: digitos verificadores nao conferem")
	}
	return CPF{valor: digits}, nil
}

func (c CPF) String() string {
	return c.valor
}

func (c CPF) Formatado() string {
	v := c.valor
	return v[:3] + "." + v[3:6] + "." + v[6:9] + "-" + v[9:]
}

func todosDigitosIguais(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func validarDigitosVerificadores(cpf string) bool {
	d1 := calcularDigito(cpf[:9], 10)
	d2 := calcularDigito(cpf[:9]+strconv.Itoa(d1), 11)
	return d1 == int(cpf[9]-'0') && d2 == int(cpf[10]-'0')
}

func calcularDigito(base string, pesoInicial int) int {
	soma := 0
	for i := 0; i < len(base); i++ {
		num, _ := strconv.Atoi(string(base[i]))
		soma += num * (pesoInicial - i)
	}
	resto := soma % 11
	if resto < 2 {
		return 0
	}
	return 11 - resto
}

func NewCPFSeguro(valor string) CPF {
	cpf, _ := NovoCPF(valor)
	return cpf
}
