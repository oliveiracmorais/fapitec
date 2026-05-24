package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oliveiracmorais/fapitec/api/internal/auditoria/infraestrutura/adaptadores"
	auditoriaPersistencia "github.com/oliveiracmorais/fapitec/api/internal/auditoria/infraestrutura/persistencia"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/casos_de_uso"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/aplicacao/dto"
	emailService "github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/email"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/hash"
	"github.com/oliveiracmorais/fapitec/api/internal/identidade_e_acesso/infraestrutura/persistencia"
)

func main() {
	repo := persistencia.NovoRepositorioDeUsuarioMemoria()
	tokenRepo := persistencia.NovoRepositorioDeTokenRedefinicaoMemoria()
	hashService := hash.NovoServicoDeHashBcrypt()
	emailLog := emailService.NovoServicoDeEmailLog()

	auditRepo := auditoriaPersistencia.NovoRepositorioDeEventosMemoria()
	auditService := adaptadores.NovoRegistradorAuditoria(auditRepo)

	cadastrar := casos_de_uso.NovoCadastrarUsuarioComAuditoria(repo, hashService, auditService)
	autenticar := casos_de_uso.NovoAutenticarUsuarioComAuditoria(repo, hashService, auditService)
	solicitarRedefinicao := casos_de_uso.NovoSolicitarRedefinicaoSenhaComAuditoria(repo, tokenRepo, emailLog, auditService)
	redefinirSenha := casos_de_uso.NovoRedefinirSenhaComAuditoria(repo, tokenRepo, hashService, auditService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/cadastro", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("POST /api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CPF   string `json:"cpf"`
			Senha string `json:"senha"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"erro":"requisicao invalida"}`, http.StatusBadRequest)
			return
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

	mux.HandleFunc("POST /api/v1/solicitar-redefinicao-senha", func(w http.ResponseWriter, r *http.Request) {
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
	})

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

	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
