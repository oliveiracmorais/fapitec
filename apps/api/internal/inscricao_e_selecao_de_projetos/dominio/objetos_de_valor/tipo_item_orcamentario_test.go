package objetos_de_valor

import "testing"

func TestNovoTipoItemOrcamentario(t *testing.T) {
	t.Run("deve criar tipo consumo", func(t *testing.T) {
		tipo, err := NovoTipoItemOrcamentario("consumo")
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if tipo != TipoItemConsumo {
			t.Errorf("esperava consumo, obteve %v", tipo)
		}
	})

	t.Run("deve rejeitar tipo invalido", func(t *testing.T) {
		_, err := NovoTipoItemOrcamentario("invalido")
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})
}
