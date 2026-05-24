package dominio

import "time"

type ResultadoEvento string

const (
	ResultadoSucesso ResultadoEvento = "sucesso"
	ResultadoFalha   ResultadoEvento = "falha"
	ResultadoNegado  ResultadoEvento = "negado"
	ResultadoErro    ResultadoEvento = "erro"
)

type EventoAuditoria struct {
	ID        int64
	Acao      string
	AtorID    int64
	AtorNome  string
	AtorCPF   string
	Perfil    string
	Resultado ResultadoEvento
	Modulo    string
	Recurso   string
	Origem    string
	DataHora  time.Time
	Contexto  map[string]string
}

func NovoEventoAuditoria(acao string, atorID int64, atorNome, atorCPF, perfil string, resultado ResultadoEvento) EventoAuditoria {
	return EventoAuditoria{
		Acao:      acao,
		AtorID:    atorID,
		AtorNome:  atorNome,
		AtorCPF:   atorCPF,
		Perfil:    perfil,
		Resultado: resultado,
		DataHora:  time.Now(),
		Contexto:  make(map[string]string),
	}
}
