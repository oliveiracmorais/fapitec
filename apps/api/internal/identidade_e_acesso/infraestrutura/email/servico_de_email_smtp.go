package email

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
)

type ServicoDeEmailSMTP struct {
	host     string
	port     string
	user     string
	pass     string
	from     string
	frontend string
}

func NovoServicoDeEmailSMTP() *ServicoDeEmailSMTP {
	return &ServicoDeEmailSMTP{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		user:     os.Getenv("SMTP_USER"),
		pass:     os.Getenv("SMTP_PASS"),
		from:     os.Getenv("SMTP_FROM"),
		frontend: FrontendPadrao(),
	}
}

func (s *ServicoDeEmailSMTP) EnviarRedefinicaoSenha(ctx context.Context, email, token string) error {
	link := fmt.Sprintf("%s/redefinir-senha?token=%s", s.frontend, token)

	subject := "Redefinicao de Senha - FAPITEC-SE"
	body := s.montarTemplateHTML(link)

	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s", s.from, email, subject, body)

	return smtp.SendMail(addr, auth, s.from, []string{email}, []byte(msg))
}

func (s *ServicoDeEmailSMTP) montarTemplateHTML(link string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="pt-BR">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif">
<table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:40px 0">
<tr><td align="center">
<table width="600" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.1)">
<tr><td style="background-color:#1e3a5f;padding:30px 40px;text-align:center">
<h1 style="color:#ffffff;margin:0;font-size:22px;font-weight:bold">FAPITEC-SE</h1>
<p style="color:#b0c4de;margin:5px 0 0;font-size:14px">Fundacao de Apoio a Pesquisa e Inovacao Tecnologica de Sergipe</p>
</td></tr>
<tr><td style="padding:40px">
<h2 style="color:#333333;font-size:20px;margin:0 0 20px">Redefinicao de Senha</h2>
<p style="color:#555555;font-size:15px;line-height:1.6;margin:0 0 20px">Recebemos uma solicitacao de redefinicao de senha para sua conta na plataforma FAPITEC-SE.</p>
<p style="color:#555555;font-size:15px;line-height:1.6;margin:0 0 20px">Clique no botao abaixo para criar uma nova senha. Este link e valido por 1 hora.</p>
<table cellpadding="0" cellspacing="0" style="margin:30px 0">
<tr><td align="center" style="background-color:#1e3a5f;border-radius:6px;padding:12px 32px">
<a href="%s" style="color:#ffffff;font-size:16px;font-weight:bold;text-decoration:none">Redefinir Senha</a>
</td></tr>
</table>
<p style="color:#999999;font-size:13px;line-height:1.5;margin:0">Se voce nao solicitou esta redefinicao, ignore este email. Nenhuma alteracao sera feita na sua conta.</p>
<hr style="border:none;border-top:1px solid #eeeeee;margin:30px 0">
<p style="color:#999999;font-size:12px;margin:0">Caso o botao nao funcione, copie e cole o link abaixo no seu navegador:</p>
<p style="color:#999999;font-size:12px;margin:10px 0 0;word-break:break-all">%s</p>
</td></tr>
<tr><td style="background-color:#f8f8f8;padding:20px 40px;text-align:center">
<p style="color:#aaaaaa;font-size:12px;margin:0">FAPITEC-SE — Fundacao de Apoio a Pesquisa e Inovacao Tecnologica de Sergipe</p>
</td></tr>
</table>
</td></tr>
</table>
</body>
</html>`, link, link)
}

func ConfigSMTPPresente() bool {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	frontend := os.Getenv("FRONTEND_URL")
	if host == "" || port == "" || user == "" || pass == "" || from == "" || frontend == "" {
		return false
	}
	return true
}

func FrontendPadrao() string {
	if f := os.Getenv("FRONTEND_URL"); f != "" {
		return f
	}
	return "http://localhost:3000"
}
