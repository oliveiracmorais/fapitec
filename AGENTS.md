# AGENTS.md - Synkra AIOX (Codex CLI)

Este arquivo define as instrucoes do projeto para o Codex CLI.

<!-- AIOX-MANAGED-START: core -->
## Core Rules

1. Siga a Constitution em `.aiox-core/constitution.md`
2. Priorize `CLI First -> Observability Second -> UI Third`
3. Trabalhe por stories em `docs/stories/`
4. Nao invente requisitos fora dos artefatos existentes
<!-- AIOX-MANAGED-END: core -->

<!-- AIOX-MANAGED-START: quality -->
## Quality Gates

- Rode `npm run lint`
- Rode `npm run typecheck`
- Rode `npm test`
- Atualize checklist e file list da story antes de concluir
<!-- AIOX-MANAGED-END: quality -->

## Project Standards

- Organizacao-alvo: `FAPITEC-SE`
- Gerenciador de pacotes padrao: `pnpm`
- Frontend padrao: `Next.js + React + Tailwind CSS + shadcn/ui`
- Backend padrao: `Go`
- Banco de dados padrao: `PostgreSQL`
- Acesso a dados no backend Go: preferir `sqlc`; nao usar ORM como padrao do projeto
- Estilo arquitetural: `DDD + Clean Architecture + Clean Code + TDD`
- Linguagem ubíqua: nomes de pastas, arquivos, modulos, tipos, atributos e casos de uso em portugues
- Localizacao: moeda, percentuais, datas, horarios e mensagens devem seguir padroes do Brasil
- Valores monetarios em Go: nunca usar `float32` ou `float64`
- Regras monetarias: preferir representacao inteira em centavos ou tipo de precisao controlada encapsulado no dominio
- Identidade e acesso: avaliar `Casdoor` e `Casbin` conforme ADR em `docs/architecture/adr-001-identidade-e-acesso.md`
- Infraestrutura: respeitar o local de hospedagem exigido no Termo de Referencia
- Planejamento: considerar o cronograma mensal de entregas exigido no Termo de Referencia

## Domain Conventions

- Usar termos do dominio da FAPITEC-SE como `edital`, `proponente`, `bolsista`, `prestacao_de_contas` e `tomada_de_contas_especial`
- Evitar nomes genericos em ingles no dominio quando houver termo institucional claro em portugues
- Regras de negocio nao devem depender diretamente de framework, banco, HTTP ou UI
- Componentes de interface nao devem concentrar logica critica de dominio

## Backend Guidance

- Estruturar o backend Go por contextos de dominio
- Separar `dominio`, `aplicacao`, `interfaces` e `infraestrutura` sempre que fizer sentido
- Centralizar SQL em consultas versionadas e geradas via `sqlc`
- Priorizar testes de dominio e casos de uso antes de adapters e integracoes
- Tratar autenticacao, autorizacao, auditoria e gestao documental como capacidades transversais

## Planning Notes

- O sistema esta sendo planejado a partir dos modulos definidos pela FAPITEC-SE
- O backlog deve respeitar dependencias entre capacidades fundacionais e modulos de negocio
- Stories devem considerar incrementos mensais compativeis com o cronograma contratual
- Decisoes de hospedagem, seguranca e integracao devem manter aderencia ao Termo de Referencia

<!-- AIOX-MANAGED-START: codebase -->
## Project Map

- Core framework: `.aiox-core/`
- CLI entrypoints: `bin/`
- Shared packages: `packages/`
- Tests: `tests/`
- Docs: `docs/`
<!-- AIOX-MANAGED-END: codebase -->

<!-- AIOX-MANAGED-START: commands -->
## Common Commands

- `npm run sync:ide`
- `npm run sync:ide:check`
- `npm run sync:skills:codex`
- `npm run sync:skills:codex:global` (opcional; neste repo o padrao e local-first)
- `npm run validate:structure`
- `npm run validate:agents`
<!-- AIOX-MANAGED-END: commands -->

<!-- AIOX-MANAGED-START: shortcuts -->
## Agent Shortcuts

Preferencia de ativacao no Codex CLI:
1. Use `/skills` e selecione `aiox-<agent-id>` vindo de `.codex/skills` (ex.: `aiox-architect`)
2. Se preferir, use os atalhos abaixo (`@architect`, `/architect`, etc.)

Interprete os atalhos abaixo carregando o arquivo correspondente em `.aiox-core/development/agents/` (fallback: `.codex/agents/`), renderize o greeting via `generate-greeting.js` e assuma a persona ate `*exit`:

- `@architect`, `/architect`, `/architect.md` -> `.aiox-core/development/agents/architect.md`
- `@dev`, `/dev`, `/dev.md` -> `.aiox-core/development/agents/dev.md`
- `@qa`, `/qa`, `/qa.md` -> `.aiox-core/development/agents/qa.md`
- `@pm`, `/pm`, `/pm.md` -> `.aiox-core/development/agents/pm.md`
- `@po`, `/po`, `/po.md` -> `.aiox-core/development/agents/po.md`
- `@sm`, `/sm`, `/sm.md` -> `.aiox-core/development/agents/sm.md`
- `@analyst`, `/analyst`, `/analyst.md` -> `.aiox-core/development/agents/analyst.md`
- `@devops`, `/devops`, `/devops.md` -> `.aiox-core/development/agents/devops.md`
- `@data-engineer`, `/data-engineer`, `/data-engineer.md` -> `.aiox-core/development/agents/data-engineer.md`
- `@ux-design-expert`, `/ux-design-expert`, `/ux-design-expert.md` -> `.aiox-core/development/agents/ux-design-expert.md`
- `@squad-creator`, `/squad-creator`, `/squad-creator.md` -> `.aiox-core/development/agents/squad-creator.md`
- `@aiox-master`, `/aiox-master`, `/aiox-master.md` -> `.aiox-core/development/agents/aiox-master.md`
<!-- AIOX-MANAGED-END: shortcuts -->
