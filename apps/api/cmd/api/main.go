package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oliveiracmorais/fapitec/api/internal/auditoria/infraestrutura/adaptadores"
	auditoriaPersistencia "github.com/oliveiracmorais/fapitec/api/internal/auditoria/infraestrutura/persistencia"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/casos_de_uso"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	gestaoEditais 	"github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/casos_de_uso"
	gestaoEditaisDTO "github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/aplicacao/dto"
	gestaoEditaisRepositorios "github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/dominio/repositorios"
	gestaoEditaisPersistencia "github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/infraestrutura/persistencia"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/entidades"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/repositorios"
	emailService "github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/email"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/verificacao"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/hash"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia"
	sqlcgedital "github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/infraestrutura/persistencia/sqlc"
	sqlcpersistencia "github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia/sqlc"
)

func main() {
	hashService := hash.NovoServicoDeHashBcrypt()
	turnstileVerificador := verificacao.NovoTurnstileVerificador()
	emailLog := emailService.NovoServicoDeEmailLog()

	auditRepo := auditoriaPersistencia.NovoRepositorioDeEventosMemoria()
	auditService := adaptadores.NovoRegistradorAuditoria(auditRepo)

	var repo repositorios.RepositorioDeUsuario
	var tokenRepo repositorios.RepositorioDeTokenRedefinicao
	var editalRepo gestaoEditaisRepositorios.RepositorioDeEdital

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://fapitec:fapitec123@localhost:5432/fapitec?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err == nil {
		if err := pool.Ping(context.Background()); err == nil {
			queries := sqlcpersistencia.New(pool)
			repo = persistencia.NovoRepositorioDeUsuarioSQLC(queries)
			tokenRepo = persistencia.NovoRepositorioDeTokenRedefinicaoSQLC(queries)

			queriesEdital := sqlcgedital.New(pool)
			editalRepo = gestaoEditaisPersistencia.NovoRepositorioDeEditalSQLC(queriesEdital)

			log.Println("Conectado ao PostgreSQL")
		} else {
			pool.Close()
		}
	}

	if repo == nil {
		log.Println("PostgreSQL indisponivel — usando repositorio em memoria")
		repo = persistencia.NovoRepositorioDeUsuarioMemoria()
		tokenRepo = persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	}

	if editalRepo == nil {
		log.Println("PostgreSQL indisponivel para editais — usando repositorio em memoria")
		editalRepo = gestaoEditaisPersistencia.NovoRepositorioDeEditalMemoria()
	}

	cadastrar := casos_de_uso.NovoCadastrarUsuarioComAuditoria(repo, hashService, auditService)
	autenticar := casos_de_uso.NovoAutenticarUsuarioComAuditoria(repo, hashService, auditService)
	solicitarRedefinicao := casos_de_uso.NovoSolicitarRedefinicaoSenhaComAuditoria(repo, tokenRepo, emailLog, auditService)
	redefinirSenha := casos_de_uso.NovoRedefinirSenhaComAuditoria(repo, tokenRepo, hashService, auditService)

	criarEdital := gestaoEditais.NovoEdital(editalRepo)
	listarEditais := gestaoEditais.NovoListarEditais(editalRepo)
	visualizarEdital := gestaoEditais.NovoVisualizarEdital(editalRepo)
	atualizarEdital := gestaoEditais.NovoAtualizarEdital(editalRepo)
	deletarEdital := gestaoEditais.NovoDeletarEdital(editalRepo)

	mux := http.NewServeMux()

	cadastroHandler := func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Nome             string `json:"nome"`
			CPF              string `json:"cpf"`
			Email            string `json:"email"`
			ConfirmacaoEmail string `json:"confirmacao_email"`
			Senha            string `json:"senha"`
			ConfirmacaoSenha string `json:"confirmacao_senha"`
			Estrangeiro      bool   `json:"estrangeiro"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		entrada := dto.CadastrarUsuarioEntrada{
			Nome:             req.Nome,
			CPF:              req.CPF,
			Email:            req.Email,
			ConfirmacaoEmail: req.ConfirmacaoEmail,
			Senha:            req.Senha,
			ConfirmacaoSenha: req.ConfirmacaoSenha,
			Estrangeiro:      req.Estrangeiro,
		}

		saida, err := cadastrar.Executar(context.Background(), entrada)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(saida)
	}

	mux.HandleFunc("POST /api/v1/cadastro", cadastroHandler)
	mux.HandleFunc("POST /api/v1/register", cadastroHandler)

	mux.HandleFunc("GET /api/v1/check-cpf", func(w http.ResponseWriter, r *http.Request) {
		cpf := r.URL.Query().Get("cpf")
		usuario, _ := repo.BuscarPorCPF(context.Background(), cpf)
		existe := usuario != nil
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"existe": existe})
	})

	mux.HandleFunc("GET /api/v1/check-email", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		usuario, _ := repo.BuscarPorEmail(context.Background(), email)
		existe := usuario != nil
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"existe": existe})
	})

	mux.HandleFunc("POST /api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CPF          string `json:"cpf"`
			Senha        string `json:"senha"`
			CaptchaToken string `json:"captcha_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		cpfLimpo := regexp.MustCompile(`\D`).ReplaceAllString(req.CPF, "")
		usuario, _ := repo.BuscarPorCPF(r.Context(), cpfLimpo)
		if usuario != nil && usuario.Tentativas >= 3 {
			valido, _ := turnstileVerificador.Verificar(req.CaptchaToken)
			if !valido {
				http.Error(w, `{"erro":"Validacao de captcha falhou. Tente novamente."}`, http.StatusUnauthorized)
				return
			}
		}

		entrada := dto.AutenticarUsuarioEntrada{
			CPF:   req.CPF,
			Senha: req.Senha,
		}

		saida, err := autenticar.Executar(context.Background(), entrada)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("GET /api/v1/user-profile", func(w http.ResponseWriter, r *http.Request) {
		cpf := r.URL.Query().Get("cpf")
		email := r.URL.Query().Get("email")

		if cpf == "" && email == "" {
			http.Error(w, `{"erro":"informe cpf ou email"}`, http.StatusBadRequest)
			return
		}

		var usuario *entidades.Usuario
		if cpf != "" {
			usuario, _ = repo.BuscarPorCPF(context.Background(), cpf)
		} else {
			usuario, _ = repo.BuscarPorEmail(context.Background(), email)
		}

		if usuario == nil {
			http.Error(w, `{"erro":"usuario nao encontrado"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":          usuario.ID,
			"nome":        usuario.Nome,
			"documento":   usuario.CPF,
			"email":       usuario.Email.String(),
			"estrangeiro": usuario.Estrangeiro,
			"criado_em":   usuario.CriadoEm.Format("2006-01-02T15:04:05-07:00"),
		})
	})

	solicitarRedefinicaoHandler := func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		entrada := dto.SolicitarRedefinicaoSenhaEntrada{
			Email: req.Email,
		}

		err := solicitarRedefinicao.Executar(context.Background(), entrada)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"mensagem": "Se o e-mail estiver cadastrado, voce recebera um link de redefinicao de senha.",
		})
	}

	mux.HandleFunc("POST /api/v1/solicitar-redefinicao-senha", solicitarRedefinicaoHandler)
	mux.HandleFunc("POST /api/v1/reset-password", solicitarRedefinicaoHandler)

	mux.HandleFunc("POST /api/v1/redefinir-senha", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Token            string `json:"token"`
			Senha            string `json:"senha"`
			ConfirmacaoSenha string `json:"confirmacao_senha"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		entrada := dto.RedefinirSenhaEntrada{
			Token:            req.Token,
			Senha:            req.Senha,
			ConfirmacaoSenha: req.ConfirmacaoSenha,
		}

		err := redefinirSenha.Executar(context.Background(), entrada)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"mensagem": "Senha redefinida com sucesso.",
		})
	})

	mux.HandleFunc("GET /api/v1/auditoria", func(w http.ResponseWriter, r *http.Request) {
		eventos, err := auditRepo.Listar(context.Background())
		if err != nil {
			http.Error(w, `{"erro":"erro ao listar eventos de auditoria"}`, http.StatusInternalServerError)
			return
		}
		resp := make([]map[string]any, 0, len(eventos))
		for _, e := range eventos {
			resp = append(resp, map[string]any{
				"id":         e.ID,
				"acao":       e.Acao,
				"ator_id":    e.AtorID,
				"ator_nome":  e.AtorNome,
				"ator_cpf":   e.AtorCPF,
				"perfil":     e.Perfil,
				"resultado":  e.Resultado,
				"modulo":     e.Modulo,
				"recurso":    e.Recurso,
				"origem":     e.Origem,
				"data_hora":  e.DataHora.Format("2006-01-02T15:04:05-07:00"),
				"contexto":   e.Contexto,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("POST /api/v1/editais", func(w http.ResponseWriter, r *http.Request) {
		var req gestaoEditaisDTO.CriarEditalEntrada
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		saida, err := criarEdital.Executar(context.Background(), req)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("GET /api/v1/editais", func(w http.ResponseWriter, r *http.Request) {
		filtros := gestaoEditais.FiltrosListarEditais{
			Titulo:      r.URL.Query().Get("titulo"),
			Status:      r.URL.Query().Get("status"),
			TipoChamada: r.URL.Query().Get("tipo_chamada"),
		}

		saida, err := listarEditais.Executar(context.Background(), filtros)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("GET /api/v1/editais/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			http.Error(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		saida, err := visualizarEdital.Executar(context.Background(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("PUT /api/v1/editais/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			http.Error(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		var req gestaoEditaisDTO.AtualizarEditalEntrada
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		saida, err := atualizarEdital.Executar(context.Background(), id, req)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("DELETE /api/v1/editais/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			http.Error(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		if err := deletarEdital.Executar(context.Background(), id); err != nil {
			http.Error(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/v1/dashboard/indicadores", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"indicadores": []map[string]any{
				{"id": "iso-27001", "nome": "Progresso ISO 27001", "valor": 65, "tipo": "porcentagem", "cor": "violeta"},
				{"id": "nao-conformidades", "nome": "Não Conformidades Tratadas", "valor": 78, "tipo": "porcentagem", "cor": "verde"},
				{"id": "valor-absoluto", "nome": "Valor em Projetos", "valor": 5000, "tipo": "numero", "cor": "azul"},
				{"id": "status-geral", "nome": "Status Geral", "valor": 82, "tipo": "barra", "cor": "laranja"},
			},
		})
	})

	mux.HandleFunc("GET /api/v1/dashboard/graficos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"graficos": []map[string]any{
				{
					"id": "donut", "tipo": "donut", "titulo": "Distribuição por Tipo de Projeto",
					"faixas": []map[string]any{
						{"nome": "Pesquisa", "valor": 60, "legenda": "60%"},
						{"nome": "Inovação", "valor": 25, "legenda": "25%"},
						{"nome": "Extensão", "valor": 15, "legenda": "15%"},
					},
				},
				{
					"id": "evolucao", "tipo": "linha", "titulo": "Evolução Mensal de Projetos",
					"dados": []map[string]any{
						{"mes": "Jan", "valor": 12},
						{"mes": "Fev", "valor": 19},
						{"mes": "Mar", "valor": 15},
						{"mes": "Abr", "valor": 22},
						{"mes": "Mai", "valor": 28},
					},
				},
			},
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
