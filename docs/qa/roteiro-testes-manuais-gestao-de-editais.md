# Roteiro de Testes Manuais — Gestão de Editais

## Pré-requisitos

```bash
# Terminal 1 — Backend
cd apps/api && go run ./cmd/api

# Terminal 2 — Frontend (opcional para testes visuais)
cd apps/web && pnpm run dev
```

> ⚠️ O backend inicia com PostgreSQL. Se o banco não estiver disponível, usa repositório em memória automaticamente. Dados são perdidos ao reiniciar.

---

## Cenário 1 — Criar Edital (POST /api/v1/editais)

**Via curl:**
```bash
curl -s -X POST http://localhost:8080/api/v1/editais \
  -H "Content-Type: application/json" \
  -d '{"nome":"Edital APQ 2026","descricao":"Edital de apoio a pesquisa","data_inicio":"2026-06-01","data_fim":"2026-12-31","tipo_chamada":"APQ"}'
```

**Esperado:** HTTP 201, JSON com `id`, `nome`, `descricao`, `data_inicio`, `data_fim`, `status: "ativo"`, `tipo_chamada: "APQ"`, `criado_em`.

### Subcenários

- [ ] **Nome vazio** → `{"erro":"nome do edital é obrigatorio"}` (HTTP 400)
  ```bash
  curl -s -X POST http://localhost:8080/api/v1/editais \
    -H "Content-Type: application/json" \
    -d '{"nome":"","descricao":"Teste","data_inicio":"2026-06-01","data_fim":"2026-12-31","tipo_chamada":"APQ"}'
  ```

- [ ] **Data inválida** → `{"erro":"data de inicio invalida: ..."}` (HTTP 400)
  ```bash
  curl -s -X POST http://localhost:8080/api/v1/editais \
    -H "Content-Type: application/json" \
    -d '{"nome":"Edital","descricao":"Teste","data_inicio":"invalida","data_fim":"2026-12-31","tipo_chamada":"APQ"}'
  ```

- [ ] **Data início após data fim** → `{"erro":"data de inicio nao pode ser posterior a data de fim"}` (HTTP 400)
  ```bash
  curl -s -X POST http://localhost:8080/api/v1/editais \
    -H "Content-Type: application/json" \
    -d '{"nome":"Edital","descricao":"Teste","data_inicio":"2027-01-01","data_fim":"2026-12-31","tipo_chamada":"APQ"}'
  ```

- [ ] **JSON malformado** → `{"erro":"requisicao invalida"}` (HTTP 400)

---

## Cenário 2 — Listar Editais (GET /api/v1/editais)

**Via curl:**
```bash
curl -s http://localhost:8080/api/v1/editais
```

**Esperado:** HTTP 200, JSON com `editais` (array) e `total`.

### Subcenários

- [ ] **Listagem vazia** → `{"editais":[],"total":0}` (antes de criar qualquer edital)
- [ ] **Listagem com 1+ editais** → array populado após criar ao menos um
- [ ] **Filtro por título:**
  ```bash
  curl -s "http://localhost:8080/api/v1/editais?titulo=APQ"
  ```
- [ ] **Filtro por status:**
  ```bash
  curl -s "http://localhost:8080/api/v1/editais?status=ativo"
  ```
- [ ] **Filtro por tipo_chamada:**
  ```bash
  curl -s "http://localhost:8080/api/v1/editais?tipo_chamada=ARC"
  ```
- [ ] **Filtros combinados:**
  ```bash
  curl -s "http://localhost:8080/api/v1/editais?titulo=APQ&status=ativo"
  ```
- [ ] **Filtro sem resultados** → `{"editais":[],"total":0}`

---

## Cenário 3 — Visualizar Edital (GET /api/v1/editais/{id})

Crie um edital primeiro e guarde o `id` retornado:

```bash
EDITAL_ID=$(curl -s -X POST http://localhost:8080/api/v1/editais \
  -H "Content-Type: application/json" \
  -d '{"nome":"Edital ARC 2026","descricao":"Edital de inovacao","data_inicio":"2026-07-01","data_fim":"2026-12-31","tipo_chamada":"ARC"}' | python3 -c "import sys,json;print(json.load(sys.stdin)['id'])")

curl -s "http://localhost:8080/api/v1/editais/$EDITAL_ID"
```

**Esperado:** HTTP 200, JSON com todos os campos do edital.

### Subcenários

- [ ] **ID inexistente** → `{"erro":"edital nao encontrado"}` (HTTP 404)
  ```bash
  curl -s http://localhost:8080/api/v1/editais/99999
  ```
- [ ] **ID inválido (não numérico)** → `{"erro":"id invalido"}` (HTTP 400)
  ```bash
  curl -s http://localhost:8080/api/v1/editais/abc
  ```

---

## Cenário 4 — Atualizar Edital (PUT /api/v1/editais/{id})

**Via curl (atualização parcial — apenas o que for enviado):**
```bash
EDITAL_ID=1  # substitua pelo ID real

curl -s -X PUT "http://localhost:8080/api/v1/editais/$EDITAL_ID" \
  -H "Content-Type: application/json" \
  -d '{"nome":"Edital APQ 2026 - Atualizado","status":"em_avaliacao"}'
```

**Esperado:** HTTP 200, JSON com dados atualizados. Status deve refletir a alteração.

