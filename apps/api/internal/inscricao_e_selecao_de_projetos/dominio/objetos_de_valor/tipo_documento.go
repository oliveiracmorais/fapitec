package objetos_de_valor

import "fmt"

type TipoDocumento string

const (
	TipoDocumentoRG                    TipoDocumento = "RG"
	TipoDocumentoCPF                   TipoDocumento = "CPF"
	TipoDocumentoComprovanteConclusao  TipoDocumento = "comprovante_conclusao"
	TipoDocumentoComprovanteResidencia TipoDocumento = "comprovante_residencia"
	TipoDocumentoDiploma               TipoDocumento = "diploma"
	TipoDocumentoHistorico             TipoDocumento = "historico"
	TipoDocumentoPlanoTrabalho         TipoDocumento = "plano_trabalho"
	TipoDocumentoOutro                 TipoDocumento = "outro"
)

func NovoTipoDocumento(s string) (TipoDocumento, error) {
	switch TipoDocumento(s) {
	case TipoDocumentoRG, TipoDocumentoCPF,
		TipoDocumentoComprovanteConclusao, TipoDocumentoComprovanteResidencia,
		TipoDocumentoDiploma, TipoDocumentoHistorico,
		TipoDocumentoPlanoTrabalho, TipoDocumentoOutro:
		return TipoDocumento(s), nil
	default:
		return "", fmt.Errorf("tipo de documento invalido: %s", s)
	}
}

func (t TipoDocumento) String() string {
	return string(t)
}
