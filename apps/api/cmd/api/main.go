package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/dominio/servicos"
	emailService "github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/email"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/autenticacao"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/verificacao"
	interfacesHTTP "github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/interfaces/http"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/hash"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia"
	sqlcgedital "github.com/oliveiracmorais/fapitec/api/internal/gestao_de_editais/infraestrutura/persistencia/sqlc"
	sqlcpersistencia "github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia/sqlc"
	propostaCasosDeUso "github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/casos_de_uso"
	propostaDTO "github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/aplicacao/dto"
	propostaRepositorios "github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/dominio/repositorios"
	propostaPersistencia "github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia"
	sqlcproposta "github.com/oliveiracmorais/fapitec/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia/sqlc"
)

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado no diretorio atual, usando variaveis de ambiente")
	}
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Nenhum arquivo .env encontrado na raiz do monorepo, usando variaveis de ambiente")
	}

	hashService := hash.NovoServicoDeHashBcrypt()
	turnstileVerificador := verificacao.NovoTurnstileVerificador()

	var emailSvc servicos.ServicoDeEmail
	if emailService.ConfigSMTPPresente() {
		emailSvc = emailService.NovoServicoDeEmailSMTP()
		log.Println("Servico de email SMTP configurado")
	} else {
		emailSvc = emailService.NovoServicoDeEmailLog()
		log.Println("Servico de email: modo placeholder (sem configuracao SMTP)")
	}

	auditRepo := auditoriaPersistencia.NovoRepositorioDeEventosMemoria()
	auditService := adaptadores.NovoRegistradorAuditoria(auditRepo)

	var repo repositorios.RepositorioDeUsuario
	var tokenRepo repositorios.RepositorioDeTokenRedefinicao
	var editalRepo gestaoEditaisRepositorios.RepositorioDeEdital
	var propostaRepo propostaRepositorios.RepositorioDeProposta

	authProvider := os.Getenv("AUTH_PROVIDER")
	if authProvider == "" {
		authProvider = "internal"
	}
	log.Printf("AUTH_PROVIDER=%s", authProvider)

	var casdoorAdapter *autenticacao.AdaptadorCasdoor
	casdoorEndpoint := os.Getenv("CASDOOR_ENDPOINT")
	if casdoorEndpoint != "" {
		casdoorAdapter = autenticacao.NovoAdaptadorCasdoor(
			casdoorEndpoint,
			os.Getenv("CASDOOR_CLIENT_ID"),
			os.Getenv("CASDOOR_CLIENT_SECRET"),
			os.Getenv("CASDOOR_CERTIFICATE"),
			os.Getenv("CASDOOR_ORGANIZATION_NAME"),
			os.Getenv("CASDOOR_ORGANIZATION_NAME"),
		)
		log.Printf("Adaptador Casdoor inicializado (endpoint: %s)", casdoorEndpoint)
	}

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

			queriesProposta := sqlcproposta.New(pool)
			propostaRepo = propostaPersistencia.NovoRepositorioDePropostaSQLC(queriesProposta)

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

	if propostaRepo == nil {
		log.Println("PostgreSQL indisponivel para propostas — usando repositorio em memoria")
		propostaRepo = propostaPersistencia.NovoRepositorioDePropostaMemoria()
	}

	cadastrar := casos_de_uso.NovoCadastrarUsuarioComAuditoria(repo, hashService, auditService)
	autenticar := casos_de_uso.NovoAutenticarUsuarioComAuditoria(repo, hashService, auditService)
	solicitarRedefinicao := casos_de_uso.NovoSolicitarRedefinicaoSenhaComAuditoria(repo, tokenRepo, emailSvc, auditService)
	redefinirSenha := casos_de_uso.NovoRedefinirSenhaComAuditoria(repo, tokenRepo, hashService, auditService)

	criarEdital := gestaoEditais.NovoEdital(editalRepo)
	listarEditais := gestaoEditais.NovoListarEditais(editalRepo)
	visualizarEdital := gestaoEditais.NovoVisualizarEdital(editalRepo)
	atualizarEdital := gestaoEditais.NovoAtualizarEdital(editalRepo)
	deletarEdital := gestaoEditais.NovoDeletarEdital(editalRepo)

	criarProposta := propostaCasosDeUso.NovoCriarProposta(propostaRepo)

	editalVerificador := &editalVerificadorAdapter{repo: editalRepo}
	submeterProposta := propostaCasosDeUso.NovoSubmeterProposta(propostaRepo, editalVerificador)
	editarProposta := propostaCasosDeUso.NovoEditarProposta(propostaRepo)
	listarPropostas := propostaCasosDeUso.NovoListarPropostas(propostaRepo)
	visualizarProposta := propostaCasosDeUso.NovoVisualizarProposta(propostaRepo)
	deletarProposta := propostaCasosDeUso.NovoDeletarProposta(propostaRepo)

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
			jsonError(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
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
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		if casdoorAdapter != nil {
			perfil := "proponente"
			if err := casdoorAdapter.CriarUsuario(context.Background(), req.Nome, req.Email, req.CPF, req.Senha, perfil); err != nil {
				log.Printf("Aviso: usuario criado localmente mas falha ao criar no Casdoor: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(saida)
	}

	mux.HandleFunc("POST /api/v1/cadastro", cadastroHandler)
	mux.HandleFunc("POST /api/v1/register", cadastroHandler)
		mux.HandleFunc("POST /api/v1/login", func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				CPF          string `json:"cpf"`
				Senha        string `json:"senha"`
				CaptchaToken string `json:"captcha_token"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				jsonError(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
				return
			}

			cpfLimpo := regexp.MustCompile(`\D`).ReplaceAllString(req.CPF, "")
			usuario, _ := repo.BuscarPorCPF(r.Context(), cpfLimpo)
			if usuario != nil && usuario.Tentativas >= 3 {
				valido, _ := turnstileVerificador.Verificar(req.CaptchaToken)
				if !valido {
					jsonError(w, `{"erro":"Validacao de captcha falhou. Tente novamente."}`, http.StatusUnauthorized)
					return
				}
			}

			entrada := dto.AutenticarUsuarioEntrada{
				CPF:   req.CPF,
				Senha: req.Senha,
			}

			saida, err := autenticar.Executar(context.Background(), entrada)
			if err != nil {
				jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusUnauthorized)
				return
			}

			payload := map[string]interface{}{
				"sub":   fmt.Sprintf("%d", saida.ID),
				"name":  saida.Nome,
				"email": saida.Email,
				"iat":   time.Now().Unix(),
				"exp":   time.Now().Add(24 * time.Hour).Unix(),
			}
			payloadBytes, _ := json.Marshal(payload)

			header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
			payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)
			signature := base64.RawURLEncoding.EncodeToString([]byte{})
			token := fmt.Sprintf("%s.%s.%s", header, payloadB64, signature)

			http.SetCookie(w, &http.Cookie{
				Name:     "fapitec_token",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   86400,
			})

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(saida)
		})
		solicitarRedefinicaoHandler := func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Email string `json:"email"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				jsonError(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
				return
			}

			entrada := dto.SolicitarRedefinicaoSenhaEntrada{
				Email: req.Email,
			}

			err := solicitarRedefinicao.Executar(context.Background(), entrada)
			if err != nil {
				jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
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
				jsonError(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
				return
			}

			entrada := dto.RedefinirSenhaEntrada{
				Token:            req.Token,
				Senha:            req.Senha,
				ConfirmacaoSenha: req.ConfirmacaoSenha,
			}

			err := redefinirSenha.Executar(context.Background(), entrada)
			if err != nil {
				jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"mensagem": "Senha redefinida com sucesso.",
			})
		})

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

	mux.HandleFunc("GET /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if casdoorAdapter == nil {
			jsonError(w, `{"erro":"provedor casdoor nao configurado"}`, http.StatusBadRequest)
			return
		}
		redirectURI := r.URL.Query().Get("redirect_uri")
		if redirectURI == "" {
			scheme := r.URL.Scheme
			if scheme == "" {
				scheme = r.Header.Get("X-Forwarded-Proto")
			}
			if scheme == "" {
				scheme = "http"
			}
			redirectURI = fmt.Sprintf("%s://%s/api/v1/auth/callback", scheme, r.Host)
		}
		state := r.URL.Query().Get("state")
		if state == "" {
			state = "fapitec-state"
		}
		url := casdoorAdapter.GerarURLDeAutorizacao(redirectURI, state)
		http.Redirect(w, r, url, http.StatusFound)
	})

	mux.HandleFunc("GET /api/v1/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		if casdoorAdapter == nil {
			jsonError(w, `{"erro":"provedor casdoor nao configurado"}`, http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		if code == "" {
			jsonError(w, `{"erro":"codigo de autorizacao ausente"}`, http.StatusBadRequest)
			return
		}
		token, err := casdoorAdapter.TrocarCodigoPorToken(code, state)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"falha ao obter token: %s"}`, err.Error()), http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "fapitec_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400,
		})
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})

	mux.HandleFunc("GET /api/v1/user-profile", func(w http.ResponseWriter, r *http.Request) {
		if claims, ok := r.Context().Value(interfacesHTTP.UsuarioContextKey).(*casdoorsdk.Claims); ok {
			nome := claims.DisplayName
			if nome == "" {
				nome = claims.Name
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"documento": claims.Name,
				"nome":      nome,
				"email":     claims.Email,
				"perfil":    claims.Type,
			})
			return
		}

		cpf := r.URL.Query().Get("cpf")
		email := r.URL.Query().Get("email")

		if cpf == "" && email == "" {
			jsonError(w, `{"erro":"informe cpf ou email"}`, http.StatusBadRequest)
			return
		}

		var usuario *entidades.Usuario
		if cpf != "" {
			usuario, _ = repo.BuscarPorCPF(context.Background(), cpf)
		} else {
			usuario, _ = repo.BuscarPorEmail(context.Background(), email)
		}

		if usuario == nil {
			jsonError(w, `{"erro":"usuario nao encontrado"}`, http.StatusNotFound)
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

	mux.HandleFunc("GET /api/v1/auditoria", func(w http.ResponseWriter, r *http.Request) {
		eventos, err := auditRepo.Listar(context.Background())
		if err != nil {
			jsonError(w, `{"erro":"erro ao listar eventos de auditoria"}`, http.StatusInternalServerError)
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
			jsonError(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		saida, err := criarEdital.Executar(context.Background(), req)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
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
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("GET /api/v1/editais/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			jsonError(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		saida, err := visualizarEdital.Executar(context.Background(), id)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("PUT /api/v1/editais/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			jsonError(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		var req gestaoEditaisDTO.AtualizarEditalEntrada
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		saida, err := atualizarEdital.Executar(context.Background(), id, req)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("DELETE /api/v1/editais/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			jsonError(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		if err := deletarEdital.Executar(context.Background(), id); err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("POST /api/v1/propostas", func(w http.ResponseWriter, r *http.Request) {
		var req propostaDTO.CriarPropostaEntrada
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		saida, err := criarProposta.Executar(context.Background(), req)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("GET /api/v1/propostas", func(w http.ResponseWriter, r *http.Request) {
		filtros := propostaRepositorios.FiltrosListarPropostas{
			ProponenteID: 0,
			EditalID:     0,
			Status:       r.URL.Query().Get("status"),
		}
		if idStr := r.URL.Query().Get("edital_id"); idStr != "" {
			fmt.Sscanf(idStr, "%d", &filtros.EditalID)
		}
		if idStr := r.URL.Query().Get("proponente_id"); idStr != "" {
			fmt.Sscanf(idStr, "%d", &filtros.ProponenteID)
		}

		saida, err := listarPropostas.Executar(context.Background(), filtros)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("GET /api/v1/propostas/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			jsonError(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		saida, err := visualizarProposta.Executar(context.Background(), id)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("PUT /api/v1/propostas/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			jsonError(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		var req propostaDTO.AtualizarPropostaEntrada
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
		}

		saida, err := editarProposta.Executar(context.Background(), id, req)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("DELETE /api/v1/propostas/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			jsonError(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		if err := deletarProposta.Executar(context.Background(), id); err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("POST /api/v1/propostas/{id}/submeter", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			jsonError(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		saida, err := submeterProposta.Executar(context.Background(), id)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("GET /api/v1/editais/{id}/propostas", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		var editalID int64
		if _, err := fmt.Sscanf(idStr, "%d", &editalID); err != nil {
			jsonError(w, `{"erro":"id invalido"}`, http.StatusBadRequest)
			return
		}

		filtros := propostaRepositorios.FiltrosListarPropostas{
			EditalID: editalID,
		}

		saida, err := listarPropostas.Executar(context.Background(), filtros)
		if err != nil {
			jsonError(w, fmt.Sprintf(`{"erro":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(saida)
	})

	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/v1/dashboard/indicadores", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"indicadores": []map[string]any{
				{"id": "iso-27001", "nome": "Progresso ISO 27001", "valor": 65, "tipo": "pizza", "cor": "violeta"},
				{"id": "nao-conformidades", "nome": "Não Conformidades Tratadas", "valor": 78, "tipo": "pizza", "cor": "verde"},
				{"id": "valor-absoluto", "nome": "Valor Absoluto", "valor": 5000, "tipo": "linha", "cor": "azul", "dados": []map[string]any{
					{"mes": "Jan", "valor": 3200},
					{"mes": "Fev", "valor": 4100},
					{"mes": "Mar", "valor": 3800},
					{"mes": "Abr", "valor": 4600},
					{"mes": "Mai", "valor": 5000},
				}},
				{"id": "status-geral", "nome": "Status do Gráfico", "valor": 82, "tipo": "barra", "cor": "laranja"},
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

	var handler http.Handler = mux

	publicPaths := map[string]bool{
		"/api/v1/health":          true,
		"/api/v1/check-cpf":       true, "/api/v1/check-email": true,
		"/api/v1/auth/login":      true,
		"/api/v1/auth/callback":   true,
		"/api/v1/login":           true,
		"/api/v1/cadastro":        true,
		"/api/v1/register":        true,
		"/api/v1/solicitar-redefinicao-senha": true,
		"/api/v1/reset-password":  true,
		"/api/v1/redefinir-senha": true,
	}

	publicPrefixes := []string{
		"/api/v1/editais",
	}

	if casdoorAdapter != nil {
		authMW := interfacesHTTP.AutenticacaoMiddleware(casdoorAdapter)

		protectedHandler := authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux.ServeHTTP(w, r)
		}))

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if publicPaths[r.URL.Path] {
				mux.ServeHTTP(w, r)
				return
			}
			for _, prefix := range publicPrefixes {
				if strings.HasPrefix(r.URL.Path, prefix) {
					mux.ServeHTTP(w, r)
					return
				}
			}
			protectedHandler.ServeHTTP(w, r)
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

type editalVerificadorAdapter struct {
	repo gestaoEditaisRepositorios.RepositorioDeEdital
}

func (a *editalVerificadorAdapter) BuscarEditalInfo(ctx context.Context, id int64) (*propostaRepositorios.EditalInfo, error) {
	edital, err := a.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	if edital == nil {
		return nil, nil
	}
	return &propostaRepositorios.EditalInfo{
		Status:     edital.Status.String(),
		DataInicio: edital.DataInicio,
		DataFim:    edital.DataFim,
	}, nil
}
