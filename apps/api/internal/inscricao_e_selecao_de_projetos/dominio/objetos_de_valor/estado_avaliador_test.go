package objetos_de_valor

import "testing"

func TestNovoEstadoAvaliador(t *testing.T) {
	t.Run("deve criar estado ativo", func(t *testing.T) {
		e, err := NovoEstadoAvaliador("ativo")
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if e != EstadoAvaliadorAtivo {
			t.Errorf("esperava ativo, obteve %s", e)
		}
	})

	t.Run("deve criar estado inativo", func(t *testing.T) {
		e, err := NovoEstadoAvaliador("inativo")
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if e != EstadoAvaliadorInativo {
			t.Errorf("esperava inativo, obteve %s", e)
		}
	})

	t.Run("deve rejeitar estado invalido", func(t *testing.T) {
		_, err := NovoEstadoAvaliador("suspenso")
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("string deve retornar o valor", func(t *testing.T) {
		e, _ := NovoEstadoAvaliador("ativo")
		if e.String() != "ativo" {
			t.Errorf("esperava 'ativo', obteve %s", e.String())
		}
	})
}
