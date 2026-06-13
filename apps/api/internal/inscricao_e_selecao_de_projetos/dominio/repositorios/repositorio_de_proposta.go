package repositorios

import (
	"context"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
)

type FiltrosListarPropostas struct {
	ProponenteID int64
	EditalID     int64
	Status       string
}

type RepositorioDeProposta interface {
	Criar(ctx context.Context, proposta *entidades.Proposta) error
	BuscarPorID(ctx context.Context, id int64) (*entidades.Proposta, error)
	BuscarPorProtocolo(ctx context.Context, protocolo string) (*entidades.Proposta, error)
	Listar(ctx context.Context, filtros FiltrosListarPropostas) ([]*entidades.Proposta, error)
	Atualizar(ctx context.Context, proposta *entidades.Proposta) error
	Deletar(ctx context.Context, id int64) error
	ContarPorEdital(ctx context.Context, editalID int64) (int64, error)
	SalvarParecer(ctx context.Context, parecer *entidades.Parecer) error
	ListarPareceresPorProposta(ctx context.Context, propostaID int64) ([]*entidades.Parecer, error)
	ListarPareceresPorAvaliador(ctx context.Context, avaliadorID int64) ([]*entidades.Parecer, error)
	ListarPropostasPorEdital(ctx context.Context, editalID int64) ([]*entidades.Proposta, error)
	AtualizarStatus(ctx context.Context, id int64, status string) error
}

type EditalInfo struct {
	Status     string
	DataInicio time.Time
	DataFim    time.Time
}

type EditalVerificador interface {
	BuscarEditalInfo(ctx context.Context, id int64) (*EditalInfo, error)
}
