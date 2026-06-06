package objetos_de_valor

import "fmt"

type TipoChamada string

const (
	TipoChamadaAPQ TipoChamada = "APQ"
	TipoChamadaARC TipoChamada = "ARC"
)

func NovoTipoChamada(s string) (TipoChamada, error) {
	switch TipoChamada(s) {
	case TipoChamadaAPQ, TipoChamadaARC:
		return TipoChamada(s), nil
	default:
		return "", fmt.Errorf("tipo de chamada invalido: %s", s)
	}
}

func (t TipoChamada) String() string {
	return string(t)
}
