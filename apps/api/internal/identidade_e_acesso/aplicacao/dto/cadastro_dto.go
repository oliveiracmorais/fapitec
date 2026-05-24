package dto

type CadastrarUsuarioEntrada struct {
	Nome             string
	CPF              string
	Email            string
	ConfirmacaoEmail string
	Senha            string
	ConfirmacaoSenha string
	Estrangeiro      bool
}

type CadastrarUsuarioSaida struct {
	ID    int64
	Nome  string
	CPF   string
	Email string
}
