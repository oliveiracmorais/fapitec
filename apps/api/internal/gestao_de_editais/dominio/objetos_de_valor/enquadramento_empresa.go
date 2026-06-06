package objetos_de_valor

type EnquadramentoEmpresa string

const (
	EnquadramentoSimplesNacional EnquadramentoEmpresa = "Simples Nacional"
	EnquadramentoLucroPresumido  EnquadramentoEmpresa = "Lucro Presumido"
)

func EnquadramentosEmpresaValidos() []EnquadramentoEmpresa {
	return []EnquadramentoEmpresa{EnquadramentoSimplesNacional, EnquadramentoLucroPresumido}
}

func EnquadramentoEmpresaValido(e string) bool {
	for _, v := range EnquadramentosEmpresaValidos() {
		if string(v) == e {
			return true
		}
	}
	return false
}
