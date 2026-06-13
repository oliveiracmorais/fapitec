package casos_de_uso

import (
	"context"
	"fmt"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type EmitirParecer struct {
	propostaRepo   repositorios.RepositorioDeProposta
	avaliadorRepo  repositorios.RepositorioDeAvaliador
	atribuicaoRepo repositorios.RepositorioDeAtribuicao
}

func NovoEmitirParecer(propostaRepo repositorios.RepositorioDeProposta, avaliadorRepo repositorios.RepositorioDeAvaliador, atribuicaoRepo repositorios.RepositorioDeAtribuicao) *EmitirParecer {
	return &EmitirParecer{
		propostaRepo:   propostaRepo,
		avaliadorRepo:  avaliadorRepo,
		atribuicaoRepo: atribuicaoRepo,
	}
}

func (uc *EmitirParecer) Executar(ctx context.Context, entrada dto.EmitirParecerEntrada) (*dto.ParecerSaida, error) {
	proposta, err := uc.propostaRepo.BuscarPorID(ctx, entrada.PropostaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar proposta: %w", err)
	}
	if proposta == nil {
		return nil, fmt.Errorf("proposta nao encontrada")
	}

	avaliador, err := uc.avaliadorRepo.BuscarPorID(ctx, entrada.AvaliadorID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar avaliador: %w", err)
	}
	if avaliador == nil {
		return nil, fmt.Errorf("avaliador nao encontrado")
	}

	if avaliador.EstaInativo() {
		return nil, fmt.Errorf("avaliador inativo nao pode emitir pareceres")
	}

	atribuicoes, err := uc.atribuicaoRepo.ListarPorAvaliador(ctx, entrada.AvaliadorID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar atribuicoes: %w", err)
	}

	atribuicaoValida := false
	for _, a := range atribuicoes {
		if a.EditalID == proposta.EditalID && a.ConviteFoiAceito() {
			atribuicaoValida = true
			break
		}
	}
	if !atribuicaoValida {
		return nil, fmt.Errorf("avaliador nao esta atribuido a este edital com convite aceito")
	}

	parecer, err := entidades.NovoParecer(entidades.NovoParecerParams{
		PropostaID:   entrada.PropostaID,
		Etapa:        entrada.Etapa,
		AvaliadorID:  entrada.AvaliadorID,
		Nota:         entrada.Nota,
		ParecerTexto: entrada.ParecerTexto,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar parecer: %w", err)
	}

	if err := uc.propostaRepo.SalvarParecer(ctx, parecer); err != nil {
		return nil, fmt.Errorf("erro ao salvar parecer: %w", err)
	}

	propostaAtualizada, err := uc.propostaRepo.BuscarPorID(ctx, entrada.PropostaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar proposta atualizada: %w", err)
	}

	if err := propostaAtualizada.AdicionarParecer(parecer); err != nil {
		return nil, fmt.Errorf("erro ao adicionar parecer a proposta: %w", err)
	}

	if err := uc.propostaRepo.Atualizar(ctx, propostaAtualizada); err != nil {
		return nil, fmt.Errorf("erro ao atualizar proposta: %w", err)
	}

	return &dto.ParecerSaida{
		ID:           parecer.ID,
		PropostaID:   parecer.PropostaID,
		Etapa:        parecer.Etapa,
		Nota:         parecer.Nota,
		ParecerTexto: parecer.ParecerTexto,
		Data:         parecer.Data.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
