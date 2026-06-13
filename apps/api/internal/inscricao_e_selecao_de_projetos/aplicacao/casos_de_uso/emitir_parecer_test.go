package casos_de_uso

import (
	"context"
	"testing"
	"time"

	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
)

func setupEmitirParecerTest(t *testing.T) (*EmitirParecer, *entidades.Proposta, int64, func()) {
	propostaRepo := persistencia.NovoRepositorioDePropostaMemoria()
	avaliadorRepo := persistencia.NovoRepositorioDeAvaliadorMemoria()
	atribuicaoRepo := persistencia.NovoRepositorioDeAtribuicaoMemoria()

	avaliador, _ := entidades.NovoAvaliador(entidades.NovoAvaliadorParams{
		UsuarioID: 1, Nome: "Avaliador Teste",
		CPF: "000.000.000-01", Email: "avaliador@teste.com",
	})
	avaliadorRepo.Criar(context.Background(), avaliador)

	proposta, _ := entidades.NovaProposta(entidades.NovaPropostaParams{
		EditalID: 1, ProponenteID: 1,
		DadosProponente: entidades.ProponenteInfo{Nome: "Proponente", CPF: "111.111.111-11"},
	})
	propostaRepo.Criar(context.Background(), proposta)
	proposta.Submeter()
	propostaRepo.Atualizar(context.Background(), proposta)

	atribuicao, _ := entidades.NovaAtribuicao(entidades.NovaAtribuicaoParams{
		AvaliadorID: avaliador.ID, EditalID: 1,
		DataInicio: time.Now().Add(-24 * time.Hour),
		DataFim:    time.Now().Add(24 * time.Hour),
	})
	atribuicaoRepo.Criar(context.Background(), atribuicao)
	atribuicao.AceitarConvite()
	atribuicaoRepo.Atualizar(context.Background(), atribuicao)

	uc := NovoEmitirParecer(propostaRepo, avaliadorRepo, atribuicaoRepo)
	return uc, proposta, avaliador.ID, func() {}
}

func TestEmitirParecer(t *testing.T) {
	t.Run("deve emitir parecer com sucesso", func(t *testing.T) {
		uc, proposta, avaliadorID, _ := setupEmitirParecerTest(t)

		saida, err := uc.Executar(context.Background(), dto.EmitirParecerEntrada{
			PropostaID:   proposta.ID,
			Etapa:        "unica",
			AvaliadorID:  avaliadorID,
			Nota:         85,
			ParecerTexto: "Projeto excelente e bem fundamentado",
		})
		if err != nil {
			t.Fatalf("esperava nil, obteve %v", err)
		}
		if saida.Nota != 85 {
			t.Errorf("nota deve ser 85, obteve %d", saida.Nota)
		}
		if saida.PropostaID != proposta.ID {
			t.Errorf("proposta ID incorreto")
		}
	})

	t.Run("deve rejeitar parecer de avaliador nao atribuido", func(t *testing.T) {
		uc, proposta, _, _ := setupEmitirParecerTest(t)

		_, err := uc.Executar(context.Background(), dto.EmitirParecerEntrada{
			PropostaID:   proposta.ID,
			Etapa:        "unica",
			AvaliadorID:  999,
			Nota:         85,
			ParecerTexto: "Teste",
		})
		if err == nil {
			t.Fatal("esperava erro para avaliador nao atribuido")
		}
	})

	t.Run("deve rejeitar nota acima de 100", func(t *testing.T) {
		uc, proposta, avaliadorID, _ := setupEmitirParecerTest(t)

		_, err := uc.Executar(context.Background(), dto.EmitirParecerEntrada{
			PropostaID:   proposta.ID,
			Etapa:        "unica",
			AvaliadorID:  avaliadorID,
			Nota:         150,
			ParecerTexto: "Nota invalida",
		})
		if err == nil {
			t.Fatal("esperava erro para nota acima de 100")
		}
	})

	t.Run("deve rejeitar parecer textual vazio", func(t *testing.T) {
		uc, proposta, avaliadorID, _ := setupEmitirParecerTest(t)

		_, err := uc.Executar(context.Background(), dto.EmitirParecerEntrada{
			PropostaID:   proposta.ID,
			Etapa:        "unica",
			AvaliadorID:  avaliadorID,
			Nota:         80,
			ParecerTexto: "",
		})
		if err == nil {
			t.Fatal("esperava erro para parecer textual vazio")
		}
	})
}
