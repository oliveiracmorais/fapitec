package dto

type AutenticarUsuarioEntrada struct {
	CPF    string
	Senha  string
}

type AutenticarUsuarioSaida struct {
	ID    int64
	Nome  string
	CPF   string
	Email string
}
