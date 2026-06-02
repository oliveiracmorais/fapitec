package main

import (
	"fmt"
	"log"
	"os"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/joho/godotenv"
)

type perfilConfig struct {
	Name        string
	DisplayName string
	Description string
}

var perfis = []perfilConfig{
	{Name: "administrador_fapitec", DisplayName: "Administrador FAPITEC", Description: "Acesso total a todos os modulos e operacoes do sistema"},
	{Name: "instituicao_ensino", DisplayName: "Instituicao de Ensino", Description: "Cadastro de pesquisadores, submissao de projetos institucionais, acompanhamento de bolsas"},
	{Name: "funcionario_fapitec", DisplayName: "Funcionario FAPITEC", Description: "Acompanhamento de solicitacoes, validacao documental, apoio a analise financeira"},
	{Name: "diretoria", DisplayName: "Diretoria", Description: "Aprovacao/rejeicao de projetos, definicao de diretrizes, aprovacao final de pagamentos"},
	{Name: "proponente", DisplayName: "Proponente", Description: "Cadastro de perfil, submissao de propostas, solicitacao de auxilios, prestacao de contas"},
	{Name: "avaliador", DisplayName: "Avaliador", Description: "Recebimento de convites, avaliacao anonimizada de projetos, feedback e pontuacao"},
}

type moduloOp struct {
	Modulo    string
	Operacoes []string
}