### Subcenários

- [ ] **Atualizar nome** → nome alterado no retorno
- [ ] **Atualizar descrição** → descrição alterada
- [ ] **Atualizar datas** → datas alteradas (formato `"2006-01-02"`)
- [ ] **Atualizar status** → `"encerrado"` ou `"em_avaliacao"`
- [ ] **Atualizar tipo_chamada** → tipo alterado
- [ ] **Data início após data fim** → `{"erro":"data de inicio nao pode ser posterior a data de fim"}` (HTTP 400)
- [ ] **Status inválido** → `{"erro":"status de edital invalido: ..."}` (HTTP 400)
  ```bash
  curl -s -X PUT "http://localhost:8080/api/v1/editais/$EDITAL_ID" \
    -H "Content-Type: application/json" \
    -d '{"status":"status_inexistente"}'
  ```
- [ ] **ID inexistente** → `{"erro":"edital nao encontrado"}` (HTTP 400)
  ```bash
  curl -s -X PUT http://localhost:8080/api/v1/editais/99999 \
    -H "Content-Type: application/json" \
    -d '{"nome":"Teste"}'
  ```

---

## Cenário 5 — Deletar Edital (DELETE /api/v1/editais/{id})

Crie um edital temporário e delete-o:

```bash
EDITAL_ID=$(curl -s -X POST http://localhost:8080/api/v1/editais \
  -H "Content-Type: application/json" \
  -d '{"nome":"Edital Temporario","descricao":"Para deletar","data_inicio":"2026-06-01","data_fim":"2026-12-31","tipo_chamada":"APQ"}' | python3 -c "import sys,json;print(json.load(sys.stdin)['id'])")

curl -s -o /dev/null -w "%{http_code}" -X DELETE "http://localhost:8080/api/v1/editais/$EDITAL_ID"
```

**Esperado:** HTTP 204 (No Content).

### Subcenários

- [ ] **Conferir exclusão** — `GET /api/v1/editais/{id}` após deletar → `{"erro":"edital nao encontrado"}` (HTTP 404)
- [ ] **ID inexistente** → `{"erro":"edital nao encontrado"}` (HTTP 404)
  ```bash
  curl -s -X DELETE http://localhost:8080/api/v1/editais/99999
  ```
- [ ] **ID inválido** → `{"erro":"id invalido"}` (HTTP 400)
  ```bash
  curl -s -X DELETE http://localhost:8080/api/v1/editais/abc
  ```

---

## Cenário 6 — Fluxo Completo (Fumográfico)

```bash
# 1. Criar 3 editais
curl -s -X POST http://localhost:8080/api/v1/editais \
  -H "Content-Type: application/json" \
  -d '{"nome":"APQ 2026.1","descricao":"Pesquisa basica","data_inicio":"2026-01-01","data_fim":"2026-06-30","tipo_chamada":"APQ"}'

curl -s -X POST http://localhost:8080/api/v1/editais \
  -H "Content-Type: application/json" \
  -d '{"nome":"ARC 2026.1","descricao":"Inovacao tecnologica","data_inicio":"2026-03-01","data_fim":"2026-09-30","tipo_chamada":"ARC"}'

curl -s -X POST http://localhost:8080/api/v1/editais \
  -H "Content-Type: application/json" \
  -d '{"nome":"APQ 2026.2","descricao":"Pesquisa aplicada","data_inicio":"2026-07-01","data_fim":"2026-12-31","tipo_chamada":"APQ"}'

# 2. Listar todos (esperado: 3)
curl -s http://localhost:8080/api/v1/editais | python3 -c "import sys,json;d=json.load(sys.stdin);print(f'Total: {d[\"total\"]}')"

# 3. Filtrar por APQ (esperado: 2)
curl -s "http://localhost:8080/api/v1/editais?tipo_chamada=APQ" | python3 -c "import sys,json;d=json.load(sys.stdin);print(f'APQ: {d[\"total\"]}')"

# 4. Atualizar um para "encerrado"
EDITAL_ID=1
curl -s -X PUT "http://localhost:8080/api/v1/editais/$EDITAL_ID" \
  -H "Content-Type: application/json" \
  -d '{"status":"encerrado"}'

# 5. Filtrar por "encerrado" (esperado: 1)
curl -s "http://localhost:8080/api/v1/editais?status=encerrado" | python3 -c "import sys,json;d=json.load(sys.stdin);print(f'Encerrados: {d[\"total\"]}')"

# 6. Deletar o encerrado
curl -s -X DELETE "http://localhost:8080/api/v1/editais/$EDITAL_ID"

# 7. Listar novamente (esperado: 2)
curl -s http://localhost:8080/api/v1/editais | python3 -c "import sys,json;d=json.load(sys.stdin);print(f'Total: {d[\"total\"]}')"
```

---

## Pontos de Atenção

- Repositório **em memória**: dados são perdidos ao reiniciar o backend
- Repositório **PostgreSQL**: persistente; para resetar, execute `TRUNCATE editais;` ou recrie o banco
- Status válidos: `ativo`, `encerrado`, `em_avaliacao`
- Tipos de chamada: `APQ`, `ARC` (string livre, sem validação no backend)
- Para resetar estado entre testes, reinicie o processo Go (`Ctrl+C` e `go run ./cmd/api`)
