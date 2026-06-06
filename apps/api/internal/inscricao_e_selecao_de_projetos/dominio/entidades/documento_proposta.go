package entidades

import (
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
)

type DocumentoProposta struct {
	Tipo        objetos_de_valor.TipoDocumento
	NomeArquivo string
	Caminho     string
	DataUpload  time.Time
}
