package objetos_de_valor

import "testing"

func TestEmailValido(t *testing.T) {
	email, err := NovoEmail("usuario@example.com")
	if err != nil {
		t.Fatalf("Email valido retornou erro: %v", err)
	}
	if email.String() != "usuario@example.com" {
		t.Fatalf("Email esperado usuario@example.com, got %s", email.String())
	}
}

func TestEmailComSubdominio(t *testing.T) {
	_, err := NovoEmail("usuario@sub.example.com.br")
	if err != nil {
		t.Fatalf("Email com subdominio deveria ser valido: %v", err)
	}
}

func TestEmailInvalido(t *testing.T) {
	casos := []struct {
		valor string
		nome  string
	}{
		{"", "vazio"},
		{"invalido", "sem arroba"},
		{"@example.com", "sem usuario"},
		{"usuario@", "sem dominio"},
		{"usuario@.com", "dominio incompleto"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			_, err := NovoEmail(c.valor)
			if err == nil {
				t.Fatalf("Email %s deveria ser invalido", c.valor)
			}
		})
	}
}
