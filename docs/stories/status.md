# Status do Projeto

## Goal
Preparar ambiente de homologação na Oracle Cloud VM e realizar deploy funcional da aplicação FAPITEC-SE.

## Progress
### Done
- VM Oracle A1 Flex ARM64 configurada com Docker 29.5.2 + Compose v5.1.4
- Certificado SSL Let's Encrypt emitido para compusa.duckdns.org (válido até 26/08/2026)
- Nginx configurado como reverse proxy: `/api/` → API Go (127.0.0.1:8080), demais rotas → Next.js (127.0.0.1:3000)
- Dockerfiles de homologação criados (API e Web)
- `docker-compose.homolog.yml` criado (PostgreSQL 16 + API Go + Next.js)
- Script `scripts/deploy-homologacao.sh` criado
- Correção de `http.Error()` → `jsonError()` para `Content-Type: application/json`
- Correção de caminhos de assets estáticos do Next.js standalone
- Migrações SQL executadas (usuários, tokens, editais)
- **Cadastro e login funcionando** (CPF: `123.456.789-09`, senha: `Teste@123`)

### In Progress
- *(none)*

### Blocked
- *(none)*

## Credenciais de Teste
- CPF: `123.456.789-09` (ou `12345678909`)
- Senha: `Teste@123`

## Próximos Passos
- Automatizar migrações SQL no script de deploy
- Testar outras funcionalidades (recuperação de senha, editais)
