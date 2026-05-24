package objetos_de_valor

import "fmt"

type StatusEdital string

const (
	StatusEditalAtivo        StatusEdital = "ativo"
	StatusEditalEncerrado    StatusEdital = "encerrado"
	StatusEditalEmAvaliacao  StatusEdital = "em_avaliacao"
)

func NovoStatusEdital(s string) (StatusEdital, error) {
	switch StatusEdital(s) {
	case StatusEditalAtivo, StatusEditalEncerrado, StatusEditalEmAvaliacao:
		return StatusEdital(s), nil
	default:
		return "", fmt.Errorf("status de edital invalido: %s", s)
	}
}

func (s StatusEdital) String() string {
	return string(s)
}
