package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type GerenciarConvite struct {
	repo repositorios.RepositorioDeAtribuicao
}

func NovoGerenciarConvite(repo repositorios.RepositorioDeAtribuicao) *GerenciarConvite {
	return &GerenciarConvite{repo: repo}
}

func (uc *GerenciarConvite) Executar(ctx context.Context, id int64, entrada dto.GerenciarConviteEntrada) (*dto.AtribuicaoSaida, error) {
	atribuicao, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar atribuicao: %w", err)
	}
	if atribuicao == nil {
		return nil, fmt.Errorf("atribuicao nao encontrada")
	}

	switch entrada.Acao {
	case "aceitar":
		if err := atribuicao.AceitarConvite(); err != nil {
			return nil, err
		}
	case "recusar":
		if err := atribuicao.RecusarConvite(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("acao invalida: use 'aceitar' ou 'recusar'")
	}

	if err := uc.repo.Atualizar(ctx, atribuicao); err != nil {
		return nil, fmt.Errorf("erro ao atualizar convite: %w", err)
	}

	return paraAtribuicaoSaida(atribuicao), nil
}
