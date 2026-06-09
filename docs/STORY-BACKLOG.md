# Story Backlog

## Estatisticas

| Metrica | Valor |
|---------|-------|
| Total | 5 |
| 🔴 HIGH | 1 |
| 🟡 MEDIUM | 2 |
| 🟢 LOW | 2 |
| Source | QA Review 2.2.3 |

---

## 🔴 HIGH

#### [2.2.3-F1] Pre-preenchimento dos dados da sessao no formulario de inscricao
- **Source**: QA Review 2.2.3
- **Priority**: 🔴 HIGH
- **Effort**: 1 hora
- **Status**: 📋 TODO
- **Assignee**: @dev
- **Sprint**: Corrente
- **Description**: Usar `useAuth()` para preencher automaticamente nome (→ nome), documento (→ cpf) e email no formulario de inscricao (`inscrever/page.tsx`). Atualmente `PROPOSTA_PADRAO` inicializa todos os campos vazios, mesmo quando o usuario ja esta logado com esses dados disponiveis na sessao.
- **Success Criteria**:
  - [ ] Ao acessar `/editais/[id]/inscrever`, os campos nome, CPF e email aparecem pre-preenchidos com dados do usuario logado
  - [ ] O usuario ainda pode editar os campos se necessario
  - [ ] Lint + typecheck passam
- **Acceptance**: Implementacao validada por @qa

---

## 🟡 MEDIUM

#### [2.2.3-F2] Upload de documentos sem endpoint backend
- **Source**: QA Review 2.2.3
- **Priority**: 🟡 MEDIUM
- **Effort**: 4 horas
- **Status**: 📋 TODO
- **Assignee**: @dev
- **Sprint**: Proxima
- **Description**: A funcao `uploadDocumento()` em `api-propostas.ts` chama `POST /api/v1/propostas/{id}/documentos` mas o backend nao expoe este endpoint. Upload funciona apenas client-side sem persistencia.
- **Success Criteria**:
  - [ ] Backend expoe endpoint `POST /api/v1/propostas/{id}/documentos` com multipart/form-data
  - [ ] Documentos enviados sao persistedos e retornados na visualizacao da proposta
  - [ ] Validacao de formato e tamanho no backend
- **Acceptance**: Fluxo completo de upload funciona (client → API → storage → retrieval)

#### [2.2.3-F5] Testes unitarios para componentes de proposta
- **Source**: QA Review 2.2.3
- **Priority**: 🟡 MEDIUM
- **Effort**: 2 horas
- **Status**: 📋 TODO
- **Assignee**: @dev
- **Sprint**: Proxima
- **Description**: Nao ha testes automatizados cobrindo os 14 novos componentes criados na Story 2.2.3. Adicionar testes unitarios para validacao, renderizacao de etapas e fluxo de navegacao.
- **Success Criteria**:
  - [ ] Testes para `lib/validacao.ts` (8 novas validacoes)
  - [ ] Testes de renderizacao para `proposta-form-multi-etapas`
  - [ ] Testes de navegacao entre etapas
  - [ ] Testes de calculo orcamentario
- **Acceptance**: `npm test` inclui e passa os novos testes

---

## 🟢 LOW

#### [2.2.3-F3] Etapa empresa vinculada ausente no formulario
- **Source**: QA Review 2.2.3 / Ressalva PO
- **Priority**: 🟢 LOW
- **Effort**: 2 horas
- **Status**: 📋 TODO
- **Assignee**: @dev
- **Sprint**: Futura
- **Description**: O tipo `EmpresaVinculada` existe no `CriarPropostaPayload` mas nao ha etapa no formulario multi-etapas para preencher dados de empresa vinculada. PO aprovou a story com esta ressalva.
- **Success Criteria**:
  - [ ] Etapa "Empresa Vinculada" adicionada ao formulario (entre Dados Pessoais e Formacao, ou conforme definicao do PO)
  - [ ] Campos: nome, CNPJ, porte, enquadramento
- **Acceptance**: Validado por @po

#### [2.2.3-F4] Fluxo de edicao nao submete proposta
- **Source**: QA Review 2.2.3
- **Priority**: 🟢 LOW
- **Effort**: 30 min
- **Status**: 📋 TODO
- **Assignee**: @dev
- **Sprint**: Futura
- **Description**: Ao editar uma proposta (modo rascunho), o fluxo `submeterEDirecionar` so dispara em criacao nova. Na edicao, salva e redireciona sem submeter. Comportamento documentado como intencional, mas pode confundir usuarios que esperam submeter apos editar.
- **Success Criteria**:
  - [ ] Decisao documentada: manter comportamento atual (edicao salva sem submeter) OU adicionar confirmacao "salvar como rascunho vs salvar e submeter"
- **Acceptance**: Comportamento claro para o usuario na UI

---

*Backlog gerado em 2026-06-09 apos QA Review da Story 2.2.3*
