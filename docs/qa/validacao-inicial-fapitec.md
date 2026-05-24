# Validação Inicial FAPITEC-SE

## Objetivo

Este documento orienta a apresentação inicial da fundação da plataforma para validação pela `FAPITEC-SE`.

A proposta desta etapa é demonstrar um ponto de partida verificável, não uma solução completa de produção. O foco é confirmar macroestrutura, linguagem institucional, módulos, perfis mínimos e trilha auditável básica antes de aprofundar fluxos de negócio.

## Escopo Da Demonstração

- Catálogo inicial dos módulos previstos para a plataforma.
- Perfis institucionais mínimos para validação da autorização.
- Comandos CLI para listar módulos e perfis.
- Comandos CLI para listar identidades demonstrativas, simular sessão, acesso permitido, acesso negado e evento de auditoria.
- Dashboard modular inicial com páginas placeholder.

## Narrativa Recomendada

1. A plataforma está sendo construída a partir das stories fundacionais.
2. O alicerce respeita a sequência `CLI First -> Observability Second -> UI Third`.
3. Autenticação, autorização e auditoria aparecem nesta etapa como base mínima para validação.
4. O dashboard é um mapa navegável da macroestrutura, sem detalhar regras de negócio ainda não validadas.
5. A validação esperada da `FAPITEC-SE` deve focar nomes dos módulos, perfis, prioridades e aderência institucional.

## Roteiro CLI

```bash
pnpm run shell:modulos
pnpm run shell:perfis
pnpm run shell:identidades
node scripts/plataforma-shell.js sessao gestor_validacao
node scripts/plataforma-shell.js acesso-sessao sessao-gestor_validacao gestao_de_editais
node scripts/plataforma-shell.js auditoria-sessao sessao-gestor_validacao gestao_de_editais
node scripts/plataforma-shell.js auditoria-sessao sessao-proponente_validacao auditoria
```

## Roteiro UI

```bash
pnpm dev
```

Acessar:

```text
http://localhost:3000
```

Sessão demonstrativa direta:

```text
http://localhost:3000/?sessao=sessao-gestor_validacao
```

## Pontos A Validar Com A FAPITEC-SE

- Os módulos apresentados correspondem à macroestrutura esperada?
- Os nomes dos módulos estão adequados à linguagem institucional?
- Os perfis mínimos fazem sentido para a primeira onda?
- Algum perfil ou módulo deve ser renomeado, removido ou priorizado?
- A trilha mínima de auditoria cobre os eventos esperados para validação inicial?
- A próxima etapa deve priorizar integração runtime de autenticação/autorização/auditoria ou gestão documental base?

## Observações QA

- A Story 1.6 permanece com gate `CONCERNS` para uso como entrega runtime completa.
- A Story 1.7 adiciona sessão demonstrativa, autorização por perfil no dashboard e eventos de auditoria inspecionáveis.
- A autenticação demonstrativa segue como base de validação, sem fechar a decisão arquitetural sobre `Casdoor` e `Casbin`.
