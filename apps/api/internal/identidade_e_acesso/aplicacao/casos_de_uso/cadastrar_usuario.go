package casos_de_uso

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/repositorios"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/servicos"
)

var nomeRegex = regexp.MustCompile(`^[A-Za-zÀ-ÿ\s]{3,}$`)

type ServicoDeHash interface {
	Hash(senha string) (string, error)
}

type CadastrarUsuario struct {
	repo       repositorios.RepositorioDeUsuario
	validador  servicos.ValidadorDeSenha
	hash       ServicoDeHash
	auditoria  RegistradorAuditoria
}

func NovoCadastrarUsuario(repo repositorios.RepositorioDeUsuario, hash ServicoDeHash) CadastrarUsuario {
	return NovoCadastrarUsuarioComAuditoria(repo, hash, nil)
}

func NovoCadastrarUsuarioComAuditoria(repo repositorios.RepositorioDeUsuario, hash ServicoDeHash, auditoria RegistradorAuditoria) CadastrarUsuario {
	return CadastrarUsuario{
		repo:      repo,
		validador: servicos.NovoValidadorDeSenha(),
		hash:      hash,
		auditoria: auditoria,
	}
}

func (c CadastrarUsuario) Executar(ctx context.Context, entrada dto.CadastrarUsuarioEntrada) (*dto.CadastrarUsuarioSaida, error) {
	nome := strings.TrimSpace(entrada.Nome)
	if !nomeRegex.MatchString(nome) {
		return nil, errors.New("Nome inválido. Deve conter apenas letras e espaços, mínimo 3 caracteres.")
	}

	var identificador string
	if entrada.Estrangeiro {
		passaporte, err := objetos_de_valor.NovoPassaporte(entrada.CPF)
		if err != nil {
			return nil, err
		}
		identificador = passaporte.String()
	} else {
		cpf, err := objetos_de_valor.NovoCPF(entrada.CPF)
		if err != nil {
			return nil, errors.New("CPF inválido. Verifique os dígitos.")
		}
		identificador = cpf.String()
	}

	email, err := objetos_de_valor.NovoEmail(entrada.Email)
	if err != nil {
		return nil, errors.New("E-mail inválido. Verifique o formato.")
	}

	if entrada.Email != entrada.ConfirmacaoEmail {
		return nil, errors.New("O e-mail deve ser IGUAL ao e-mail principal.")
	}

	if err := c.validador.Validar(entrada.Senha); err != nil {
		return nil, err
	}

	if entrada.Senha != entrada.ConfirmacaoSenha {
		return nil, errors.New("A senha deve ser IGUAL à primeira senha fornecida.")
	}

	existente, _ := c.repo.BuscarPorCPF(ctx, identificador)
	if existente != nil {
		if entrada.Estrangeiro {
			return nil, errors.New("Passaporte já cadastrado no sistema.")
		}
		return nil, errors.New("CPF já cadastrado no sistema.")
	}

	existente, _ = c.repo.BuscarPorEmail(ctx, email.String())
	if existente != nil {
		return nil, errors.New("E-mail já cadastrado. Utilize outro endereço ou recupere sua senha.")
	}

	hashSenha, err := c.hash.Hash(entrada.Senha)
	if err != nil {
		return nil, errors.New("erro ao processar senha")
	}

	usuario := &entidades.Usuario{
		Nome:        nome,
		CPF:         identificador,
		Estrangeiro: entrada.Estrangeiro,
		Email:       email,
		SenhaHash:   objetos_de_valor.NovaSenhaHash(hashSenha),
	}

	if err := c.repo.Inserir(ctx, usuario); err != nil {
		return nil, err
	}

	if c.auditoria != nil {
		c.auditoria.Registrar(ctx, RegistrarEventoInput{
			Acao: "cadastro_de_usuario", AtorID: usuario.ID, AtorNome: usuario.Nome,
			AtorCPF: usuario.CPF, Resultado: "sucesso",
		})
	}

	return &dto.CadastrarUsuarioSaida{
		ID:          usuario.ID,
		Nome:        usuario.Nome,
		Documento:   usuario.CPF,
		Email:       email.String(),
		Estrangeiro: usuario.Estrangeiro,
	}, nil
}
