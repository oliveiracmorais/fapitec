package objetos_de_valor

import "testing"

func TestNovoStatusConvite(t *testing.T) {
	t.Run("deve criar status pendente", func(t *testing.T) {
		s, err := NovoStatusConvite("pendente")
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if s != StatusConvitePendente {
			t.Errorf("esperava pendente, obteve %s", s)
		}
	})

	t.Run("deve criar status aceito", func(t *testing.T) {
		s, err := NovoStatusConvite("aceito")
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if s != StatusConviteAceito {
			t.Errorf("esperava aceito, obteve %s", s)
		}
	})

	t.Run("deve criar status recusado", func(t *testing.T) {
		s, err := NovoStatusConvite("recusado")
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if s != StatusConviteRecusado {
			t.Errorf("esperava recusado, obteve %s", s)
		}
	})

	t.Run("deve rejeitar status invalido", func(t *testing.T) {
		_, err := NovoStatusConvite("cancelado")
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("string deve retornar o valor", func(t *testing.T) {
		s, _ := NovoStatusConvite("pendente")
		if s.String() != "pendente" {
			t.Errorf("esperava 'pendente', obteve %s", s.String())
		}
	})
}
