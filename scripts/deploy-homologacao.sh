#!/usr/bin/env bash
set -euo pipefail

# Uso: ./scripts/deploy-homologacao.sh
# Requer: chave SSH adicionada ao ssh-agent (ssh-add ~/.ssh/sua-chave.key)

SSH_USER="${SSH_USER:-ubuntu}"
VM_IP="${VM_IP:-137.131.166.115}"
REMOTE_DIR="${REMOTE_DIR:-/home/ubuntu/apps/fapitec}"
REMOTE="$SSH_USER@$VM_IP"

echo "=== Verificando ssh-agent ==="
ssh -o BatchMode=yes "$REMOTE" "echo Conexao SSH OK" || {
	echo "ERRO: Conexao SSH falhou. Certifique-se de que a chave esta carregada no agente:"
	echo "  ssh-add ~/.ssh/sua-chave.key"
	exit 1
}

echo "=== Build local dos artefatos ==="

echo "[1/4] Compilando API Go para linux/arm64..."
cd "$(dirname "$0")/.."
mkdir -p apps/api/bin
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -C apps/api -o bin/api-linux-arm64 ./cmd/api

echo "[2/4] Buildando Next.js (standalone)..."
cd apps/web
NEXT_OUTPUT=standalone npm run build
cd ../..

echo "[3/4] Copiando artefatos para VM..."
ssh "$REMOTE" "mkdir -p $REMOTE_DIR/api $REMOTE_DIR/web $REMOTE_DIR/api/db/migracoes"

scp apps/api/bin/api-linux-arm64 "$REMOTE:$REMOTE_DIR/api/api"
scp apps/api/Dockerfile.homolog "$REMOTE:$REMOTE_DIR/api/Dockerfile"

rsync -az --delete apps/web/.next/standalone/ "$REMOTE:$REMOTE_DIR/web/"
rsync -az --delete apps/web/.next/static/ "$REMOTE:$REMOTE_DIR/web/apps/web/.next/static/"
rsync -az --delete apps/web/public/ "$REMOTE:$REMOTE_DIR/web/apps/web/public/"
scp apps/web/Dockerfile.homolog "$REMOTE:$REMOTE_DIR/web/Dockerfile"

scp docker-compose.homolog.yml "$REMOTE:$REMOTE_DIR/docker-compose.yml"
scp .env.homolog "$REMOTE:$REMOTE_DIR/.env" 2>/dev/null || true
rsync -az --delete apps/api/db/migracoes/ "$REMOTE:$REMOTE_DIR/api/db/migracoes/"
scp apps/api/db/init-casdoor-db.sql "$REMOTE:$REMOTE_DIR/api/db/init-casdoor-db.sql"

echo "[4/4] Realizando deploy e migracoes na VM..."
ssh "$REMOTE" "cd $REMOTE_DIR && docker compose pull && docker compose up -d --build && echo 'Aguardando PostgreSQL...' && until docker compose exec -T postgres pg_isready -U fapitec 2>/dev/null; do sleep 1; done && echo 'Executando migracoes...' && for f in api/db/migracoes/*.sql; do echo \"  Aplicando \$f\" && cat \"\$f\" | docker compose exec -T postgres psql -U fapitec -d fapitec; done && echo 'Migracoes concluidas'"

echo "=== Deploy concluido ==="
echo "Acesse: https://compusa.duckdns.org"
