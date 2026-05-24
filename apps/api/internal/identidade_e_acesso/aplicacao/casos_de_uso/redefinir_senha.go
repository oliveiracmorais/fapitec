package casos_de_uso

import (
	"context"
	"errors"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/repositorios"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/servicos"
)

type ServicoDeHashRedefinicao interface {
	Hash(senha string) (string, error)
}

type RedefinirSenha struct {
	repo       repositorios.RepositorioDeUsuario
	tokenRepo  repositorios.RepositorioDeTokenRedefinicao
	validador  servicos.ValidadorDeSenha
	hash       ServicoDeHashRedefinicao
	auditoria  RegistradorAuditoria
}

func NovoRedefinirSenha(repo repositorios.RepositorioDeUsuario, tokenRepo repositorios.RepositorioDeTokenRedefinicao, hash ServicoDeHashRedefinicao) RedefinirSenha {
	return NovoRedefinirSenhaComAuditoria(repo, tokenRepo, hash, nil)
}

func NovoRedefinirSenhaComAuditoria(repo repositorios.RepositorioDeUsuario, tokenRepo repositorios.RepositorioDeTokenRedefinicao, hash ServicoDeHashRedefinicao, auditoria RegistradorAuditoria) RedefinirSenha {
	return RedefinirSenha{
		repo:      repo,
		tokenRepo: tokenRepo,
		validador: servicos.NovoValidadorDeSenha(),
		hash:      hash,
		auditoria: auditoria,
	}
}

func (r RedefinirSenha) Executar(ctx context.Context, entrada dto.RedefinirSenhaEntrada) error {
	token, err := r.tokenRepo.BuscarPorToken(ctx, entrada.Token)
	if err != nil {
		return errors.New("token invalido ou expirado")
	}
	if token == nil {
		return errors.New("token invalido ou expirado")
	}

	if token.Expirado() {
		r.tokenRepo.Remover(ctx, entrada.Token)
		return errors.New("token invalido ou expirado")
	}

	if err := r.validador.Validar(entrada.Senha); err != nil {
		return err
	}

	if entrada.Senha != entrada.ConfirmacaoSenha {
		return errors.New("A senha deve ser IGUAL a primeira senha fornecida.")
	}

	hashSenha, err := r.hash.Hash(entrada.Senha)
	if err != nil {
		return errors.New("erro ao processar senha")
	}

	if err := r.repo.AtualizarSenha(ctx, token.UsuarioID, hashSenha); err != nil {
		return err
	}

	r.tokenRepo.Remover(ctx, entrada.Token)

	if r.auditoria != nil {
		r.auditoria.Registrar(ctx, RegistrarEventoInput{
			Acao: "redefinicao_de_senha", AtorID: token.UsuarioID,
			Resultado: "sucesso",
		})
	}

	return nil
}
