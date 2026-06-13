package objetos_de_valor

import (
	"testing"
)

func TestNovoHashAnonimizacao(t *testing.T) {
	t.Run("deve gerar hash nao vazio", func(t *testing.T) {
		h := NovoHashAnonimizacao(1, 10)
		if h.IsZero() {
			t.Error("hash nao deve ser zero")
		}
		if len(h.String()) != 64 {
			t.Errorf("hash SHA-256 deve ter 64 caracteres hex, obteve %d", len(h.String()))
		}
	})

	t.Run("deve ser deterministico", func(t *testing.T) {
		h1 := NovoHashAnonimizacao(1, 10)
		h2 := NovoHashAnonimizacao(1, 10)
		if h1.String() != h2.String() {
			t.Error("hash deve ser deterministico para mesmos parametros")
		}
	})

	t.Run("deve ser diferente para avaliadores diferentes", func(t *testing.T) {
		h1 := NovoHashAnonimizacao(1, 10)
		h2 := NovoHashAnonimizacao(2, 10)
		if h1.String() == h2.String() {
			t.Error("hash deve ser diferente para avaliadores diferentes")
		}
	})

	t.Run("deve ser diferente para editais diferentes", func(t *testing.T) {
		h1 := NovoHashAnonimizacao(1, 10)
		h2 := NovoHashAnonimizacao(1, 20)
		if h1.String() == h2.String() {
			t.Error("hash deve ser diferente para editais diferentes")
		}
	})
}

func TestHashAnonimizacaoExistente(t *testing.T) {
	t.Run("deve criar hash a partir de valor existente", func(t *testing.T) {
		original := NovoHashAnonimizacao(1, 10)
		h, err := HashAnonimizacaoExistente(original.String())
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if h.String() != original.String() {
			t.Errorf("hash deve ser igual ao original")
		}
	})

	t.Run("deve rejeitar valor vazio", func(t *testing.T) {
		_, err := HashAnonimizacaoExistente("")
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})
}

func TestHashIsZero(t *testing.T) {
	t.Run("hash vazio deve ser zero", func(t *testing.T) {
		var h HashAnonimizacao
		if !h.IsZero() {
			t.Error("hash vazio deve ser zero")
		}
	})

	t.Run("hash com valor nao deve ser zero", func(t *testing.T) {
		h := NovoHashAnonimizacao(1, 10)
		if h.IsZero() {
			t.Error("hash com valor nao deve ser zero")
		}
	})
}
