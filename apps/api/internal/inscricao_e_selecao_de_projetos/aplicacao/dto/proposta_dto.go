package dto

type CriarPropostaEntrada struct {
	EditalID         int64                    `json:"edital_id"`
	ProponenteID     int64                    `json:"proponente_id"`
	DadosProponente  ProponenteInfoDTO        `json:"dados_proponente"`
	DadosAcademicos  DadosAcademicosDTO       `json:"dados_academicos"`
	EmpresaVinculada string                   `json:"empresa_vinculada"`
	ItensOrcamentarios []ItemOrcamentarioDTO  `json:"itens_orcamentarios"`
	Documentos       []DocumentoPropostaDTO   `json:"documentos"`
}

type ProponenteInfoDTO struct {
	Nome          string `json:"nome"`
	CPF           string `json:"cpf"`
	RG            string `json:"rg"`
	Genero        string `json:"genero"`
	Etnia         string `json:"etnia"`
	DataNascimento string `json:"data_nascimento"`
	Endereco      string `json:"endereco"`
	CEP           string `json:"cep"`
	Logradouro    string `json:"logradouro"`
	Numero        string `json:"numero"`
	Complemento   string `json:"complemento"`
	Bairro        string `json:"bairro"`
	Cidade        string `json:"cidade"`
	UF            string `json:"uf"`
	Telefone      string `json:"telefone"`
	Email         string `json:"email"`
}

type DadosAcademicosDTO struct {
	MaiorTitulacao  string `json:"maior_titulacao"`
	Curso           string `json:"curso"`
	Instituicao     string `json:"instituicao"`
	AnoConclusao    int    `json:"ano_conclusao"`
	AreaConhecimento string `json:"area_conhecimento"`
}

type ItemOrcamentarioDTO struct {
	Descricao     string `json:"descricao"`
	Tipo          string `json:"tipo"`
	Quantidade    int    `json:"quantidade"`
	ValorUnitario int64  `json:"valor_unitario"`
	ValorTotal    int64  `json:"valor_total"`
}

type DocumentoPropostaDTO struct {
	Tipo        string `json:"tipo"`
	NomeArquivo string `json:"nome_arquivo"`
	Caminho     string `json:"caminho"`
}

type PropostaSaida struct {
	ID                 int64                `json:"id"`
	EditalID           int64                `json:"edital_id"`
	ProponenteID       int64                `json:"proponente_id"`
	Protocolo          string               `json:"protocolo"`
	Status             string               `json:"status"`
	Versao             int                  `json:"versao"`
	DadosProponente    ProponenteInfoDTO    `json:"dados_proponente"`
	DadosAcademicos    DadosAcademicosDTO   `json:"dados_academicos"`
	EmpresaVinculada   string               `json:"empresa_vinculada"`
	ItensOrcamentarios []ItemOrcamentarioDTO `json:"itens_orcamentarios"`
	Documentos         []DocumentoPropostaDTO `json:"documentos"`
	ValorTotal         int64                `json:"valor_total"`
	Pareceres          []ParecerDTO         `json:"pareceres"`
	DataSubmissao      string               `json:"data_submissao"`
	DataAtualizacao    string               `json:"data_atualizacao"`
	CriadoEm           string               `json:"criado_em"`
}

type ParecerDTO struct {
	Etapa        string `json:"etapa"`
	AvaliadorID  int64  `json:"avaliador_id"`
	Nota         int    `json:"nota"`
	ParecerTexto string `json:"parecer_texto"`
	Data         string `json:"data"`
}

type AtualizarPropostaEntrada struct {
	DadosProponente    *ProponenteInfoDTO      `json:"dados_proponente,omitempty"`
	DadosAcademicos    *DadosAcademicosDTO     `json:"dados_academicos,omitempty"`
	EmpresaVinculada   *string                 `json:"empresa_vinculada,omitempty"`
	ItensOrcamentarios *[]ItemOrcamentarioDTO  `json:"itens_orcamentarios,omitempty"`
	Documentos         *[]DocumentoPropostaDTO `json:"documentos,omitempty"`
}
