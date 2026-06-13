package casos_de_uso

import (
	"context"
	"testing"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func setupListarPropostasParaAvaliarTest(t *testing.T) (*ListarPropostasParaAvaliar, int64) {
	propostaRepo := persistencia.NovoRepositorioDePropostaMemoria()
	atribuicaoRepo := persistencia.NovoRepositorioDeAtribuicaoMemoria()

	avaliadorRepo := persistencia.NovoRepositorioDeAvaliadorMemoria()
	avaliador, _ := entidades.NovoAvaliador(entidades.NovoAvaliadorParams{
		UsuarioID: 1, Nome: "Avaliador",
		CPF: "000.000.000-01", Email: "a@teste.com",
	})
	avaliadorRepo.Criar(context.Background(), avaliador)

	proposta, _ := entidades.NovaProposta(entidades.NovaPropostaParams{
		EditalID: 1, ProponenteID: 1,
		DadosProponente: entidades.ProponenteInfo{Nome: "Proponente", CPF: "111.111.111-11"},
	})
	propostaRepo.Criar(context.Background(), proposta)
	proposta.Submeter()
	propostaRepo.Atualizar(context.Background(), proposta)

	atr, _ := entidades.NovaAtribuicao(entidades.NovaAtribuicaoParams{
		AvaliadorID: avaliador.ID, EditalID: 1,
		DataInicio: time.Now().Add(-24 * time.Hour),
		DataFim:    time.Now().Add(24 * time.Hour),
	})
	atribuicaoRepo.Criar(context.Background(), atr)
	atr.AceitarConvite()
	atribuicaoRepo.Atualizar(context.Background(), atr)

	uc := NovoListarPropostasParaAvaliar(propostaRepo, atribuicaoRepo)
	return uc, avaliador.ID
}

func TestListarPropostasParaAvaliar(t *testing.T) {
	t.Run("deve listar propostas para avaliador", func(t *testing.T) {
		uc, avaliadorID := setupListarPropostasParaAvaliarTest(t)

		resultado, err := uc.Executar(context.Background(), avaliadorID)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(resultado) == 0 {
			t.Fatal("esperava pelo menos uma proposta")
		}
		if resultado[0].Status != "submetida" {
			t.Errorf("status deve ser submetida, obteve %s", resultado[0].Status)
		}
	})

	t.Run("deve retornar lista vazia para avaliador sem atribuicoes", func(t *testing.T) {
		uc, _ := setupListarPropostasParaAvaliarTest(t)

		resultado, err := uc.Executar(context.Background(), 999)
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if len(resultado) != 0 {
			t.Errorf("esperava 0 propostas, obteve %d", len(resultado))
		}
	})
}
