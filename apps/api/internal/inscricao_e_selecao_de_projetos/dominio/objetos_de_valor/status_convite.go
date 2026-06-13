package objetos_de_valor

import "fmt"

type StatusConvite string

const (
	StatusConvitePendente StatusConvite = "pendente"
	StatusConviteAceito   StatusConvite = "aceito"
	StatusConviteRecusado StatusConvite = "recusado"
)

func NovoStatusConvite(s string) (StatusConvite, error) {
	switch StatusConvite(s) {
	case StatusConvitePendente, StatusConviteAceito, StatusConviteRecusado:
		return StatusConvite(s), nil
	default:
		return "", fmt.Errorf("status de convite invalido: %s", s)
	}
}

func (s StatusConvite) String() string {
	return string(s)
}
