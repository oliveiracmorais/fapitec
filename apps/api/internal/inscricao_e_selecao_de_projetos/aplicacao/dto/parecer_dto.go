package dto

type EmitirParecerEntrada struct {
	PropostaID   int64  `json:"proposta_id"`
	Etapa        string `json:"etapa"`
	AvaliadorID  int64  `json:"avaliador_id"`
	Nota         int    `json:"nota"`
	ParecerTexto string `json:"parecer_texto"`
}

type ParecerSaida struct {
	ID           int64  `json:"id"`
	PropostaID   int64  `json:"proposta_id"`
	Etapa        string `json:"etapa"`
	Nota         int    `json:"nota"`
	ParecerTexto string `json:"parecer_texto"`
	Data         string `json:"data"`
}

type ParecerAnonimizadoSaida struct {
	ID           int64  `json:"id"`
	PropostaID   int64  `json:"proposta_id"`
	Etapa        string `json:"etapa"`
	HashAvaliador string `json:"hash_avaliador"`
	Nota         int    `json:"nota"`
	ParecerTexto string `json:"parecer_texto"`
	Data         string `json:"data"`
}

type PropostaParaAvaliarSaida struct {
	ID                 int64  `json:"id"`
	EditalID           int64  `json:"edital_id"`
	Protocolo          string `json:"protocolo"`
	Status             string `json:"status"`
	DadosProponente    ProponenteInfoDTO `json:"dados_proponente"`
	DadosAcademicos    DadosAcademicosDTO `json:"dados_academicos"`
	ValorTotal         int64  `json:"valor_total"`
	Pareceres          []ParecerSaida `json:"pareceres"`
	DataSubmissao      string `json:"data_submissao"`
}

type FinalizarAvaliacaoEntrada struct {
	NotaDeCorte int `json:"nota_de_corte"`
}

type ClassificacaoSaida struct {
	PropostaID  int64  `json:"proposta_id"`
	Protocolo   string `json:"protocolo"`
	NotaFinal   int    `json:"nota_final"`
	Status      string `json:"status"`
}

func ParaParecerSaida(parecer *ParecerSaida) *ParecerSaida {
	return parecer
}
