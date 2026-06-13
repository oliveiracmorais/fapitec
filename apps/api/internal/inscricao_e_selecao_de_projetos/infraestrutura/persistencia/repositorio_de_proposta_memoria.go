package persistencia

import (
	"context"
	"sync"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/objetos_de_valor"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
)

type RepositorioDePropostaMemoria struct {
	mu           sync.RWMutex
	propostas    map[int64]*entidades.Proposta
	versoes      map[int64][]*entidades.Proposta
	pareceres    map[int64]*entidades.Parecer
	proximoID    int64
	proximoParecerID int64
}

func NovoRepositorioDePropostaMemoria() *RepositorioDePropostaMemoria {
	return &RepositorioDePropostaMemoria{
		propostas:    make(map[int64]*entidades.Proposta),
		versoes:      make(map[int64][]*entidades.Proposta),
		pareceres:    make(map[int64]*entidades.Parecer),
		proximoID:    1,
		proximoParecerID: 1,
	}
}

func (r *RepositorioDePropostaMemoria) Criar(ctx context.Context, proposta *entidades.Proposta) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	proposta.ID = r.proximoID
	r.proximoID++
	r.propostas[proposta.ID] = proposta
	return nil
}

func (r *RepositorioDePropostaMemoria) BuscarPorID(ctx context.Context, id int64) (*entidades.Proposta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proposta, ok := r.propostas[id]
	if !ok {
		return nil, nil
	}
	return proposta, nil
}

func (r *RepositorioDePropostaMemoria) BuscarPorProtocolo(ctx context.Context, protocolo string) (*entidades.Proposta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.propostas {
		if p.Protocolo.String() == protocolo {
			return p, nil
		}
	}
	return nil, nil
}

func (r *RepositorioDePropostaMemoria) Listar(ctx context.Context, filtros repositorios.FiltrosListarPropostas) ([]*entidades.Proposta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resultado := make([]*entidades.Proposta, 0, len(r.propostas))
	for _, p := range r.propostas {
		if filtros.EditalID != 0 && p.EditalID != filtros.EditalID {
			continue
		}
		if filtros.ProponenteID != 0 && p.ProponenteID != filtros.ProponenteID {
			continue
		}
		if filtros.Status != "" && p.Status.String() != filtros.Status {
			continue
		}
		resultado = append(resultado, p)
	}
	return resultado, nil
}

func (r *RepositorioDePropostaMemoria) Atualizar(ctx context.Context, proposta *entidades.Proposta) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existente, ok := r.propostas[proposta.ID]
	if ok {
		versaoAnterior := copiarProposta(existente)
		r.versoes[proposta.ID] = append(r.versoes[proposta.ID], versaoAnterior)
	}

	proposta.Versao = len(r.versoes[proposta.ID]) + 1
	r.propostas[proposta.ID] = proposta
	return nil
}

func copiarProposta(p *entidades.Proposta) *entidades.Proposta {
	if p == nil {
		return nil
	}
	itens := make([]entidades.ItemOrcamentario, len(p.ItensOrcamentarios))
	copy(itens, p.ItensOrcamentarios)
	documentos := make([]entidades.DocumentoProposta, len(p.Documentos))
	copy(documentos, p.Documentos)
	pareceres := make([]entidades.Parecer, len(p.Pareceres))
	copy(pareceres, p.Pareceres)
	return &entidades.Proposta{
		ID:                   p.ID,
		EditalID:             p.EditalID,
		ProponenteID:         p.ProponenteID,
		Protocolo:            p.Protocolo,
		Status:               p.Status,
		Versao:               p.Versao,
		DadosProponente:      p.DadosProponente,
		DadosAcademicos:      p.DadosAcademicos,
		EmpresaVinculada:     p.EmpresaVinculada,
		ItensOrcamentarios:   itens,
		Documentos:           documentos,
		ValorTotalSolicitado: p.ValorTotalSolicitado,
		Pareceres:            pareceres,
		DataSubmissao:        p.DataSubmissao,
		DataAtualizacao:      p.DataAtualizacao,
		CriadoEm:             p.CriadoEm,
	}
}

func (r *RepositorioDePropostaMemoria) Deletar(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.propostas, id)
	return nil
}

func (r *RepositorioDePropostaMemoria) ContarPorEdital(ctx context.Context, editalID int64) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int64
	for _, p := range r.propostas {
		if p.EditalID == editalID && p.Status.String() != "rascunho" {
			count++
		}
	}
	return count, nil
}

func (r *RepositorioDePropostaMemoria) SalvarParecer(_ context.Context, parecer *entidades.Parecer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	parecer.ID = r.proximoParecerID
	r.proximoParecerID++
	c := *parecer
	r.pareceres[parecer.ID] = &c
	return nil
}

func (r *RepositorioDePropostaMemoria) ListarPareceresPorProposta(_ context.Context, propostaID int64) ([]*entidades.Parecer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resultado []*entidades.Parecer
	for _, p := range r.pareceres {
		if p.PropostaID == propostaID {
			c := *p
			resultado = append(resultado, &c)
		}
	}
	if resultado == nil {
		resultado = []*entidades.Parecer{}
	}
	return resultado, nil
}

func (r *RepositorioDePropostaMemoria) ListarPareceresPorAvaliador(_ context.Context, avaliadorID int64) ([]*entidades.Parecer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resultado []*entidades.Parecer
	for _, p := range r.pareceres {
		if p.AvaliadorID == avaliadorID {
			c := *p
			resultado = append(resultado, &c)
		}
	}
	if resultado == nil {
		resultado = []*entidades.Parecer{}
	}
	return resultado, nil
}

func (r *RepositorioDePropostaMemoria) ListarPropostasPorEdital(_ context.Context, editalID int64) ([]*entidades.Proposta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resultado []*entidades.Proposta
	for _, p := range r.propostas {
		if p.EditalID == editalID {
			resultado = append(resultado, p)
		}
	}
	if resultado == nil {
		resultado = []*entidades.Proposta{}
	}
	return resultado, nil
}

func (r *RepositorioDePropostaMemoria) AtualizarStatus(_ context.Context, id int64, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	proposta, ok := r.propostas[id]
	if !ok {
		return nil
	}
	novoStatus, _ := objetos_de_valor.NovoStatusProposta(status)
	proposta.Status = novoStatus
	proposta.DataAtualizacao = time.Now()
	return nil
}


