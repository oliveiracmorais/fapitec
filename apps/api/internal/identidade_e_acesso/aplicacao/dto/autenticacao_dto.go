package dto

type AutenticarUsuarioEntrada struct {
	CPF   string `json:"cpf"`
	Senha string `json:"senha"`
}

type AutenticarUsuarioSaida struct {
	ID    int64  `json:"id"`
	Nome  string `json:"nome"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
}
