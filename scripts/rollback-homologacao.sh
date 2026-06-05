#!/usr/bin/env bash
set -euo pipefail

# =============================================================================
# rollback-homologacao.sh
# Reverte o deploy para a versão anterior dos containers Docker.
#
# Uso:
#   bash scripts/rollback-homologacao.sh
#
# Requer:
#   - Execução na própria VM (ou via SSH)
#   - Docker Compose com imagens taggeadas por versão
#   - Pelo menos uma imagem anterior disponível com tag 'previous'
# =============================================================================

REMOTE_DIR="${REMOTE_DIR:-/home/ubuntu/apps/fapitec}"
COMPOSE_FILE="${REMOTE_DIR}/docker-compose.yml"
PROJECT_NAME="fapitec"

echo "=============================================="
echo "  Rollback Homologação"
echo "=============================================="

cd "$REMOTE_DIR"

# --- 1. Parar containers atuais ---
echo "[1/4] Parando containers atuais..."
docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down --timeout 30 || true

# --- 2. Restaurar imagens anteriores (se existirem) ---
echo "[2/4] Restaurando imagens anteriores..."

# Estratégia: usar imagens com tag 'previous' para cada serviço
# Se não houver tag 'previous', tenta 'latest' como fallback
SERVICES=$(docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" config --services 2>/dev/null || echo "api web")

for service in $SERVICES; do
    IMAGE=$(docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" images "$service" 2>/dev/null | awk "NR>1 && /\$service/ {print \$2}" || true)

    if docker image inspect "${PROJECT_NAME}_${service}:previous" &>/dev/null; then
        echo "  Restaurando ${service}: previous → current"
        docker tag "${PROJECT_NAME}_${service}:previous" "${PROJECT_NAME}_${service}:current" 2>/dev/null || true
    elif docker image inspect "${PROJECT_NAME}_${service}:latest" &>/dev/null; then
        echo "  Usando latest para ${service} (previous não disponível)"
    else
        echo "  Aviso: Nenhuma imagem anterior encontrada para ${service}"
    fi
done

# --- 3. Reiniciar containers com versão anterior ---
echo "[3/4] Reiniciando containers..."
docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" up -d

# --- 4. Aguardar saúde ---
echo "[4/4] Aguardando serviços ficarem saudáveis..."
for i in $(seq 1 12); do
    HEALTH_API=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/api/v1/health 2>/dev/null || true)
    HEALTH_WEB=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:3000/ 2>/dev/null || true)

    echo "  Tentativa $i: API=$HEALTH_API Web=$HEALTH_WEB"

    if [ "$HEALTH_API" = "200" ]; then
        echo "  API saudável!"
        echo ""
        echo "=============================================="
        echo "  Rollback concluído com sucesso!"
        echo "=============================================="
        exit 0
    fi
    sleep 5
done

echo ""
echo "ALERTA: Rollback executado, mas API não respondeu."
echo "Verifique manualmente: docker compose logs"
exit 1
