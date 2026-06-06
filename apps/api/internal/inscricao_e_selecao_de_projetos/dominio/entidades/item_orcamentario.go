package entidades

import "github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"

type ItemOrcamentario struct {
	Descricao     string
	Tipo          objetos_de_valor.TipoItemOrcamentario
	Quantidade    int
	ValorUnitario int64
	ValorTotal    int64
}
