package dto

type AutenticarUsuarioEntrada struct {
	CPF   string `json:"cpf"`
	Senha string `json:"senha"`
}

type AutenticarUsuarioSaida struct {
	ID          int64  `json:"id"`
	Nome        string `json:"nome"`
	Documento   string `json:"documento"`
	Email       string `json:"email"`
	Estrangeiro bool   `json:"estrangeiro"`
}
