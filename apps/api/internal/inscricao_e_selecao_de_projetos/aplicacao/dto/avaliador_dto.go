package dto

type CadastrarAvaliadorEntrada struct {
	UsuarioID         int64  `json:"usuario_id"`
	Nome              string `json:"nome"`
	CPF               string `json:"cpf"`
	Email             string `json:"email"`
	TitulacaoMaxima   string `json:"titulacao_maxima"`
	AreaConhecimento  string `json:"area_conhecimento"`
	Instituicao       string `json:"instituicao"`
	CurriculoResumido string `json:"curriculo_resumido"`
}

type AtualizarAvaliadorEntrada struct {
	Nome              *string `json:"nome,omitempty"`
	CPF               *string `json:"cpf,omitempty"`
	Email             *string `json:"email,omitempty"`
	TitulacaoMaxima   *string `json:"titulacao_maxima,omitempty"`
	AreaConhecimento  *string `json:"area_conhecimento,omitempty"`
	Instituicao       *string `json:"instituicao,omitempty"`
	CurriculoResumido *string `json:"curriculo_resumido,omitempty"`
	Estado            *string `json:"estado,omitempty"`
}

type AvaliadorSaida struct {
	ID                int64  `json:"id"`
	UsuarioID         int64  `json:"usuario_id"`
	Nome              string `json:"nome"`
	CPF               string `json:"cpf"`
	Email             string `json:"email"`
	TitulacaoMaxima   string `json:"titulacao_maxima"`
	AreaConhecimento  string `json:"area_conhecimento"`
	Instituicao       string `json:"instituicao"`
	CurriculoResumido string `json:"curriculo_resumido"`
	Estado            string `json:"estado"`
	DataCadastro      string `json:"data_cadastro"`
	DataAtualizacao   string `json:"data_atualizacao"`
	TotalPropostas    int64  `json:"total_propostas,omitempty"`
	AtribuicoesAtivas int64  `json:"atribuicoes_ativas,omitempty"`
}

type AtribuirEditalEntrada struct {
	EditalID  int64  `json:"edital_id"`
	DataInicio string `json:"data_inicio"`
	DataFim   string `json:"data_fim"`
}

type AtribuicaoSaida struct {
	ID               int64  `json:"id"`
	AvaliadorID      int64  `json:"avaliador_id"`
	EditalID         int64  `json:"edital_id"`
	DataInicio       string `json:"data_inicio"`
	DataFim          string `json:"data_fim"`
	StatusConvite    string `json:"status_convite"`
	HashAnonimizacao string `json:"hash_anonimizacao"`
	CriadoEm         string `json:"criado_em"`
}

type GerenciarConviteEntrada struct {
	Acao string `json:"acao"`
}
