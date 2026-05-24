package dto

type SolicitarRedefinicaoSenhaEntrada struct {
	Email string
}

type RedefinirSenhaEntrada struct {
	Token            string
	Senha            string
	ConfirmacaoSenha string
}
