# GitHub Actions - Workflows

Este diretório contém os workflows de CI/CD do projeto FAPITEC-SE.

## Workflows

### `deploy-homologacao.yml`

Pipeline de deploy automatizado para o ambiente de homologação na Oracle Cloud VM.

**Trigger:** Push na branch `main` ou manual via `workflow_dispatch`.

**Jobs:**

| Job | Descrição | Dependente de |
|-----|-----------|---------------|
| `quality` | Lint, typecheck, testes Node.js e Go | — |
| `build-and-deploy` | Build, cópia de artefatos, deploy SSH, healthcheck | `quality` |
| `rollback` | Rollback automático se healthcheck falhar | `build-and-deploy` (failure) |

**Fluxo:**

1. **quality** — Executa `pnpm run lint`, `pnpm run typecheck`, `pnpm test`, `go test ./...`
2. **build-and-deploy** — Faz checkout, build do Next.js (standalone) e API Go, copia artefatos via rsync/SSH, executa `scripts/deploy-homologacao.sh` na VM, valida healthcheck (`/api/v1/health`)
3. **rollback** — Se `build-and-deploy` falhar, executa `scripts/rollback-homologacao.sh` na VM e revalida saúde

## Secrets Necessários

Configurar em **Settings → Secrets and variables → Actions**:

| Secret | Descrição | Obrigatório |
|--------|-----------|-------------|
| `SSH_HOST` | IP da Oracle Cloud VM (137.131.166.115) | Sim |
| `SSH_KEY` | Chave privada SSH para acesso à VM | Sim |
| `DOCKER_COMPOSE_VARS` | Variáveis de ambiente para deploy (ex.: DB_PASSWORD) | Sim |
| `SLACK_WEBHOOK` | Webhook do Slack para notificações | Não |

## Deploy Manual

O deploy manual via script `scripts/deploy-homologacao.sh` continua funcionando independentemente da pipeline.
