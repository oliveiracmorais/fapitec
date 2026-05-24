package email

import (
	"context"
	"log"
)

type ServicoDeEmailLog struct{}

func NovoServicoDeEmailLog() *ServicoDeEmailLog {
	return &ServicoDeEmailLog{}
}

func (s *ServicoDeEmailLog) EnviarRedefinicaoSenha(ctx context.Context, email, token string) error {
	log.Printf("[EMAIL-PLACEHOLDER] Para: %s | Token de redefinicao: %s", email, token)
	return nil
}
