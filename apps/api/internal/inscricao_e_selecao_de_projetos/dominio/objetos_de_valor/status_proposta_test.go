package objetos_de_valor

import "testing"

func TestNovoStatusProposta(t *testing.T) {
	t.Run("deve criar status valido", func(t *testing.T) {
		status, err := NovoStatusProposta("rascunho")
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if status != StatusPropostaRascunho {
			t.Errorf("esperava rascunho, obteve %v", status)
		}
	})

	t.Run("deve rejeitar status invalido", func(t *testing.T) {
		_, err := NovoStatusProposta("invalido")
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve converter para string", func(t *testing.T) {
		status, _ := NovoStatusProposta("submetida")
		if status.String() != "submetida" {
			t.Errorf("esperava submetida, obteve %s", status.String())
		}
	})
}
