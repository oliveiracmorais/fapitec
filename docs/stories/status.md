# Status do Projeto

## Goal
Preparar ambiente de homologação na Oracle Cloud VM e realizar deploy funcional da aplicação FAPITEC-SE.

## Progress
### Done
- VM Oracle A1 Flex ARM64 configurada com Docker 29.5.2 + Compose v5.1.4
- Certificado SSL Let's Encrypt emitido para fapitec.duckdns.org (válido até 02/09/2026)
- Nginx configurado como reverse proxy: `/api/` → API Go (127.0.0.1:8080), demais rotas → Next.js (127.0.0.1:3000)
- Dockerfiles de homologação criados (API e Web)
- `docker-compose.homolog.yml` criado (PostgreSQL 16 + API Go + Web + Dozzle + Casdoor)
- Script `scripts/deploy-homologacao.sh` criado (inclui cópia e execução de migrações SQL + init-casdoor-db)
- Correção de `http.Error()` → `jsonError()` para `Content-Type: application/json`
- Correção de caminhos de assets estáticos do Next.js standalone
- Migrações SQL automatizadas no deploy
- Dozzle instalado para logs em tempo real (`/logs/`)
- Estrutura reorganizada para `~/apps/fapitec/` na VM
- **Story 1.18 — Integração Casdoor IAM concluída** (Waves 1-9):
  - Infraestrutura local + homologação (Docker, init DB, feature flag)
  - Adapter Casdoor Go (JWT, OAuth, Enforce, criar usuário, GerarURL)
  - Middleware de autenticação (valida JWT, injeta claims no contexto)
  - Middleware de autorização (path-based route permissions: editais, dashboard, auditoria)
  - Fluxo OIDC (GET /auth/login redirect, POST /auth/callback code→token)
  - Endpoints antigos desabilitados com 410 Gone quando AUTH_PROVIDER=casdoor
  - Frontend Next.js: middleware, API routes, AuthContext, login button
  - Seed de 6 perfis (roles) + admin + permissões (8 módulos x 5 operações)
  - 14 testes automatizados (adapter + middleware authN + authZ) — todos PASS
  - go vet/build PASS, tsc --noEmit PASS
- **Cadastro e login funcionando** (CPF: `123.456.789-09`) — modo `internal`
- **Recuperação de senha testada** (solicitar token + redefinir senha + login)

### In Progress
- *(none)*

### Blocked
- *(none)*

### Ready for Deploy (Story 1.18 — Homologação)
- Casdoor configurado no `docker-compose.homolog.yml` (porta 8000, database `casdoor`)
- Init script `api/db/init-casdoor-db.sql` copiado no deploy
- `.env.homolog` com vars `AUTH_PROVIDER`, `CASDOOR_*`
- Nginx config snippet em `docs/infra/nginx-casdoor-subdomain.conf`
- **Pending VM-side**: configurar nginx host para `auth.fapitec.duckdns.org` → `127.0.0.1:8000`

## Próximos Passos
- Módulos do negócio (bolsistas, prestação de contas) — Onda 2+

## Epic Ativa
- [Epic 2.1: Finalização da Homologação e Pipeline CI/CD](../epics/2.1.finalizacao-homologacao-e-ci-cd.md)
  - Story 1: Configuração do subdomínio Casdoor no nginx da VM
  - Story 2: Testes de resiliência do sistema IAM
  - Story 3: Pipeline CI/CD para deploy em homologação
