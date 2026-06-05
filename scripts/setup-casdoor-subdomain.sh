#!/usr/bin/env bash
set -euo pipefail

# =============================================================================
# setup-casdoor-subdomain.sh
# Configura o subdomínio auth.fapitec.duckdns.org no nginx da VM Oracle Cloud.
#
# Uso:
#   sudo ./scripts/setup-casdoor-subdomain.sh
#
# Requer:
#   - Execução como root (sudo)
#   - Portas 80 e 443 liberadas no firewall da VM e na Oracle Cloud Security List
#   - duckdns.org DNS apontando auth.fapitec.duckdns.org → IP da VM
# =============================================================================

# --- Configurações ---
DOMAIN="auth.fapitec.duckdns.org"
NGINX_CONF_SOURCE="$(dirname "$0")/../docs/infra/nginx-casdoor-subdomain.conf"
NGINX_AVAILABLE="/etc/nginx/sites-available/${DOMAIN}"
NGINX_ENABLED="/etc/nginx/sites-enabled/${DOMAIN}"
LETSENCRYPT_DIR="/etc/letsencrypt/live/${DOMAIN}"
CERTBOT_EMAIL="admin@fapitec.se.gov.br"

echo "=============================================="
echo "  Setup: ${DOMAIN}"
echo "=============================================="

# --- Verificação de pré-requisitos ---
if [[ $EUID -ne 0 ]]; then
    echo "ERRO: Execute como root (sudo)."
    exit 1
fi

if [[ ! -f "$NGINX_CONF_SOURCE" ]]; then
    echo "ERRO: ${NGINX_CONF_SOURCE} não encontrado."
    exit 1
fi

if ! command -v nginx &>/dev/null; then
    echo "ERRO: nginx não instalado."
    exit 1
fi

if ! command -v certbot &>/dev/null; then
    echo "certbot não encontrado. Instalando..."
    apt-get update -qq && apt-get install -y -qq certbot python3-certbot-nginx
fi

# --- 1. Copiar configuração nginx ---
echo ""
echo "[1/5] Copiando configuração nginx..."
cp "$NGINX_CONF_SOURCE" "$NGINX_AVAILABLE"
echo "  → ${NGINX_AVAILABLE}"

# --- 2. Ativar site via symlink ---
echo "[2/5] Ativando site..."
ln -sf "$NGINX_AVAILABLE" "$NGINX_ENABLED"
echo "  → ${NGINX_ENABLED}"

# --- 3. Emitir/renovar certificado SSL ---
echo "[3/5] Certificado SSL..."

if [[ -d "$LETSENCRYPT_DIR" ]]; then
    echo "  Certificado já existe. Verificando validade..."
    if openssl x509 -checkend 2592000 -noout -in "${LETSENCRYPT_DIR}/fullchain.pem" 2>/dev/null; then
        echo "  Certificado válido por mais de 30 dias. Nenhuma ação necessária."
    else
        echo "  Certificado próximo de expirar. Renovando..."
        certbot certonly --nginx -d "$DOMAIN" --non-interactive --agree-tos --email "$CERTBOT_EMAIL"
    fi
else
    echo "  Emitindo novo certificado para ${DOMAIN}..."
    certbot certonly --nginx -d "$DOMAIN" --non-interactive --agree-tos --email "$CERTBOT_EMAIL"
fi

# --- 4. Validar sintaxe e recarregar nginx ---
echo "[4/5] Validando e recarregando nginx..."
nginx -t
nginx -s reload
echo "  nginx recarregado com sucesso."

# --- 5. Configurar auto-renew via cron (se não existir) ---
echo "[5/5] Auto-renew do certificado..."
if crontab -l 2>/dev/null | grep -q "certbot renew"; then
    echo "  Cron job para certbot renew já existe."
else
    (crontab -l 2>/dev/null; echo "0 3 * * * /usr/bin/certbot renew --quiet --post-hook 'systemctl reload nginx'") | crontab -
    echo "  Cron job adicionado: diariamente às 03:00."
fi

echo ""
echo "=============================================="
echo "  Setup concluído!"
echo "  Acesse: https://${DOMAIN}"
echo "=============================================="

# Teste rápido
echo ""
echo "--- Teste rápido ---"
curl -sI "https://${DOMAIN}" | head -5 || echo "  (Aguardando propagação DNS / serviço Casdoor)"
