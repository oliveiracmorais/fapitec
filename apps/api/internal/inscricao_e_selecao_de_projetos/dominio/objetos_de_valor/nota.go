package objetos_de_valor

import "fmt"

type Nota int

const (
	NotaMinima Nota = 0
	NotaMaxima Nota = 100
)

func NovaNota(valor int) (Nota, error) {
	if valor < int(NotaMinima) || valor > int(NotaMaxima) {
		return 0, fmt.Errorf("nota deve estar entre %d e %d", NotaMinima, NotaMaxima)
	}
	return Nota(valor), nil
}

func (n Nota) Valor() int {
	return int(n)
}

func (n Nota) String() string {
	return fmt.Sprintf("%d", n)
}
