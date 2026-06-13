package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type ListarPropostasParaAvaliar struct {
	propostaRepo   repositorios.RepositorioDeProposta
	atribuicaoRepo repositorios.RepositorioDeAtribuicao
}

func NovoListarPropostasParaAvaliar(propostaRepo repositorios.RepositorioDeProposta, atribuicaoRepo repositorios.RepositorioDeAtribuicao) *ListarPropostasParaAvaliar {
	return &ListarPropostasParaAvaliar{
		propostaRepo:   propostaRepo,
		atribuicaoRepo: atribuicaoRepo,
	}
}

func (uc *ListarPropostasParaAvaliar) Executar(ctx context.Context, avaliadorID int64) ([]*dto.PropostaParaAvaliarSaida, error) {
	atribuicoes, err := uc.atribuicaoRepo.ListarPorAvaliador(ctx, avaliadorID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar atribuicoes: %w", err)
	}

	var saida []*dto.PropostaParaAvaliarSaida
	editaisVistos := make(map[int64]bool)

	for _, atr := range atribuicoes {
		if !atr.ConviteFoiAceito() {
			continue
		}
		if editaisVistos[atr.EditalID] {
			continue
		}
		editaisVistos[atr.EditalID] = true

		propostas, err := uc.propostaRepo.ListarPropostasPorEdital(ctx, atr.EditalID)
		if err != nil {
			continue
		}

		for _, proposta := range propostas {
			if proposta.Status.String() != "submetida" && proposta.Status.String() != "em_avaliacao" {
				continue
			}

			pareceresSaida := make([]dto.ParecerSaida, len(proposta.Pareceres))
			for i, p := range proposta.Pareceres {
				pareceresSaida[i] = dto.ParecerSaida{
					ID:           p.ID,
					PropostaID:   p.PropostaID,
					Etapa:        p.Etapa,
					Nota:         p.Nota,
					ParecerTexto: p.ParecerTexto,
					Data:         p.Data.Format("2006-01-02T15:04:05Z07:00"),
				}
			}

			saida = append(saida, &dto.PropostaParaAvaliarSaida{
				ID:       proposta.ID,
				EditalID: proposta.EditalID,
				Protocolo: proposta.Protocolo.String(),
				Status:   proposta.Status.String(),
				DadosProponente: dto.ProponenteInfoDTO{
					Nome:          proposta.DadosProponente.Nome,
					CPF:           proposta.DadosProponente.CPF,
					Email:         proposta.DadosProponente.Email,
				},
				DadosAcademicos: dto.DadosAcademicosDTO{
					MaiorTitulacao:  proposta.DadosAcademicos.MaiorTitulacao,
					AreaConhecimento: proposta.DadosAcademicos.AreaConhecimento,
				},
				ValorTotal:    proposta.ValorTotalSolicitado,
				Pareceres:     pareceresSaida,
				DataSubmissao: proposta.DataSubmissao.Format("2006-01-02T15:04:05Z07:00"),
			})
		}
	}

	if saida == nil {
		saida = []*dto.PropostaParaAvaliarSaida{}
	}

	return saida, nil
}
