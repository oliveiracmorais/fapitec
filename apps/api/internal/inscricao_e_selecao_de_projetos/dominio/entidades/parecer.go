package entidades

import "time"

type Parecer struct {
	Etapa        string
	AvaliadorID  int64
	Nota         int
	ParecerTexto string
	Data         time.Time
}
