# ADR-001 - Estratégia de Identidade e Acesso

## Status

Aceito

## Contexto

O sistema exigirá autenticação, autorização por perfis e possivelmente regras mais finas de acesso entre módulos, operações e perfis institucionais.

Também há requisitos implícitos de:

- trilha de auditoria;
- segurança da informação;
- conformidade com contexto brasileiro e com LGPD;
- aderência às restrições operacionais da `FAPITEC-SE`;
- suporte a evolução futura sem acoplamento excessivo ao frontend.

O usuário demonstrou interesse em avaliar:

- `Casdoor` para autenticação e gestão de identidades;
- `Casbin` para autorização.

## Decisão

Adotar **Casdoor como plataforma única de IAM** (identidade, autenticação, autorização e auditoria).

### Motivação

A avaliação técnica revelou que **Casdoor já inclui Casbin internamente** para seu modelo de permissões, tornando a pilha Casdoor + Casbin standalone (Alternativa C original) redundante. A decisão consolida-se em Casdoor-only pelos seguintes fatores:

1. **Cobertura completa**: Casdoor atende authN (login, registro, MFA, recuperação de senha) + authZ (RBAC com suporte a ABAC futuro) em um único serviço.
2. **UI de gestão nativa**: O painel administrador do Casdoor permite gerenciar usuários, papéis e permissões sem construir interface própria.
3. **SSO-ready**: Suporte nativo a OAuth 2.0 / OIDC, facilitando integração futura com instituições parceiras.
4. **Auditoria embutida**: Record API registra todas as operações (login, criação de usuário, alteração de permissões).
5. **ARM64 compatível**: Imagem Docker multi-arquitetura (~68MB) desde v1.621.0, adequada para a VM Oracle Cloud.
6. **SDK Go maduro**: `casdoor-go-sdk` v1.46.0 com 134 releases, compatível com Go 1.26.3 do projeto.
7. **Separação de responsabilidades**: Autenticação via OIDC/JWT, autorização via Casdoor Enforce API — backend Go valida tokens e consulta decisões sem acoplamento.

### O que será mantido do backend Go atual

- **Validação de CPF/passaporte** (lógica de domínio pura, permanece em `identidade_e_acesso/dominio/objetos_de_valor/`)
- **Validação de força de senha** (serviço de domínio existente)
- **CAPTCHA Turnstile** (pode permanecer como middleware pré-Casdoor)
- **Adaptadores**: as interfaces de domínio existentes (`ServicoDeAutenticacao`, etc.) serão implementadas por um adapter Casdoor em `identidade_e_acesso/infraestrutura/autenticacao/adaptador_casdoor.go`

### O que o Casdoor substituirá

- `POST /api/v1/login` → fluxo OIDC Casdoor
- `POST /api/v1/cadastro` → Casdoor User Management API
- `POST /api/v1/solicitar-redefinicao-senha` → fluxo nativo Casdoor
- Hash de senha (bcrypt) → Casdoor gerencia internamente
- Sessão via `localStorage` → JWT/OIDC tokens gerenciados pelo Casdoor
- Autorização (inexistente hoje) → Casdoor Enforce API + modelo RBAC

## Consequências

- **Arquitetura**: o contexto `identidade_e_acesso` ganha um adapter para Casdoor nos níveis de infraestrutura; o domínio permanece independente da ferramenta.
- **Deploy**: Casdoor será adicionado como serviço no `docker-compose.homolog.yml`, com seu próprio banco PostgreSQL (pode ser no mesmo servidor).
- **Migração**: a autenticação própria atual será mantida como fallback durante a transição; o Casdoor será ativado por feature flag.
- **Custo de operação**: +1 container Docker (~68MB), +1 schema no PostgreSQL, necessidade de monitorar disponibilidade do Casdoor.
- **Perfis**: os 6 perfis institucionais mapeados na Story 1.4 serão cadastrados como roles no Casdoor.
- **Testes**: a camada de adapter deve ser coberta por testes de integração; o domínio segue testável sem dependência externa.

## Próximos passos

### Imediatos

1. ~~Levantar requisitos de perfis e permissões~~ ✓ (Story 1.4)
2. ~~Avaliar aderência de Casdoor e Casbin~~ ✓ (relatório em `/tmp/adr-001-casdoor-casbin-evaluation-report.md`)
3. **Criar story de implementação**: deploy do Casdoor + adapter Go + migração de fluxos
4. **Alinhar Story 1.4** (autorização por perfis) à decisão: perfis viram roles no Casdoor
5. **Alinhar Story 1.13** (ponta-a-ponta) para usar Casdoor em vez de auth própria

### Durante implementação

6. Deploy do Casdoor via Docker Compose na Oracle Cloud VM
7. Implementar `adaptador_casdoor.go` no backend Go
8. Implementar middleware de autorização usando Casdoor Enforce API
9. Seed dos 6 perfis institucionais no Casdoor
10. Migração gradual dos endpoints de auth própria para Casdoor
11. Remover código de auth própria após validação (opcional, manter como fallback)

## Referências

- Relatório completo: `/tmp/adr-001-casdoor-casbin-evaluation-report.md`
- Casdoor: https://github.com/casdoor/casdoor
- Go SDK: https://github.com/casdoor/casdoor-go-sdk
- Docker: `casbin/casdoor` — multi-arch (amd64 + arm64)
- Próxima story: `docs/stories/` (a criar)

## Change Log

| Data | Versão | Descrição | Autor |
|------|--------|-----------|-------|
| 2026-06-01 | 1.0 | Proposta inicial | — |
| 2026-06-01 | 2.0 | Decisão final: Casdoor-only. ADR aceito após avaliação técnica de Casdoor + Casbin. Relatório em `/tmp/adr-001-casdoor-casbin-evaluation-report.md` | Pax (PO) |
