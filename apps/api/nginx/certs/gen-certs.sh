#!/bin/sh
# Gera certificado auto-assinado para ambiente de homologacao/desenvolvimento
# Uso: ./gen-certs.sh
# Producao: substituir cert.pem e key.pem por certificado ICP-Brasil ou Let's Encrypt

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout key.pem \
  -out cert.pem \
  -subj "/C=BR/ST=SE/L=Aracaju/O=FAPITEC-SE/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,DNS:api.fapitec.se.gov.br"

echo "Certificados gerados: cert.pem, key.pem"
