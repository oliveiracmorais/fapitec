# Status do Projeto

## Goal
Concluir Epic 2.2 — Gestão de Editais e Submissão de Propostas.

## Progress
### Done
- **Epic 2.1 — Finalização da Homologação e Pipeline CI/CD — CONCLUÍDA**
- **Story 2.1.1 — Configuração do subdomínio Casdoor no nginx da VM:**
  - Nginx host configurado: `auth.fapitec.duckdns.org` → `127.0.0.1:8000` (sites-available + sites-enabled)
  - Certificado Let's Encrypt emitido e auto-renew configurado via webroot
  - UI Casdoor acessível via HTTPS
  - Fluxo OIDC completo validado (login → Casdoor → callback → JWT → dashboard)
  - ACME challenge adicionado ao config do auth para renovação automática
- **Story 2.1.2 — Testes de Resiliência do Sistema IAM:**
  - 8 novos testes automatizados de resiliência (410 Gone, JWT expiry, JWT inválido, token ausente)
  - 410 Gone handlers implementados para 6 endpoints legados no modo Casdoor
  - Relatório consolidado em `docs/qa/testes-resiliencia-iam.md`
  - 10/10 acceptance criteria validados (8 automatizados + 2 manuais documentados)
- **Story 2.1.3 — Pipeline CI/CD para Deploy em Homologação:**
  - Pipeline GitHub Actions criada: lint → typecheck → test → build → deploy → healthcheck → rollback
  - Secrets `SSH_HOST` e `SSH_KEY` configurados no GitHub
  - Pipeline executada com sucesso em push para `main` (run #27036212417)
  - Rollback automático testado (execuções anteriores com falha)
  - Script `scripts/deploy-vm.sh` criado para execução remota
- **Story 1.18 — Integração Casdoor IAM concluída** (Waves 1-9):
  - Infraestrutura local + homologação (Docker, init DB, feature flag)
  - Adapter Casdoor Go (JWT, OAuth, Enforce, criar usuário, GerarURL)
  - Middleware de autenticação (valida JWT, injeta claims no contexto, fallback internal)
  - Fluxo OIDC (GET /auth/login redirect, callback code→token, cookie JWT)
  - Frontend Next.js: middleware, API routes (proxy via backend), login button gov.br
  - Seed de 6 perfis (roles) + admin + permissões (8 módulos x 5 operações)
  - 14 testes automatizados (adapter + middleware authN + authZ) — todos PASS
- Infraestrutura base da homologação:
  - VM Oracle A1 Flex ARM64 configurada com Docker 29.5.2 + Compose v5.1.4
  - Nginx host com TLS Let's Encrypt para `fapitec.duckdns.org` e `auth.fapitec.duckdns.org`
  - Certificado SSL renovado e auto-renew funcional
  - Dockerfiles de homologação (API Go + Web Next.js)
  - `docker-compose.yml` na VM (PostgreSQL 16 + API Go + Web + Dozzle + Casdoor + Mailpit)
  - Migrações SQL automatizadas no deploy
  - Dozzle para logs em tempo real, Mailpit para email
  - Documentação do ambiente em `docs/infra/ambiente-homologacao.md`
- Cadastro e login funcionando via Casdoor (OIDC) e internal (fallback)
- Recuperação de senha testada

### In Progress
- *(none)*

### Blocked
- *(none)*

### Ready for Deploy
- *(none)*

## Próximos Passos
- **Epic 2.2: Gestão de Editais e Submissão de Propostas** — EM ANDAMENTO
  - Story 2.2.1: Extensão do Modelo de Editais — ✅ CONCLUÍDA (inclui correções C-01, C-02, C-03, C-04)
  - Story 2.2.2: Submissão de Propostas — Backend — ✅ CONCLUÍDA (QA fixes aplicados: versionamento e validação de edital)
  - Story 2.2.3: Submissão de Propostas — Frontend — A AGENDAR
  - Story 2.2.4: Gestão de Avaliadores e Atribuição — A AGENDAR
  - Story 2.2.5: Avaliação de Projetos — A AGENDAR
  - Story 2.2.6: Classificação, Seleção e Divulgação — A AGENDAR
- Módulos futuros: bolsistas, prestação de contas — Onda 2+

## Epic Ativa
- **Epic 2.2 — Gestão de Editais e Submissão de Propostas** — Em Andamento
