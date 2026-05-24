package casos_de_uso

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/repositorios"
)

type ServicoDeComparacaoHash interface {
	Comparar(hash, senha string) bool
}

type AutenticarUsuario struct {
	repo          repositorios.RepositorioDeUsuario
	hash          ServicoDeComparacaoHash
	tempoBloqueio time.Duration
	auditoria     RegistradorAuditoria
}

func NovoAutenticarUsuario(repo repositorios.RepositorioDeUsuario, hash ServicoDeComparacaoHash) AutenticarUsuario {
	return NovoAutenticarUsuarioComAuditoria(repo, hash, nil)
}

func NovoAutenticarUsuarioComAuditoria(repo repositorios.RepositorioDeUsuario, hash ServicoDeComparacaoHash, auditoria RegistradorAuditoria) AutenticarUsuario {
	return AutenticarUsuario{
		repo:          repo,
		hash:          hash,
		tempoBloqueio: 15 * time.Minute,
		auditoria:     auditoria,
	}
}

func (a AutenticarUsuario) Executar(ctx context.Context, entrada dto.AutenticarUsuarioEntrada) (*dto.AutenticarUsuarioSaida, error) {
	cpf := regexp.MustCompile(`\D`).ReplaceAllString(entrada.CPF, "")
	usuario, err := a.repo.BuscarPorCPF(ctx, cpf)
	if err != nil {
		return nil, errors.New("CPF/Passaporte ou Senha inválidos. Tente novamente.")
	}
	if usuario == nil {
		return nil, errors.New("CPF/Passaporte ou Senha inválidos. Tente novamente.")
	}

	if usuario.EstaBloqueado() {
		if a.auditoria != nil {
			a.auditoria.Registrar(ctx, RegistrarEventoInput{
				Acao: "bloqueio_de_conta", AtorID: usuario.ID, AtorNome: usuario.Nome,
				AtorCPF: usuario.CPF, Resultado: "negado",
			})
		}
		return nil, errors.New("conta temporariamente bloqueada devido a multiplas tentativas falhas")
	}

	if !a.hash.Comparar(usuario.SenhaHash.String(), entrada.Senha) {
		usuario.RegistrarTentativaFalha(a.tempoBloqueio)
		var bloqueadoAte *string
		if usuario.BloqueadoAte != nil {
			b := usuario.BloqueadoAte.Format(time.RFC3339)
			bloqueadoAte = &b
		}
		_ = a.repo.AtualizarTentativas(ctx, usuario.ID, usuario.Tentativas, bloqueadoAte)
		if a.auditoria != nil {
			a.auditoria.Registrar(ctx, RegistrarEventoInput{
				Acao: "falha_de_login", AtorID: usuario.ID, AtorNome: usuario.Nome,
				AtorCPF: usuario.CPF, Resultado: "falha",
			})
		}
		if usuario.EstaBloqueado() {
			if a.auditoria != nil {
				a.auditoria.Registrar(ctx, RegistrarEventoInput{
					Acao: "bloqueio_de_conta", AtorID: usuario.ID, AtorNome: usuario.Nome,
					AtorCPF: usuario.CPF, Resultado: "negado",
				})
			}
		}
		return nil, errors.New("CPF/Passaporte ou Senha inválidos. Tente novamente.")
	}

	usuario.ResetarTentativas()
	_ = a.repo.AtualizarTentativas(ctx, usuario.ID, usuario.Tentativas, nil)

	if a.auditoria != nil {
		a.auditoria.Registrar(ctx, RegistrarEventoInput{
			Acao: "login", AtorID: usuario.ID, AtorNome: usuario.Nome,
			AtorCPF: usuario.CPF, Resultado: "sucesso",
		})
	}

	return &dto.AutenticarUsuarioSaida{
		ID:          usuario.ID,
		Nome:        usuario.Nome,
		Documento:   usuario.CPF,
		Email:       usuario.Email.String(),
		Estrangeiro: usuario.Estrangeiro,
	}, nil
}
