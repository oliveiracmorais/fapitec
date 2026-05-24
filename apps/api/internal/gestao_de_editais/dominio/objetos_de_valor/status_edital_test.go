package objetos_de_valor

import "testing"

func TestNovoStatusEdital(t *testing.T) {
	t.Run("deve criar status valido", func(t *testing.T) {
		status, err := NovoStatusEdital("ativo")
		if err != nil {
			t.Fatalf("esperava nil, got %v", err)
		}
		if status != StatusEditalAtivo {
			t.Errorf("esperava 'ativo', got %s", status)
		}
	})

	t.Run("deve rejeitar status invalido", func(t *testing.T) {
		_, err := NovoStatusEdital("invalido")
		if err == nil {
			t.Fatal("esperava erro, got nil")
		}
	})

	t.Run("deve criar status encerrado", func(t *testing.T) {
		status, err := NovoStatusEdital("encerrado")
		if err != nil {
			t.Fatalf("esperava nil, got %v", err)
		}
		if status != StatusEditalEncerrado {
			t.Errorf("esperava 'encerrado', got %s", status)
		}
	})
}
