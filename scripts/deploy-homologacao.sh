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

echo "[1/3] Buildando Next.js (standalone)..."
cd apps/web
NEXT_OUTPUT=standalone npm run build
cd ../..

echo "[2/3] Copiando artefatos para VM..."
ssh "$REMOTE" "mkdir -p $REMOTE_DIR/api $REMOTE_DIR/web $REMOTE_DIR/api/db/migracoes"

rsync -az --delete apps/api/ "$REMOTE:$REMOTE_DIR/api/" --exclude bin/

rsync -azL --delete apps/web/.next/standalone/apps/web/ "$REMOTE:$REMOTE_DIR/web/"
rsync -az --delete apps/web/.next/static/ "$REMOTE:$REMOTE_DIR/web/.next/static/"
rsync -az --delete apps/web/public/ "$REMOTE:$REMOTE_DIR/web/public/"
scp apps/web/Dockerfile.homolog "$REMOTE:$REMOTE_DIR/web/Dockerfile"

scp docker-compose.homolog.yml "$REMOTE:$REMOTE_DIR/docker-compose.yml"
scp .env.homolog "$REMOTE:$REMOTE_DIR/.env" 2>/dev/null || true
rsync -az --delete apps/api/db/migracoes/ "$REMOTE:$REMOTE_DIR/api/db/migracoes/"
scp apps/api/db/init-casdoor-db.sql "$REMOTE:$REMOTE_DIR/api/db/init-casdoor-db.sql"

echo "[3/3] Realizando deploy e migracoes na VM..."
ssh "$REMOTE" "cd $REMOTE_DIR && docker compose pull && docker compose up -d --build && echo 'Aguardando PostgreSQL...' && until docker compose exec -T postgres pg_isready -U fapitec 2>/dev/null; do sleep 1; done && echo 'Executando migracoes...' && for f in api/db/migracoes/*.sql; do echo \"  Aplicando \$f\" && cat \"\$f\" | docker compose exec -T postgres psql -U fapitec -d fapitec; done && echo 'Migracoes concluidas'"

echo "=== Deploy concluido ==="
echo "Acesse: https://fapitec.duckdns.org"
