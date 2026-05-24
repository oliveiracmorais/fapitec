package servicos

import "testing"

func TestSenhaValida(t *testing.T) {
	v := NovoValidadorDeSenha()
	err := v.Validar("Senha@123")
	if err != nil {
		t.Fatalf("Senha valida retornou erro: %v", err)
	}
}

func TestSenhaCurta(t *testing.T) {
	v := NovoValidadorDeSenha()
	err := v.Validar("Ab@1")
	if err == nil {
		t.Fatal("Senha curta deveria ser invalida")
	}
}

func TestSenhaSemMaiuscula(t *testing.T) {
	v := NovoValidadorDeSenha()
	err := v.Validar("senha@123")
	if err == nil {
		t.Fatal("Senha sem maiuscula deveria ser invalida")
	}
}

func TestSenhaSemMinuscula(t *testing.T) {
	v := NovoValidadorDeSenha()
	err := v.Validar("SENHA@123")
	if err == nil {
		t.Fatal("Senha sem minuscula deveria ser invalida")
	}
}

func TestSenhaSemNumero(t *testing.T) {
	v := NovoValidadorDeSenha()
	err := v.Validar("Senha@abc")
	if err == nil {
		t.Fatal("Senha sem numero deveria ser invalida")
	}
}

func TestSenhaSemEspecial(t *testing.T) {
	v := NovoValidadorDeSenha()
	err := v.Validar("Senha1234")
	if err == nil {
		t.Fatal("Senha sem caractere especial deveria ser invalida")
	}
}
