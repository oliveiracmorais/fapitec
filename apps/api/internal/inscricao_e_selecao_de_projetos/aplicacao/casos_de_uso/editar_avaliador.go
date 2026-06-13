package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type EditarAvaliador struct {
	repo repositorios.RepositorioDeAvaliador
}

func NovoEditarAvaliador(repo repositorios.RepositorioDeAvaliador) *EditarAvaliador {
	return &EditarAvaliador{repo: repo}
}

func (uc *EditarAvaliador) Executar(ctx context.Context, id int64, entrada dto.AtualizarAvaliadorEntrada) (*dto.AvaliadorSaida, error) {
	avaliador, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar avaliador: %w", err)
	}
	if avaliador == nil {
		return nil, fmt.Errorf("avaliador nao encontrado")
	}

	if entrada.Nome != nil {
		avaliador.Nome = *entrada.Nome
	}
	if entrada.CPF != nil {
		avaliador.CPF = *entrada.CPF
	}
	if entrada.Email != nil {
		avaliador.Email = *entrada.Email
	}
	if entrada.TitulacaoMaxima != nil {
		avaliador.TitulacaoMaxima = *entrada.TitulacaoMaxima
	}
	if entrada.AreaConhecimento != nil {
		avaliador.AreaConhecimento = *entrada.AreaConhecimento
	}
	if entrada.Instituicao != nil {
		avaliador.Instituicao = *entrada.Instituicao
	}
	if entrada.CurriculoResumido != nil {
		avaliador.CurriculoResumido = *entrada.CurriculoResumido
	}
	if entrada.Estado != nil {
		switch *entrada.Estado {
		case "ativo":
			avaliador.Ativar()
		case "inativo":
			avaliador.Inativar()
		default:
			return nil, fmt.Errorf("estado invalido: %s", *entrada.Estado)
		}
	}

	if err := uc.repo.Atualizar(ctx, avaliador); err != nil {
		return nil, fmt.Errorf("erro ao atualizar avaliador: %w", err)
	}

	return paraAvaliadorSaida(avaliador, 0, 0), nil
}