var permissoesPorPerfil = map[string][]moduloOp{
	"administrador_fapitec": {
		{Modulo: "editais", Operacoes: []string{"visualizar", "criar", "editar", "excluir", "administrar"}},
		{Modulo: "dashboard", Operacoes: []string{"visualizar", "administrar"}},
		{Modulo: "usuarios", Operacoes: []string{"visualizar", "criar", "editar", "excluir", "administrar"}},
		{Modulo: "financeiro", Operacoes: []string{"visualizar", "criar", "editar", "excluir", "administrar"}},
		{Modulo: "relatorios", Operacoes: []string{"visualizar", "criar", "administrar"}},
		{Modulo: "comunicacao", Operacoes: []string{"visualizar", "criar", "editar", "administrar"}},
		{Modulo: "documentos", Operacoes: []string{"visualizar", "criar", "editar", "excluir", "administrar"}},
		{Modulo: "conformidade", Operacoes: []string{"visualizar", "criar", "editar", "administrar"}},
	},
	"diretoria": {
		{Modulo: "editais", Operacoes: []string{"visualizar", "administrar"}},
		{Modulo: "dashboard", Operacoes: []string{"visualizar"}},
		{Modulo: "relatorios", Operacoes: []string{"visualizar", "administrar"}},
		{Modulo: "financeiro", Operacoes: []string{"visualizar", "editar", "administrar"}},
	},
	"funcionario_fapitec": {
		{Modulo: "editais", Operacoes: []string{"visualizar", "criar", "editar"}},
		{Modulo: "dashboard", Operacoes: []string{"visualizar"}},
		{Modulo: "documentos", Operacoes: []string{"visualizar", "criar", "editar"}},
		{Modulo: "comunicacao", Operacoes: []string{"visualizar", "criar"}},
	},
	"instituicao_ensino": {
		{Modulo: "editais", Operacoes: []string{"visualizar"}},
		{Modulo: "dashboard", Operacoes: []string{"visualizar"}},
		{Modulo: "documentos", Operacoes: []string{"visualizar", "criar"}},
	},
	"proponente": {
		{Modulo: "editais", Operacoes: []string{"visualizar"}},
		{Modulo: "dashboard", Operacoes: []string{"visualizar"}},
		{Modulo: "documentos", Operacoes: []string{"visualizar", "criar"}},
	},
	"avaliador": {
		{Modulo: "editais", Operacoes: []string{"visualizar"}},
		{Modulo: "dashboard", Operacoes: []string{"visualizar"}},
	},
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado, usando variaveis de ambiente")
	}
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Nenhum arquivo .env na raiz do monorepo")
	}

	endpoint := os.Getenv("CASDOOR_ENDPOINT")
	clientID := os.Getenv("CASDOOR_CLIENT_ID")
	clientSecret := os.Getenv("CASDOOR_CLIENT_SECRET")
	certificate := os.Getenv("CASDOOR_CERTIFICATE")
	orgName := os.Getenv("CASDOOR_ORGANIZATION_NAME")

	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}
	if orgName == "" {
		orgName = "fapitec"
	}
	if clientID == "" || certificate == "" {
		log.Fatal("Variaveis CASDOOR_CLIENT_ID e CASDOOR_CERTIFICATE sao obrigatorias")
	}

	client := casdoorsdk.NewClient(endpoint, clientID, clientSecret, certificate, orgName, orgName)
	fmt.Printf("Conectado ao Casdoor em %s, organizacao: %s\n", endpoint, orgName)

	for _, p := range perfis {
		existing, err := client.GetRole(p.Name)
		if err == nil && existing != nil && existing.Name != "" {
			fmt.Printf("  Role ja existe: %s, pulando\n", p.Name)
			continue
		}

		role := &casdoorsdk.Role{
			Owner:       orgName,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			Description: p.Description,
			IsEnabled:   true,
		}

		sucesso, err := client.AddRole(role)
		if err != nil {
			fmt.Printf("  ERRO ao criar role %s: %v\n", p.Name, err)
			continue
		}
		if sucesso {
			fmt.Printf("  Role criada: %s (%s)\n", p.Name, p.DisplayName)
		} else {
			fmt.Printf("  Falha ao criar role: %s\n", p.Name)
		}
	}

	for perfil, modulos := range permissoesPorPerfil {
		for _, mp := range modulos {
			permName := fmt.Sprintf("perm_%s_%s", perfil, mp.Modulo)

			existing, err := client.GetPermission(permName)
			if err == nil && existing != nil && existing.Name != "" {
				fmt.Printf("  Permissao ja existe: %s, pulando\n", permName)
				continue
			}

			roleID := fmt.Sprintf("%s/%s", orgName, perfil)
			perm := &casdoorsdk.Permission{
				Owner:        orgName,
				Name:         permName,
				DisplayName:  fmt.Sprintf("Permissao %s - %s", perfil, mp.Modulo),
				Description:  fmt.Sprintf("Permissoes do perfil %s para o modulo %s", perfil, mp.Modulo),
				Roles:        []string{roleID},
				ResourceType: "modulo",
				Resources:    []string{mp.Modulo},
				Actions:      mp.Operacoes,
				Effect:       "Allow",
				IsEnabled:    true,
			}

			sucesso, err := client.AddPermission(perm)
			if err != nil {
				fmt.Printf("  ERRO ao criar permissao %s: %v\n", permName, err)
				continue
			}
			if sucesso {
				fmt.Printf("  Permissao criada: %s\n", permName)
			}
		}
	}

	adminName := "admin_fapitec"
	existingUser, err := client.GetUser(adminName)
	if err == nil && existingUser != nil && existingUser.Name != "" {
		fmt.Printf("  Usuario admin ja existe: %s\n", adminName)
	} else {
		admin := &casdoorsdk.User{
			Owner:       orgName,
			Name:        adminName,
			DisplayName: "Administrador FAPITEC",
			Email:       "admin@fapitec.se.gov.br",
			Password:    "Fapitec@2026",
			Type:        "administrador_fapitec",
			IsAdmin:     true,
		}

		sucesso, err := client.AddUser(admin)
		if err != nil {
			fmt.Printf("  ERRO ao criar usuario admin: %v\n", err)
		} else if sucesso {
			fmt.Printf("  Usuario admin criado: %s (email: admin@fapitec.se.gov.br, senha: Fapitec@2026)\n", adminName)
			fmt.Println("  AVISO: Altere a senha do admin no primeiro acesso via UI do Casdoor!")
		}
	}

	fmt.Println("\nSeed concluido!")
}
