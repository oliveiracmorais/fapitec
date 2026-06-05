# Ambiente de Homologação

## Infraestrutura

VM Oracle Cloud (Ubuntu), domínios via DuckDNS, TLS via Let's Encrypt.

## Domínios

| Domínio | Serviço | Proxy |
|---------|---------|-------|
| `fapitec.duckdns.org` | Dashboard (Next.js) + API (Go) | nginx do host → `web:3000` / `api:8080` |
| `auth.fapitec.duckdns.org` | Casdoor UI | nginx do host → `127.0.0.1:8000` |

## Serviços (Docker Compose)

```
docker-compose.yml  →  ~/apps/fapitec/
```

| Serviço    | Porta host          | Porta container | Dockerfile                  |
|------------|---------------------|-----------------|-----------------------------|
| `postgres` | `127.0.0.1:5432`    | 5432            | postgres:16-alpine          |
| `casdoor`  | `127.0.0.1:8000`    | 8000            | casbin/casdoor:latest       |
| `mailpit`  | `127.0.0.1:8025`    | 8025            | axllent/mailpit:latest      |
| `dozzle`   | `127.0.0.1:9999`    | 8080            | amir20/dozzle:v8            |
| `api`      | `127.0.0.1:8080`    | 8080            | `api/Dockerfile`            |
| `web`      | `127.0.0.1:3000`    | 3000            | `web/Dockerfile`            |

Volumes: `postgres_data` (named volume, external: true)

Rede: `app_network` (bridge)

## TLS/HTTPS

- Dashboard + API: nginx do host (`/etc/nginx/sites-available/fapitec`)
- Casdoor: nginx do host (`/etc/nginx/sites-available/auth.fapitec.duckdns.org`)
- Certificados Let's Encrypt em `/etc/letsencrypt/live/<dominio>/`
- Redirect HTTP (80) → HTTPS (443) com HSTS

## Configurações Nginx (Host)

### `/etc/nginx/sites-available/fapitec`
Proxy para dashboard (Next.js `:3000`), API Go (`:8080`), Dozzle logs (`:9999`/logs), Mailpit (`:8025`/mail), ACME challenge.

### `/etc/nginx/sites-available/auth.fapitec.duckdns.org`
Proxy para Casdoor (`:8000`), ACME challenge.

## Fluxo de Deploy

1. Desenvolvimento local → commit + push para GitHub
2. Na VM: `cd ~/apps/fapitec && git pull` (quando houver git) ou copiar artefatos
3. `docker compose up -d --build` para rebuildar serviços alterados
4. Se alterar nginx host: `sudo nginx -t && sudo systemctl reload nginx`

## Variáveis de Ambiente (`.env`)

| Variável | Descrição |
|----------|-----------|
| `DB_PASSWORD` | Senha do PostgreSQL |
| `DATABASE_URL` | URL de conexão com o banco |
| `AUTH_PROVIDER` | `internal` ou `casdoor` |
| `CASDOOR_ENDPOINT` | `https://auth.fapitec.duckdns.org` |
| `CASDOOR_CLIENT_ID` | Client ID do Casdoor |
| `CASDOOR_CLIENT_SECRET` | Client Secret do Casdoor |
| `CASDOOR_CERTIFICATE` | Certificado do Casdoor |
| `CASDOOR_ORGANIZATION_NAME` | `fapitec` |
| `FRONTEND_URL` | `https://fapitec.duckdns.org` |
| `SMTP_*` | Configurações de email |

## Observações

- `web` usa Next.js em modo servidor (Node), não static export
- API Go não faz TLS nativo — delega ao nginx do host
- Cookies JWT com `Secure: true` (exigem HTTPS)
- Banco PostgreSQL atende tanto a aplicação (`fapitec`) quanto Casdoor (`casdoor`), databases separadas no mesmo servidor
