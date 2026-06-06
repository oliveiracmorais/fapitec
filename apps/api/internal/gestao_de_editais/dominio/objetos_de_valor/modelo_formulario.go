package objetos_de_valor

import "fmt"

type ModeloFormulario int

const (
	ModeloFormularioNaoDefinido ModeloFormulario = 0
	ModeloFormulario1           ModeloFormulario = 1
	ModeloFormulario2           ModeloFormulario = 2
	ModeloFormulario3           ModeloFormulario = 3
	ModeloFormulario4           ModeloFormulario = 4
	ModeloFormulario5           ModeloFormulario = 5
	ModeloFormulario6           ModeloFormulario = 6
)

func NovoModeloFormulario(v int) (ModeloFormulario, error) {
	if v < 0 || v > 6 {
		return ModeloFormularioNaoDefinido, fmt.Errorf("modelo de formulario invalido: %d", v)
	}
	return ModeloFormulario(v), nil
}
