package verificacao

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TurnstileVerificador struct {
	secretKey string
}

func NovoTurnstileVerificador() *TurnstileVerificador {
	sk := os.Getenv("TURNSTILE_SECRET_KEY")
	if sk == "" {
		sk = "1x0000000000000000000000000000000000000000000000000000000000000000AA"
	}
	return &TurnstileVerificador{secretKey: sk}
}

func (t *TurnstileVerificador) Verificar(token string) (bool, error) {
	if strings.HasPrefix(t.secretKey, "1x0000000000000000000000000000000") {
		return true, nil
	}

	if strings.TrimSpace(token) == "" {
		return false, nil
	}

	data := url.Values{}
	data.Set("secret", t.secretKey)
	data.Set("response", token)

	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", data)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar captcha: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("erro ao decodificar resposta do captcha: %w", err)
	}

	return result.Success, nil
}
