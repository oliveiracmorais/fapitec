#!/usr/bin/env bash
set -euo pipefail

# Uso: bash scripts/deploy-vm.sh
# Executado na VM pelo pipeline CI/CD (GitHub Actions).
# Pressupoe que artefatos ja foram copiados pelo job anterior.

REMOTE_DIR="${REMOTE_DIR:-/home/ubuntu/apps/fapitec}"

cd "$REMOTE_DIR"

echo "[1/2] Rebuildando e reiniciando containers..."
docker compose up -d --build

echo "[2/2] Executando migracoes..."
until docker compose exec -T postgres pg_isready -U fapitec 2>/dev/null; do
    sleep 1
done

for f in api/db/migracoes/*.sql; do
    echo "  Aplicando $(basename "$f")"
    cat "$f" | docker compose exec -T postgres psql -U fapitec -d fapitec
done

echo "=== Deploy na VM concluido ==="
