# Pipeline CI/CD — Deploy Homologação

## Visão Geral

Pipeline automatizada via GitHub Actions que executa validações de qualidade, build e deploy no ambiente de homologação (Oracle Cloud VM) sempre que houver push na branch `main`.

## Arquitetura

```
push → main
  │
  ├── Job: quality
  │   ├── lint (pnpm run lint)
  │   ├── typecheck (pnpm run typecheck)
  │   ├── test Node.js (pnpm test)
  │   └── test Go (go test ./...)
  │
  ├── Job: build-and-deploy (após quality)
  │   ├── build Next.js (standalone)
  │   ├── build API Go
  │   ├── copy artifacts via rsync/SSH
  │   ├── executa deploy-homologacao.sh na VM
  │   ├── healthcheck API (/api/v1/health)
  │   └── healthcheck Frontend (/)
  │
  └── Job: rollback (se build-and-deploy falhar)
      ├── executa rollback-homologacao.sh na VM
      └── verifica saúde pós-rollback
```

## Arquivos

| Arquivo | Descrição |
|---------|-----------|
| `.github/workflows/deploy-homologacao.yml` | Definição do pipeline GitHub Actions |
| `.github/workflows/README.md` | Documentação dos workflows e secrets |
| `scripts/deploy-homologacao.sh` | Script de deploy executado na VM |
| `scripts/rollback-homologacao.sh` | Script de rollback (restaura versão anterior) |

## Secrets do GitHub

Configurar em **Settings → Secrets and variables → Actions**:

| Secret | Descrição |
|--------|-----------|
| `SSH_KEY` | Chave privada para acesso SSH à VM |
| `SLACK_WEBHOOK` | (Opcional) Webhook do Slack para notificações |

> **Nota:** `SSH_HOST` e `DOCKER_COMPOSE_VARS` estão definidos como variáveis `env` no workflow. Podem ser movidos para secrets se necessário.

## Estratégia de Rollback

1. Healthcheck da API (`/api/v1/health`) é verificado após o deploy
2. Se falhar por mais de 60s (12 tentativas × 5s), o job `rollback` é acionado
3. O rollback para os containers atuais e tenta restaurar imagens com tag `previous`
4. Se não houver imagem `previous`, usa `latest` como fallback
5. Pós-rollback, healthcheck é revalidado
6. Notificação de falha é enviada (se Slack configurado)

## Testando a Pipeline

### Em branch separada
```bash
git checkout -b test/deploy-pipeline
# fazer alterações
git push origin test/deploy-pipeline
# O workflow pode ser acionado manualmente via GitHub UI (workflow_dispatch)
```

### Healthcheck com serviço quebrado
Para testar o rollback, faça uma alteração que quebre intencionalmente o healthcheck da API e faça push para main. O pipeline deve:
1. Tentar o deploy
2. Falhar no healthcheck
3. Executar rollback automaticamente
4. Enviar notificação de falha

## Manutenção

### Adicionar novo passo ao pipeline
Editar `.github/workflows/deploy-homologacao.yml` e adicionar o step no job apropriado.

### Atualizar script de deploy
Modificar `scripts/deploy-homologacao.sh` — o pipeline copia este script para a VM antes de executá-lo.

### Atualizar script de rollback
Modificar `scripts/rollback-homologacao.sh` — o pipeline copia este script para a VM antes de executá-lo em caso de falha.
