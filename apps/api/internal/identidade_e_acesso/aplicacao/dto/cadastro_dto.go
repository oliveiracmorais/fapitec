package dto

type CadastrarUsuarioEntrada struct {
	Nome             string `json:"nome"`
	CPF              string `json:"cpf"`
	Email            string `json:"email"`
	ConfirmacaoEmail string `json:"confirmacao_email"`
	Senha            string `json:"senha"`
	ConfirmacaoSenha string `json:"confirmacao_senha"`
	Estrangeiro      bool   `json:"estrangeiro"`
}

type CadastrarUsuarioSaida struct {
	ID          int64  `json:"id"`
	Nome        string `json:"nome"`
	Documento   string `json:"documento"`
	Email       string `json:"email"`
	Estrangeiro bool   `json:"estrangeiro"`
}
