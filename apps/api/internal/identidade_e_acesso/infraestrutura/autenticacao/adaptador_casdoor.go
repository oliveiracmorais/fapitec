package autenticacao

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"golang.org/x/oauth2"
)

type AdaptadorCasdoor struct {
	client *casdoorsdk.Client
	owner  string
}

func NovoAdaptadorCasdoor(endpoint, clientId, clientSecret, certificate, organizationName, applicationName string) *AdaptadorCasdoor {
	c := casdoorsdk.NewClient(endpoint, clientId, clientSecret, certificate, organizationName, applicationName)
	return &AdaptadorCasdoor{client: c, owner: organizationName}
}

func (a *AdaptadorCasdoor) ValidarJWT(token string) (*casdoorsdk.Claims, error) {
	return a.client.ParseJwtToken(token)
}

func (a *AdaptadorCasdoor) TrocarCodigoPorToken(code string, state string) (string, error) {
	t, err := a.client.GetOAuthToken(code, state)
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}

func (a *AdaptadorCasdoor) VerificarPermissao(ctx context.Context, usuarioID, perfil, modulo, operacao string) (bool, error) {
	request := casdoorsdk.CasbinRequest{usuarioID, perfil, modulo, operacao}
	return a.client.Enforce("", "", "", "", a.owner, request)
}

func (a *AdaptadorCasdoor) GerarURLDeAutorizacao(redirectURI, state string) string {
	config := oauth2.Config{
		ClientID:     a.client.ClientId,
		ClientSecret: a.client.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/api/login/oauth/authorize", a.client.Endpoint),
			TokenURL:  fmt.Sprintf("%s/api/login/oauth/access_token", a.client.Endpoint),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: redirectURI,
		Scopes:      []string{"openid", "profile", "email"},
	}
	return config.AuthCodeURL(state)
}

func (a *AdaptadorCasdoor) CriarUsuario(ctx context.Context, nome, email, cpf, senha, perfil string) error {
	if strings.HasPrefix(a.client.Endpoint, "http://") {
		log.Printf("AVISO: Criando usuario no Casdoor via HTTP (sem SSL) — senha trafega em texto plano")
	}

	user := &casdoorsdk.User{
		Owner:       a.owner,
		Name:        cpf,
		DisplayName: nome,
		Email:       email,
		Password:    senha,
		IdCard:      cpf,
		Type:        perfil,
	}
	sucesso, err := a.client.AddUser(user)
	if err != nil {
		return fmt.Errorf("erro ao criar usuario no Casdoor: %w", err)
	}
	if !sucesso {
		return fmt.Errorf("casdoor retornou falha ao criar usuario")
	}
	return nil
}
