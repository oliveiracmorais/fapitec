package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type CadastrarAvaliador struct {
	repo repositorios.RepositorioDeAvaliador
}

func NovoCadastrarAvaliador(repo repositorios.RepositorioDeAvaliador) *CadastrarAvaliador {
	return &CadastrarAvaliador{repo: repo}
}

func (uc *CadastrarAvaliador) Executar(ctx context.Context, entrada dto.CadastrarAvaliadorEntrada) (*dto.AvaliadorSaida, error) {
	existente, _ := uc.repo.BuscarPorCPF(ctx, entrada.CPF)
	if existente != nil {
		return nil, fmt.Errorf("ja existe um avaliador com este CPF")
	}

	avaliador, err := entidades.NovoAvaliador(entidades.NovoAvaliadorParams{
		UsuarioID:         entrada.UsuarioID,
		Nome:              entrada.Nome,
		CPF:               entrada.CPF,
		Email:             entrada.Email,
		TitulacaoMaxima:   entrada.TitulacaoMaxima,
		AreaConhecimento:  entrada.AreaConhecimento,
		Instituicao:       entrada.Instituicao,
		CurriculoResumido: entrada.CurriculoResumido,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar avaliador: %w", err)
	}

	if err := uc.repo.Criar(ctx, avaliador); err != nil {
		return nil, fmt.Errorf("erro ao salvar avaliador: %w", err)
	}

	return paraAvaliadorSaida(avaliador, 0, 0), nil
}
