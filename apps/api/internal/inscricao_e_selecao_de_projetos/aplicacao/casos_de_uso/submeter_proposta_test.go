package casos_de_uso

import (
	"context"
	"testing"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

type editalVerificadorMock struct {
	info *repositorios.EditalInfo
	err  error
}

func (m *editalVerificadorMock) BuscarEditalInfo(ctx context.Context, id int64) (*repositorios.EditalInfo, error) {
	return m.info, m.err
}

func editalAtivo() *repositorios.EditalInfo {
	agora := time.Now()
	return &repositorios.EditalInfo{
		Status:     "ativo",
		DataInicio: agora.Add(-24 * time.Hour),
		DataFim:    agora.Add(24 * time.Hour),
	}
}

func TestSubmeterProposta(t *testing.T) {
	t.Run("deve submeter proposta existente", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		verificador := &editalVerificadorMock{info: editalAtivo()}
		submeter := NovoSubmeterProposta(repo, verificador)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criada, _ := criar.Executar(context.Background(), entrada)

		saida, err := submeter.Executar(context.Background(), criada.ID)
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if saida.Status != "submetida" {
			t.Errorf("status deve ser submetida, obteve %s", saida.Status)
		}
		if saida.Protocolo == "" {
			t.Error("protocolo nao deve ser vazio apos submissao")
		}
	})

	t.Run("deve rejeitar submeter proposta inexistente", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		verificador := &editalVerificadorMock{info: editalAtivo()}
		submeter := NovoSubmeterProposta(repo, verificador)

		_, err := submeter.Executar(context.Background(), 999)
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})

	t.Run("deve rejeitar submeter quando edital inativo", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		verificador := &editalVerificadorMock{
			info: &repositorios.EditalInfo{
				Status:     "inativo",
				DataInicio: time.Now().Add(-24 * time.Hour),
				DataFim:    time.Now().Add(24 * time.Hour),
			},
		}
		submeter := NovoSubmeterProposta(repo, verificador)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criada, _ := criar.Executar(context.Background(), entrada)

		_, err := submeter.Executar(context.Background(), criada.ID)
		if err == nil {
			t.Error("esperava erro para edital inativo, obteve nil")
		}
	})

	t.Run("deve rejeitar submeter quando edital fora do prazo", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		verificador := &editalVerificadorMock{
			info: &repositorios.EditalInfo{
				Status:     "ativo",
				DataInicio: time.Now().Add(-48 * time.Hour),
				DataFim:    time.Now().Add(-24 * time.Hour),
			},
		}
		submeter := NovoSubmeterProposta(repo, verificador)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criada, _ := criar.Executar(context.Background(), entrada)

		_, err := submeter.Executar(context.Background(), criada.ID)
		if err == nil {
			t.Error("esperava erro para edital fora do prazo, obteve nil")
		}
	})

	t.Run("deve rejeitar submeter quando edital nao encontrado", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		verificador := &editalVerificadorMock{info: nil}
		submeter := NovoSubmeterProposta(repo, verificador)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criada, _ := criar.Executar(context.Background(), entrada)

		_, err := submeter.Executar(context.Background(), criada.ID)
		if err == nil {
			t.Error("esperava erro para edital inexistente, obteve nil")
		}
	})
}

func TestListarPropostas(t *testing.T) {
	t.Run("deve listar propostas vazias", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		uc := NovoListarPropostas(repo)

		saida, err := uc.Executar(context.Background(), repositorios.FiltrosListarPropostas{})
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if len(saida) != 0 {
			t.Errorf("esperava 0 propostas, obteve %d", len(saida))
		}
	})

	t.Run("deve listar propostas criadas", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		listar := NovoListarPropostas(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criar.Executar(context.Background(), entrada)

		saida, err := listar.Executar(context.Background(), repositorios.FiltrosListarPropostas{})
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if len(saida) != 1 {
			t.Errorf("esperava 1 proposta, obteve %d", len(saida))
		}
	})
}

func TestVisualizarProposta(t *testing.T) {
	t.Run("deve visualizar proposta existente", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		visualizar := NovoVisualizarProposta(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criada, _ := criar.Executar(context.Background(), entrada)

		saida, err := visualizar.Executar(context.Background(), criada.ID)
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
		if saida.DadosProponente.Nome != "Joao Silva" {
			t.Errorf("nome deve ser Joao Silva, obteve %s", saida.DadosProponente.Nome)
		}
	})

	t.Run("deve retornar erro para proposta inexistente", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		visualizar := NovoVisualizarProposta(repo)

		_, err := visualizar.Executar(context.Background(), 999)
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})
}

func TestDeletarProposta(t *testing.T) {
	t.Run("deve deletar proposta em rascunho", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		criar := NovoCriarProposta(repo)
		deletar := NovoDeletarProposta(repo)

		entrada := dto.CriarPropostaEntrada{
			EditalID:     1,
			ProponenteID: 1,
			DadosProponente: dto.ProponenteInfoDTO{
				Nome: "Joao Silva",
				CPF:  "123.456.789-00",
			},
		}

		criada, _ := criar.Executar(context.Background(), entrada)

		err := deletar.Executar(context.Background(), criada.ID)
		if err != nil {
			t.Errorf("esperava nil, obteve %v", err)
		}
	})

	t.Run("deve rejeitar deletar proposta inexistente", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDePropostaMemoria()
		deletar := NovoDeletarProposta(repo)

		err := deletar.Executar(context.Background(), 999)
		if err == nil {
			t.Error("esperava erro, obteve nil")
		}
	})
}
