package objetos_de_valor

import "testing"

func TestPassaporteValido(t *testing.T) {
	casos := []string{"AB123456", "ABC123", "A1B2C3D4E5F6", "abc123"}
	for _, v := range casos {
		t.Run(v, func(t *testing.T) {
			p, err := NovoPassaporte(v)
			if err != nil {
				t.Fatalf("Passaporte valido %s retornou erro: %v", v, err)
			}
			if p.String() != v {
				t.Fatalf("Passaporte esperado %s, got %s", v, p.String())
			}
		})
	}
}

func TestPassaporteInvalido(t *testing.T) {
	casos := []struct {
		valor string
		nome  string
	}{
		{"", "vazio"},
		{"AB", "curto"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ12345", "longo"},
		{"AB 123", "com espaco"},
		{"AB-123", "com especial"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			_, err := NovoPassaporte(c.valor)
			if err == nil {
				t.Fatalf("Passaporte %s deveria ser invalido", c.valor)
			}
		})
	}
}
