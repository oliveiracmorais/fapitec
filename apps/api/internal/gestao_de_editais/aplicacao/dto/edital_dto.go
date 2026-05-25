package dto

type CriarEditalEntrada struct {
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	DataInicio  string `json:"data_inicio"`
	DataFim     string `json:"data_fim"`
	TipoChamada string `json:"tipo_chamada"`
}

type AtualizarEditalEntrada struct {
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	DataInicio  string `json:"data_inicio"`
	DataFim     string `json:"data_fim"`
	Status      string `json:"status"`
	TipoChamada string `json:"tipo_chamada"`
}

type EditalSaida struct {
	ID          int64  `json:"id"`
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	DataInicio  string `json:"data_inicio"`
	DataFim     string `json:"data_fim"`
	Status      string `json:"status"`
	TipoChamada string `json:"tipo_chamada"`
	CriadoEm    string `json:"criado_em"`
}

type ListarEditaisSaida struct {
	Editais []EditalSaida `json:"editais"`
	Total   int           `json:"total"`
}
