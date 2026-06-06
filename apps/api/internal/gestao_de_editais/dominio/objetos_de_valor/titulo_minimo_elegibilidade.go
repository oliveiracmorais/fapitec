package objetos_de_valor

import "fmt"

type TituloMinimoElegibilidade string

const (
	TituloMinimoNaoExigido TituloMinimoElegibilidade = ""
	TituloMinimoGraduado   TituloMinimoElegibilidade = "Graduado"
	TituloMinimoMestre     TituloMinimoElegibilidade = "Mestre"
	TituloMinimoDoutor     TituloMinimoElegibilidade = "Doutor"
)

func NovoTituloMinimoElegibilidade(s string) (TituloMinimoElegibilidade, error) {
	switch TituloMinimoElegibilidade(s) {
	case TituloMinimoNaoExigido, TituloMinimoGraduado, TituloMinimoMestre, TituloMinimoDoutor:
		return TituloMinimoElegibilidade(s), nil
	default:
		return TituloMinimoNaoExigido, fmt.Errorf("titulo minimo de elegibilidade invalido: %s", s)
	}
}

func (t TituloMinimoElegibilidade) String() string {
	return string(t)
}
