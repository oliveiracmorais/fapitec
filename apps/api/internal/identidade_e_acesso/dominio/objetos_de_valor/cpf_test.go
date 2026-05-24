package objetos_de_valor

import (
	"testing"
)

func TestCPFValido(t *testing.T) {
	cpf, err := NovoCPF("529.982.247-25")
	if err != nil {
		t.Fatalf("CPF valido retornou erro: %v", err)
	}
	if cpf.String() != "52998224725" {
		t.Fatalf("CPF esperado 52998224725, got %s", cpf.String())
	}
	if cpf.Formatado() != "529.982.247-25" {
		t.Fatalf("CPF formatado esperado 529.982.247-25, got %s", cpf.Formatado())
	}
}

func TestCPFInvalidos(t *testing.T) {
	casos := []struct {
		valor string
		nome  string
	}{
		{"", "vazio"},
		{"123", "curto"},
		{"123456789012", "longo"},
		{"111.111.111-11", "todos digitos iguais"},
		{"123.456.789-00", "digitos verificadores invalidos"},
		{"abc.def.ghi-jk", "nao numerico"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			_, err := NovoCPF(c.valor)
			if err == nil {
				t.Fatalf("CPF %s deveria ser invalido", c.valor)
			}
		})
	}
}

func TestCPFSemPontuacao(t *testing.T) {
	cpf, err := NovoCPF("52998224725")
	if err != nil {
		t.Fatalf("CPF sem pontuacao retornou erro: %v", err)
	}
	if cpf.String() != "52998224725" {
		t.Fatalf("CPF esperado 52998224725, got %s", cpf.String())
	}
}
