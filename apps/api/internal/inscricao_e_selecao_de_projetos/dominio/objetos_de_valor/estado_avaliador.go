package objetos_de_valor

import "fmt"

type EstadoAvaliador string

const (
	EstadoAvaliadorAtivo   EstadoAvaliador = "ativo"
	EstadoAvaliadorInativo EstadoAvaliador = "inativo"
)

func NovoEstadoAvaliador(s string) (EstadoAvaliador, error) {
	switch EstadoAvaliador(s) {
	case EstadoAvaliadorAtivo, EstadoAvaliadorInativo:
		return EstadoAvaliador(s), nil
	default:
		return "", fmt.Errorf("estado de avaliador invalido: %s", s)
	}
}

func (e EstadoAvaliador) String() string {
	return string(e)
}
