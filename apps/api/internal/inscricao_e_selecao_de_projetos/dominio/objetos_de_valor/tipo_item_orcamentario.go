package objetos_de_valor

import "fmt"

type TipoItemOrcamentario string

const (
	TipoItemConsumo    TipoItemOrcamentario = "consumo"
	TipoItemPermanente TipoItemOrcamentario = "permanente"
)

func NovoTipoItemOrcamentario(s string) (TipoItemOrcamentario, error) {
	switch TipoItemOrcamentario(s) {
	case TipoItemConsumo, TipoItemPermanente:
		return TipoItemOrcamentario(s), nil
	default:
		return "", fmt.Errorf("tipo de item orcamentario invalido: %s", s)
	}
}

func (t TipoItemOrcamentario) String() string {
	return string(t)
}
