package objetos_de_valor

import (
	"strings"
	"testing"
)

func TestProtocolo(t *testing.T) {
	t.Run("deve criar protocolo valido", func(t *testing.T) {
		p := NovoProtocolo(5, 1)
		parts := strings.Split(p.String(), ".")
		if len(parts) != 3 {
			t.Errorf("esperava 3 partes, obteve %d", len(parts))
		}
		if parts[1] != "5" {
			t.Errorf("esperava edital 5, obteve %s", parts[1])
		}
		if parts[2] != "001" {
			t.Errorf("esperava sequencial 001, obteve %s", parts[2])
		}
	})

	t.Run("deve criar protocolo existente", func(t *testing.T) {
		p, err := ProtocoloExistente("2026.5.001")
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if p.String() != "2026.5.001" {
			t.Errorf("esperava 2026.5.001, obteve %s", p.String())
		}
	})

	t.Run("deve rejeitar protocolo vazio", func(t *testing.T) {
		_, err := ProtocoloExistente("")
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})
}
