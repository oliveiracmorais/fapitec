package casos_de_uso

import (
	"context"
	"testing"

	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/infraestrutura/persistencia"
)

func TestCriarEdital(t *testing.T) {
	t.Run("deve criar edital com dados validos", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDeEditalMemoria()
		caso := NovoEdital(repo)

		entrada := dto.CriarEditalEntrada{
			Nome:                      "Edital APQ 2026",
			Descricao:                 "Edital de apoio a pesquisa",
			DataInicio:                "2026-06-01",
			DataFim:                   "2026-12-31",
			TipoChamada:               "APQ",
			ModeloFormulario:          1,
			TituloMinimoElegibilidade: "Doutor",
			ExigeEmpresa:              true,
			PorteEmpresa:              []string{"MEI", "ME"},
			DocumentosObrigatorios:    []string{"RG", "CPF"},
		}

		saida, err := caso.Executar(context.Background(), entrada)
		if err != nil {
			t.Fatalf("esperava nil, got %v", err)
		}
		if saida.Nome != "Edital APQ 2026" {
			t.Errorf("esperava 'Edital APQ 2026', got %s", saida.Nome)
		}
		if saida.Status != "ativo" {
			t.Errorf("esperava 'ativo', got %s", saida.Status)
		}
		if saida.ID == 0 {
			t.Error("esperava ID > 0")
		}
		if saida.ModeloFormulario != 1 {
			t.Errorf("esperava modelo_formulario=1, got %d", saida.ModeloFormulario)
		}
		if saida.TituloMinimoElegibilidade != "Doutor" {
			t.Errorf("esperava titulo_minimo_elegibilidade=Doutor, got %s", saida.TituloMinimoElegibilidade)
		}
		if !saida.ExigeEmpresa {
			t.Error("esperava exige_empresa=true")
		}
		if len(saida.PorteEmpresa) != 2 || saida.PorteEmpresa[0] != "MEI" {
			t.Error("porte_empresa incorreto")
		}
	})

	t.Run("deve rejeitar data invalida", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDeEditalMemoria()
		caso := NovoEdital(repo)

		entrada := dto.CriarEditalEntrada{
			Nome:        "Edital Teste",
			Descricao:   "Descricao",
			DataInicio:  "data-invalida",
			DataFim:     "2026-12-31",
			TipoChamada: "APQ",
		}

		_, err := caso.Executar(context.Background(), entrada)
		if err == nil {
			t.Fatal("esperava erro, got nil")
		}
	})

	t.Run("deve rejeitar nome vazio", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDeEditalMemoria()
		caso := NovoEdital(repo)

		entrada := dto.CriarEditalEntrada{
			Nome:        "",
			Descricao:   "Descricao",
			DataInicio:  "2026-06-01",
			DataFim:     "2026-12-31",
			TipoChamada: "APQ",
		}

		_, err := caso.Executar(context.Background(), entrada)
		if err == nil {
			t.Fatal("esperava erro, got nil")
		}
	})

	t.Run("deve criar edital com campos opcionais vazios", func(t *testing.T) {
		repo := persistencia.NovoRepositorioDeEditalMemoria()
		caso := NovoEdital(repo)

		entrada := dto.CriarEditalEntrada{
			Nome:        "Edital Simples",
			Descricao:   "Descricao basica",
			DataInicio:  "2026-07-01",
			DataFim:     "2026-12-31",
			TipoChamada: "APQ",
		}

		saida, err := caso.Executar(context.Background(), entrada)
		if err != nil {
			t.Fatalf("esperava nil, got %v", err)
		}
		if saida.ModeloFormulario != 0 {
			t.Errorf("esperava modelo_formulario=0, got %d", saida.ModeloFormulario)
		}
		if saida.ExigeEmpresa {
			t.Error("esperava exige_empresa=false")
		}
		if saida.TituloMinimoElegibilidade != "" {
			t.Errorf("esperava titulo_minimo_elegibilidade vazio, got %s", saida.TituloMinimoElegibilidade)
		}
	})
}
