package casos_de_uso

import (
	"context"
	"fmt"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type SubmeterProposta struct {
	repo            repositorios.RepositorioDeProposta
	editalVerificador repositorios.EditalVerificador
}

func NovoSubmeterProposta(repo repositorios.RepositorioDeProposta, editalVerificador repositorios.EditalVerificador) *SubmeterProposta {
	return &SubmeterProposta{repo: repo, editalVerificador: editalVerificador}
}

func (uc *SubmeterProposta) Executar(ctx context.Context, id int64) (*dto.PropostaSaida, error) {
	proposta, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar proposta: %w", err)
	}
	if proposta == nil {
		return nil, fmt.Errorf("proposta nao encontrada")
	}

	editalInfo, err := uc.editalVerificador.BuscarEditalInfo(ctx, proposta.EditalID)
	if err != nil {
		return nil, fmt.Errorf("erro ao validar edital: %w", err)
	}
	if editalInfo == nil {
		return nil, fmt.Errorf("edital nao encontrado")
	}
	if editalInfo.Status != "ativo" {
		return nil, fmt.Errorf("edital nao esta ativo")
	}
	agora := time.Now()
	if agora.Before(editalInfo.DataInicio) || agora.After(editalInfo.DataFim) {
		return nil, fmt.Errorf("edital fora do prazo de vigencia")
	}

	if err := proposta.Submeter(); err != nil {
		return nil, err
	}

	sequencial, err := uc.repo.ContarPorEdital(ctx, proposta.EditalID)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar protocolo: %w", err)
	}

	proposta.Protocolo = objetos_de_valor.NovoProtocolo(proposta.EditalID, int(sequencial)+1)

	if err := uc.repo.Atualizar(ctx, proposta); err != nil {
		return nil, fmt.Errorf("erro ao submeter proposta: %w", err)
	}

	return paraPropostaSaida(proposta), nil
}
