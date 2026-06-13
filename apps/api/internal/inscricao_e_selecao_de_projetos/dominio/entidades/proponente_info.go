package entidades

import "github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"

type ProponenteInfo struct {
	Nome          string
	CPF           string
	RG            string
	Genero        string
	Etnia         string
	DataNascimento string
	Endereco      objetos_de_valor.Endereco
	Telefone      string
	Email         string
}
