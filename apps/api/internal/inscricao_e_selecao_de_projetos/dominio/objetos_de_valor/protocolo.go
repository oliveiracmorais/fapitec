package objetos_de_valor

import (
	"fmt"
	"time"
)

type Protocolo struct {
	valor string
}

func NovoProtocolo(editalID int64, sequencial int) Protocolo {
	ano := time.Now().Year()
	valor := fmt.Sprintf("%d.%d.%03d", ano, editalID, sequencial)
	return Protocolo{valor: valor}
}

func ProtocoloExistente(valor string) (Protocolo, error) {
	if valor == "" {
		return Protocolo{}, fmt.Errorf("protocolo nao pode ser vazio")
	}
	return Protocolo{valor: valor}, nil
}

func (p Protocolo) String() string {
	return p.valor
}

func (p Protocolo) IsZero() bool {
	return p.valor == ""
}
