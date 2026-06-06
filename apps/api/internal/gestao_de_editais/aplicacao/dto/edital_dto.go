package dto

type CriarEditalEntrada struct {
	Nome                      string   `json:"nome"`
	Descricao                 string   `json:"descricao"`
	DataInicio                string   `json:"data_inicio"`
	DataFim                   string   `json:"data_fim"`
	TipoChamada               string   `json:"tipo_chamada"`
	ModeloFormulario          int      `json:"modelo_formulario"`
	RelatoriosExigidos        []string `json:"relatorios_exigidos"`
	TituloMinimoElegibilidade string   `json:"titulo_minimo_elegibilidade"`
	ExigeEmpresa              bool     `json:"exige_empresa"`
	PorteEmpresa              []string `json:"porte_empresa"`
	EnquadramentoEmpresa      []string `json:"enquadramento_empresa"`
	DocumentosObrigatorios    []string `json:"documentos_obrigatorios"`
}

type AtualizarEditalEntrada struct {
	Nome                      string   `json:"nome"`
	Descricao                 string   `json:"descricao"`
	DataInicio                string   `json:"data_inicio"`
	DataFim                   string   `json:"data_fim"`
	Status                    string   `json:"status"`
	TipoChamada               string   `json:"tipo_chamada"`
	NotaDeCorte               *int     `json:"nota_de_corte"`
	ValorGlobal               *int64   `json:"valor_global"`
	ModeloFormulario          *int     `json:"modelo_formulario"`
	RelatoriosExigidos        []string `json:"relatorios_exigidos"`
	TituloMinimoElegibilidade *string  `json:"titulo_minimo_elegibilidade"`
	ExigeEmpresa              *bool    `json:"exige_empresa"`
	PorteEmpresa              []string `json:"porte_empresa"`
	EnquadramentoEmpresa      []string `json:"enquadramento_empresa"`
	DocumentosObrigatorios    []string `json:"documentos_obrigatorios"`
}

type EditalSaida struct {
	ID                        int64    `json:"id"`
	Nome                      string   `json:"nome"`
	Descricao                 string   `json:"descricao"`
	DataInicio                string   `json:"data_inicio"`
	DataFim                   string   `json:"data_fim"`
	Status                    string   `json:"status"`
	TipoChamada               string   `json:"tipo_chamada"`
	NotaDeCorte               int      `json:"nota_de_corte"`
	ValorGlobal               int64    `json:"valor_global"`
	ModeloFormulario          int      `json:"modelo_formulario"`
	RelatoriosExigidos        []string `json:"relatorios_exigidos"`
	TituloMinimoElegibilidade string   `json:"titulo_minimo_elegibilidade"`
	ExigeEmpresa              bool     `json:"exige_empresa"`
	PorteEmpresa              []string `json:"porte_empresa"`
	EnquadramentoEmpresa      []string `json:"enquadramento_empresa"`
	DocumentosObrigatorios    []string `json:"documentos_obrigatorios"`
	CriadoEm                  string   `json:"criado_em"`
}

type ListarEditaisSaida struct {
	Editais []EditalSaida `json:"editais"`
	Total   int           `json:"total"`
}
