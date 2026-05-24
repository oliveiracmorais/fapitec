package casos_de_uso

import (
	"context"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/repositorios"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/servicos"
)

type SolicitarRedefinicaoSenha struct {
	repo            repositorios.RepositorioDeUsuario
	tokenRepo       repositorios.RepositorioDeTokenRedefinicao
	emailService    servicos.ServicoDeEmail
	auditoria       RegistradorAuditoria
}
func NovoSolicitarRedefinicaoSenha(repo repositorios.RepositorioDeUsuario, tokenRepo repositorios.RepositorioDeTokenRedefinicao, emailService servicos.ServicoDeEmail) SolicitarRedefinicaoSenha {
	return NovoSolicitarRedefinicaoSenhaComAuditoria(repo, tokenRepo, emailService, nil)
}

func NovoSolicitarRedefinicaoSenhaComAuditoria(repo repositorios.RepositorioDeUsuario, tokenRepo repositorios.RepositorioDeTokenRedefinicao, emailService servicos.ServicoDeEmail, auditoria RegistradorAuditoria) SolicitarRedefinicaoSenha {
	return SolicitarRedefinicaoSenha{
		repo:         repo,
		tokenRepo:    tokenRepo,
		emailService: emailService,
		auditoria:    auditoria,
	}
}

func (s SolicitarRedefinicaoSenha) Executar(ctx context.Context, entrada dto.SolicitarRedefinicaoSenhaEntrada) error {
	mensagemGenerica := "Se o e-mail estiver cadastrado, voce recebera um link de redefinicao de senha."

	email, err := objetos_de_valor.NovoEmail(entrada.Email)
	if err != nil {
		return nil
	}

	usuario, _ := s.repo.BuscarPorEmail(ctx, email.String())
	if usuario == nil {
		return nil
	}

	token := objetos_de_valor.NovoTokeRedefinicao(usuario.ID, 1*time.Hour)

	if err := s.tokenRepo.Inserir(ctx, token); err != nil {
		return err
	}

	if err := s.emailService.EnviarRedefinicaoSenha(ctx, email.String(), token.Token); err != nil {
		return err
	}

	if s.auditoria != nil {
		s.auditoria.Registrar(ctx, RegistrarEventoInput{
			Acao: "solicitacao_redefinicao_senha", AtorID: usuario.ID, AtorNome: usuario.Nome,
			AtorCPF: usuario.CPF, Resultado: "sucesso",
		})
	}

	_ = mensagemGenerica
	return nil
}
