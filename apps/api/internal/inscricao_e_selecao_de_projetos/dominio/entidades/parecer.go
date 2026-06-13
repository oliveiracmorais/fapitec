package entidades

import (
	"fmt"
	"time"
)

type Parecer struct {
	ID           int64
	PropostaID   int64
	Etapa        string
	AvaliadorID  int64
	Nota         int
	ParecerTexto string
	Data         time.Time
}

type NovoParecerParams struct {
	PropostaID   int64
	Etapa        string
	AvaliadorID  int64
	Nota         int
	ParecerTexto string
}

func NovoParecer(params NovoParecerParams) (*Parecer, error) {
	if params.PropostaID == 0 {
		return nil, fmt.Errorf("proposta é obrigatoria")
	}
	if params.Etapa == "" {
		return nil, fmt.Errorf("etapa é obrigatoria")
	}
	if params.AvaliadorID == 0 {
		return nil, fmt.Errorf("avaliador é obrigatorio")
	}
	if params.Nota < 0 || params.Nota > 100 {
		return nil, fmt.Errorf("nota deve estar entre 0 e 100")
	}
	if params.ParecerTexto == "" {
		return nil, fmt.Errorf("parecer textual é obrigatorio")
	}

	return &Parecer{
		PropostaID:   params.PropostaID,
		Etapa:        params.Etapa,
		AvaliadorID:  params.AvaliadorID,
		Nota:         params.Nota,
		ParecerTexto: params.ParecerTexto,
		Data:         time.Now(),
	}, nil
}
