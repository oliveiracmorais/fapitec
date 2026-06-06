package objetos_de_valor

type PorteEmpresa string

const (
	PorteEmpresaMEI  PorteEmpresa = "MEI"
	PorteEmpresaME   PorteEmpresa = "ME"
	PorteEmpresaEPP  PorteEmpresa = "EPP"
)

func PortesEmpresaValidos() []PorteEmpresa {
	return []PorteEmpresa{PorteEmpresaMEI, PorteEmpresaME, PorteEmpresaEPP}
}

func PorteEmpresaValido(p string) bool {
	for _, v := range PortesEmpresaValidos() {
		if string(v) == p {
			return true
		}
	}
	return false
}
