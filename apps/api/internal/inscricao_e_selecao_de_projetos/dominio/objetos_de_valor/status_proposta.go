package objetos_de_valor

import "fmt"

type StatusProposta string

const (
	StatusPropostaRascunho     StatusProposta = "rascunho"
	StatusPropostaSubmetida    StatusProposta = "submetida"
	StatusPropostaEmAvaliacao  StatusProposta = "em_avaliacao"
	StatusPropostaEmRecurso    StatusProposta = "em_recurso"
	StatusPropostaAprovada     StatusProposta = "aprovada"
	StatusPropostaReprovada    StatusProposta = "reprovada"
)

func NovoStatusProposta(s string) (StatusProposta, error) {
	switch StatusProposta(s) {
	case StatusPropostaRascunho, StatusPropostaSubmetida,
		StatusPropostaEmAvaliacao, StatusPropostaEmRecurso,
		StatusPropostaAprovada, StatusPropostaReprovada:
		return StatusProposta(s), nil
	default:
		return "", fmt.Errorf("status de proposta invalido: %s", s)
	}
}

func (s StatusProposta) String() string {
	return string(s)
}
