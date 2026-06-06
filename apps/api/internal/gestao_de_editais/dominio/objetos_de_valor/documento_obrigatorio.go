package objetos_de_valor

type DocumentoObrigatorio string

const (
	DocumentoRG                DocumentoObrigatorio = "RG"
	DocumentoCPF               DocumentoObrigatorio = "CPF"
	DocumentoComprovanteConclusao DocumentoObrigatorio = "comprovante_conclusao"
	DocumentoComprovanteResidencia DocumentoObrigatorio = "comprovante_residencia"
)

func DocumentosObrigatoriosValidos() []DocumentoObrigatorio {
	return []DocumentoObrigatorio{DocumentoRG, DocumentoCPF, DocumentoComprovanteConclusao, DocumentoComprovanteResidencia}
}

func DocumentoObrigatorioValido(d string) bool {
	for _, v := range DocumentosObrigatoriosValidos() {
		if string(v) == d {
			return true
		}
	}
	return false
}
