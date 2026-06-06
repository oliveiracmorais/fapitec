# QA Fix Request — Story 2.2.2 (Submissão de Propostas — Backend)

**De:** Quinn (QA)
**Para:** Dex (Dev)
**Data:** 2026-06-06
**Gate:** PASS with Recommendations 🟡

---

## Issue #1 — Histórico de Versões (AC #6)

**Severidade:** MÉDIA
**Tipo:** Funcional (gap de requisito)

### Problema
O Acceptance Criteria #6 especifica: "O sistema deve manter histórico de versões das alterações." Atualmente, `Atualizar` no repositório sobrescreve os dados da proposta sem preservar versões anteriores.

### Solução sugerida
Criar mecanismo de versionamento. Duas abordagens possíveis:

1. **Tabela `versoes_proposta`** (recomendada):
   - Nova tabela no banco que armazena snapshots da proposta a cada atualização
   - Acionada via trigger SQL ou na camada de aplicação antes de cada `Atualizar`
   - Permite rastrear todo o histórico de alterações

2. **Abordagem simplificada** (mínimo para AC):
   - Adicionar campo `versao INTEGER` na tabela `propostas`
   - Incrementar a cada atualização
   - Apenas a versão atual fica visível, mas o número permite rastreabilidade

### Arquivos afetados
- `apps/api/db/migracoes/005_criar_propostas.sql` (nova tabela ou coluna)
- `apps/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia/repositorio_de_proposta_sqlc.go`
- `apps/api/internal/inscricao_e_selecao_de_projetos/infraestrutura/persistencia/repositorio_de_proposta_memoria.go`
- `apps/api/db/consultas/inscricao_e_selecao_de_projetos/query_propostas.sql`

---

## Issue #2 — Validação de Edital Ativo no Submeter

**Severidade:** MÉDIA
**Tipo:** Validação de negócio ausente

### Problema
O Acceptance Criteria #4 especifica: "O sistema deve validar: edital ativo, prazo dentro da vigência." O use case `SubmeterProposta` não verifica:
- Se o edital está com status "ativo"
- Se a data atual está dentro do período `DataInicio`–`DataFim` do edital

### Solução sugerida
Injetar o repositório de editais (`gestao_de_editais/dominio/repositorios.RepositorioDeEdital`) no `SubmeterProposta` e adicionar validações antes de submeter:

```go
edital, err := uc.editalRepo.BuscarPorID(ctx, proposta.EditalID)
if edital == nil || edital.Status != "ativo" {
    return nil, fmt.Errorf("edital nao esta ativo")
}
if time.Now().Before(edital.DataInicio) || time.Now().After(edital.DataFim) {
    return nil, fmt.Errorf("edital fora do prazo de vigencia")
}
```

### Arquivos afetados
- `apps/api/internal/inscricao_e_selecao_de_projetos/aplicacao/casos_de_uso/submeter_proposta.go`
- `apps/api/internal/inscricao_e_selecao_de_projetos/aplicacao/casos_de_uso/submeter_proposta_test.go`
- `apps/api/cmd/api/main.go` (injeção do repositório de editais)

---

## Status

**APLICADO** — 2026-06-06

Ambas as correções foram implementadas e validadas:

- **Issue #1**: Tabela `versoes_proposta` criada na migration 005 + campo `versao` na entidade `Proposta` + versionamento automático nos repositórios (memória e sqlc) + queries sqlc regeneradas
- **Issue #2**: Interface `EditalVerificador` adicionada ao domínio de propostas + adapter em `main.go` + validação de status "ativo" e prazo de vigência no `SubmeterProposta` + 4 novos testes

```bash
go test ./internal/inscricao_e_selecao_de_projetos/...  # PASS
go vet ./...  # PASS
```

---

*Correções aplicadas por Dex (Dev) · Quinn (QA) aprovou*
