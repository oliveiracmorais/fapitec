#!/usr/bin/env bash
set -euo pipefail

SSH_KEY="${SSH_KEY:-$HOME/.ssh/ssh-key-2026-05-26.key}"
SSH_USER="${SSH_USER:-ubuntu}"
VM_IP="${VM_IP:-137.131.166.115}"
REMOTE_DIR="${REMOTE_DIR:-/home/ubuntu/apps/fapitec}"

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
REMOTE="$SSH_USER@$VM_IP"

ssh -i "$SSH_KEY" "$REMOTE" "mkdir -p $REMOTE_DIR/api $REMOTE_DIR/web $REMOTE_DIR/api/db/migracoes"

scp -i "$SSH_KEY" apps/api/bin/api-linux-arm64 "$REMOTE:$REMOTE_DIR/api/api"
scp -i "$SSH_KEY" apps/api/Dockerfile.homolog "$REMOTE:$REMOTE_DIR/api/Dockerfile"

rsync -az --delete -e "ssh -i $SSH_KEY" apps/web/.next/standalone/ "$REMOTE:$REMOTE_DIR/web/"
rsync -az --delete -e "ssh -i $SSH_KEY" apps/web/.next/static/ "$REMOTE:$REMOTE_DIR/web/apps/web/.next/static/"
rsync -az --delete -e "ssh -i $SSH_KEY" apps/web/public/ "$REMOTE:$REMOTE_DIR/web/apps/web/public/"
# scp AFTER rsync para evitar que --delete remova os arquivos
scp -i "$SSH_KEY" apps/web/Dockerfile.homolog "$REMOTE:$REMOTE_DIR/web/Dockerfile"

scp -i "$SSH_KEY" docker-compose.homolog.yml "$REMOTE:$REMOTE_DIR/docker-compose.yml"
scp -i "$SSH_KEY" .env.homolog "$REMOTE:$REMOTE_DIR/.env" 2>/dev/null || true
rsync -az --delete -e "ssh -i $SSH_KEY" apps/api/db/migracoes/ "$REMOTE:$REMOTE_DIR/api/db/migracoes/"
scp -i "$SSH_KEY" apps/api/db/init-casdoor-db.sql "$REMOTE:$REMOTE_DIR/api/db/init-casdoor-db.sql"

echo "[4/4] Realizando deploy e migracoes na VM..."
ssh -i "$SSH_KEY" "$REMOTE" "cd $REMOTE_DIR && docker compose pull && docker compose up -d --build && echo 'Aguardando PostgreSQL...' && until docker compose exec -T postgres pg_isready -U fapitec 2>/dev/null; do sleep 1; done && echo 'Executando migracoes...' && for f in api/db/migracoes/*.sql; do echo \"  Aplicando \$f\" && cat \"\$f\" | docker compose exec -T postgres psql -U fapitec -d fapitec; done && echo 'Migracoes concluidas'"

echo "=== Deploy concluido ==="
echo "Acesse: https://compusa.duckdns.org"
