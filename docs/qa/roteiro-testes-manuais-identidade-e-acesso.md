# Roteiro de Testes Manuais — Identidade e Acesso

## Modos de Operação

O sistema possui dois modos de autenticação controlados pela env var `AUTH_PROVIDER`:

| Modo | AUTH_PROVIDER | Quando usar |
|------|--------------|-------------|
| **Internal** | `internal` (default) | Desenvolvimento local sem Casdoor |
| **Casdoor** | `casdoor` | Com Casdoor rodando via Docker |

**Pré-requisitos:**

```bash
# Terminal 1 — Iniciar dependências (PostgreSQL + Casdoor)
cd apps/api && docker compose up -d

# Terminal 2 — Backend (modo Casdoor)
cd apps/api && AUTH_PROVIDER=casdoor go run ./cmd/api

# Terminal 3 — Frontend
cd apps/web && pnpm run dev
```

Para testar no modo `internal` (fallback), omita `AUTH_PROVIDER` ou use `AUTH_PROVIDER=internal`.

---

## Cenário 1 — Cadastro de Usuário (modo internal)

> **Nota:** Desabilitado no modo Casdoor (retorna 410 Gone).

**Via CLI:**
```bash
pnpm run shell:cadastrar "João Silva" 52998224725 joao@example.com "Senha@123"
# Esperado: "Usuario cadastrado: João Silva (ID: 1)"
```

**Via curl:**
```bash
curl -s -X POST http://localhost:8080/api/v1/cadastro \
  -H "Content-Type: application/json" \
  -d '{"nome":"Maria Souza","cpf":"123.456.789-09","email":"maria@example.com","confirmacao_email":"maria@example.com","senha":"Teste@123","confirmacao_senha":"Teste@123"}'
```

### Subcenários

- [ ] CPF já cadastrado → `"CPF já cadastrado no sistema."`
- [ ] E-mail já cadastrado → `"E-mail já cadastrado. Utilize outro endereço ou recupere sua senha."`
- [ ] CPF inválido → `"CPF inválido. Verifique os dígitos."`
- [ ] Senha fraca → mensagem de validação
- [ ] E-mail e confirmação diferentes → `"O e-mail deve ser IGUAL ao e-mail principal."`
- [ ] Senha e confirmação diferentes → `"A senha deve ser IGUAL à primeira senha fornecida."`
- [ ] Modo Casdoor: `POST /api/v1/cadastro` → HTTP 410, `{"erro":"autenticacao propria desabilitada — use /api/v1/auth/login com Casdoor"}`

---

## Cenário 2 — Autenticação (Login) — modo internal

> **Nota:** Desabilitado no modo Casdoor (retorna 410 Gone).

```bash
pnpm run shell:login "52998224725" "Senha@123"
# Esperado: "Autenticado: João Silva (joao@example.com)"
```

### Subcenários

- [ ] Login com CPF/Passaporte ou senha inválidos → `"CPF/Passaporte ou Senha inválidos. Tente novamente."`
- [ ] 5 tentativas falhas consecutivas → bloqueio de 15 minutos
- [ ] Login com conta bloqueada → `"conta temporariamente bloqueada"`
- [ ] Modo Casdoor: `POST /api/v1/login` → HTTP 410, `{"erro":"autenticacao propria desabilitada..."}`

---

## Cenário 3 — Recuperação de Senha (modo internal)

> **Nota:** Desabilitado no modo Casdoor (retorna 410 Gone).

**Passo 1 — Solicitar redefinição:**
```bash
pnpm run shell:solicitar-redefinicao "joao@example.com"
# Verifique o log do backend — deve mostrar o token gerado
```

**Passo 2 — Redefinir senha (copie o token do log):**
```bash
pnpm run shell:redefinir-senha "<token>" "NovaSenha@456"
```

**Passo 3 — Login com nova senha:**
```bash
pnpm run shell:login "52998224725" "NovaSenha@456"
```

### Subcenários

- [ ] Solicitar redefinição para e-mail inexistente → mensagem genérica
- [ ] Token expirado → `"token invalido ou expirado"`
- [ ] Modo Casdoor: `POST /api/v1/solicitar-redefinicao-senha` → HTTP 410

---

## Cenário 4 — Auditoria

```bash
curl -s http://localhost:8080/api/v1/auditoria | python3 -m json.tool
```

- [ ] Eventos de cadastro, login, redefinição aparecem (modo internal)

---

## Cenário 5 — Login OIDC Casdoor (Frontend)

### 5.1 Fluxo completo

1. Abrir `http://localhost:3000`
2. Clicar em **"Entrar com FAPITEC"**
3. Deve redirecionar para `http://localhost:8000/login/oauth/authorize...` (Casdoor)
4. Fazer login com usuário existente no Casdoor:
   - **Admin:** `admin` / `123`
   - **Proponente:** (criar via UI do Casdoor em `http://localhost:8000`)
5. Após autenticar, o Casdoor redireciona de volta para a aplicação
6. **Esperado:** Dashboard é exibido com dados do usuário logado

### 5.2 Callback OIDC (via API)

```bash
# 1. Obter URL de autorização
curl -v "http://localhost:8080/api/v1/auth/login" 2>&1 | grep Location

# 2. Copiar a URL e abrir no navegador, fazer login
# 3. Após o redirect, o código estará na URL de callback
# 4. Trocar o código por token:
curl -s -X POST http://localhost:8080/api/v1/auth/callback \
  -H "Content-Type: application/json" \
  -d '{"code":"<codigo_do_casdoor>","state":"fapitec-state"}'
# Esperado: {"access_token":"<jwt>","token_type":"Bearer"}
```

### Subcenários

- [ ] **Code inválido** → `{"erro":"falha ao obter token: ..."}` (HTTP 401)
- [ ] **State inválido** → Casdoor rejeita (CSRF protection)
- [ ] **Modo internal:** `GET /api/v1/auth/login` → `{"erro":"provedor casdoor nao configurado"}` (HTTP 400)

---

## Cenário 6 — Validação de JWT

### 6.1 JWT válido

```bash
# Obter token via fluxo OIDC (Cenário 5)
TOKEN="<jwt_do_callback>"

# Acessar rota protegida
curl -s http://localhost:8080/api/v1/dashboard/indicadores \
  -H "Authorization: Bearer $TOKEN"
# Esperado: HTTP 200, JSON com indicadores
```

### 6.2 Token ausente

```bash
curl -s http://localhost:8080/api/v1/dashboard/indicadores
# Esperado: HTTP 401, {"erro":"token ausente"}
```

### 6.3 Formato inválido

```bash
curl -s http://localhost:8080/api/v1/dashboard/indicadores \
  -H "Authorization: FormatoErrado"
# Esperado: HTTP 401, {"erro":"formato invalido"}
```

### 6.4 Token expirado

```bash
# Usar JWT com expiracao vencida
curl -s http://localhost:8080/api/v1/dashboard/indicadores \
  -H "Authorization: Bearer eyJhbGciOiJFUzI1NiIs...token_expirado..."
# Esperado: HTTP 401, {"erro":"token invalido"}
```

### 6.5 Rota pública (health) não exige token

```bash
curl -s http://localhost:8080/api/v1/health
# Esperado: HTTP 200, {"status":"ok"} — sem token
```

---

## Cenário 7 — Autorização por Perfil

### 7.1 Acesso permitido

Logar como `proponente`, obter token e acessar:

```bash
# Dashboard (permitido para proponente)
curl -s http://localhost:8080/api/v1/dashboard/indicadores \
  -H "Authorization: Bearer $TOKEN_PROPONENTE"
# Esperado: HTTP 200

# Editais (visualizar permitido para proponente)
curl -s http://localhost:8080/api/v1/editais \
  -H "Authorization: Bearer $TOKEN_PROPONENTE"
# Esperado: HTTP 200
```

### 7.2 Acesso negado

```bash
# Criar edital (POST) — proponente não tem permissao "gerenciar"
curl -s -X POST http://localhost:8080/api/v1/editais \
  -H "Authorization: Bearer $TOKEN_PROPONENTE" \
  -H "Content-Type: application/json" \
  -d '{"titulo":"Teste"}'
# Esperado: HTTP 403, {"erro":"acesso negado"}
```

### 7.3 Admin tem acesso total

Logar como `admin` no Casdoor e testar:

```bash
curl -s -X POST http://localhost:8080/api/v1/editais \
  -H "Authorization: Bearer $TOKEN_ADMIN" \
  -H "Content-Type: application/json" \
  -d '{"titulo":"Edital Admin","descricao":"teste"}'
# Esperado: HTTP 201
```

### Subcenários

- [ ] Proponente tenta POST /api/v1/editais → 403
- [ ] Proponente tenta DELETE /api/v1/editais/1 → 403
- [ ] Admin tenta POST /api/v1/editais → 201
- [ ] Admin tenta DELETE /api/v1/editais/1 → 204
- [ ] Rota sem prefixo mapeado → 403 ("rota sem permissao mapeada")

---

## Cenário 8 — Feature Flag: Rollback AUTH_PROVIDER

### 8.1 Alternar de Casdoor para Internal

```bash
# 1. Com AUTH_PROVIDER=casdoor, testar login OIDC (Cenário 5) — funciona
# 2. Parar o backend, reiniciar com AUTH_PROVIDER=internal
# 3. Testar login interno (Cenário 2) — funciona
# 4. Testar login OIDC (GET /api/v1/auth/login) — HTTP 400, Casdoor nao configurado
```

### 8.2 Alternar de Internal para Casdoor

```bash
# 1. Com AUTH_PROVIDER=internal, testar login interno — funciona
# 2. Parar o backend, reiniciar com AUTH_PROVIDER=casdoor
# 3. Testar login OIDC — funciona
# 4. Testar login interno (POST /api/v1/login) — HTTP 410, desabilitado
```

- [ ] Rollback sem downtime (dados preservados)
- [ ] Mensagens de erro claras em cada modo

---

## Cenário 9 — Seed de Perfis Casdoor

### 9.1 Verificar roles criadas

1. Acessar UI do Casdoor: `http://localhost:8000`
2. Login: `admin` / `123`
3. Navegar em **Roles**
4. Verificar se as 6 roles existem:
   - [ ] `administrador_fapitec`
   - [ ] `instituicao_ensino`
   - [ ] `funcionario_fapitec`
   - [ ] `diretoria`
   - [ ] `proponente`
   - [ ] `avaliador`

### 9.2 Executar seed

```bash
cd apps/api && go run ./cmd/seed
# Esperado: "Role ja existe" para cada role (idempotente)
```

### 9.3 Verificar admin

- [ ] Usuário `admin_fapitec` criado no Casdoor
- [ ] Login `admin_fapitec` / senha (definida via `CASDOOR_ADMIN_PASSWORD` ou gerada)

---

## Cenário 10 — User Profile com Casdoor

```bash
# Obter token JWT via OIDC
TOKEN="<jwt>"

# Consultar profile (usa dados do JWT quando Casdoor ativo)
curl -s http://localhost:8080/api/v1/user-profile \
  -H "Authorization: Bearer $TOKEN"
# Esperado: HTTP 200, JSON com documento, nome, email, perfil do JWT
```

- [ ] Modo Casdoor: retorna dados do JWT (documento, nome, email, perfil)
- [ ] Modo internal: retorna dados do banco (id, nome, documento, email, criado_em)

---

## Cenário 11 — Endpoint Aliases (modo internal)

### 11.1 POST /api/v1/register

```bash
curl -s -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"nome":"Joao Aliasing","cpf":"111.222.333-44","email":"joao.alias@example.com","confirmacao_email":"joao.alias@example.com","senha":"Teste@123","confirmacao_senha":"Teste@123"}'
# Esperado: HTTP 201
```

### 11.2 POST /api/v1/reset-password

```bash
curl -s -X POST http://localhost:8080/api/v1/reset-password \
  -H "Content-Type: application/json" \
  -d '{"email":"joao.alias@example.com"}'
# Esperado: HTTP 200, {"mensagem":"Se o e-mail estiver cadastrado..."}
```

### Subcenários

- [ ] Alias em modo Casdoor → HTTP 410 (desabilitado)
- [ ] Modo internal: mesmo comportamento dos endpoints originais

---

## Cenário 12 — Captcha (Turnstile) — modo internal

**Pré-requisito:** `AUTH_PROVIDER=internal`, `TURNSTILE_SECRET_KEY` configurada.

### 12.1 Captcha não aparece em tentativas normais

```bash
curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"cpf":"000.000.000-00","senha":"errada"}'
# Esperado: HTTP 401, sem captcha
```

### 12.2 Captcha após 3 tentativas falhas

```bash
curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"cpf":"000.000.000-00","senha":"errada","captcha_token":"token_invalido"}'
# Esperado: HTTP 401, {"erro":"Validacao de captcha falhou..."}
```

---

## Cenário 13 — Validação de Conformidade (Story 1.10)

**Passaporte (modo internal, via curl):**
```bash
curl -s -X POST http://localhost:8080/api/v1/cadastro \
  -H "Content-Type: application/json" \
  -d '{"nome":"John Doe","cpf":"AB123456","email":"john@example.com","confirmacao_email":"john@example.com","senha":"Teste@123","confirmacao_senha":"Teste@123","estrangeiro":true}'
# Esperado: HTTP 201
```

**Check-email e check-cpf (ambos os modos):**
```bash
curl -s "http://localhost:8080/api/v1/check-email?email=joao@example.com"
# {"existe":true}

curl -s "http://localhost:8080/api/v1/check-cpf?cpf=52998224725"
# {"existe":true}
```

---

## Pontos de Atenção

- **Modo Casdoor**: o Casdoor deve estar rodando (`docker compose up -d` em `apps/api/`)
- **Modo internal**: dados ficam em memória — são perdidos ao reiniciar
- **JWT expirado**: configurar expiração para testar (ou usar token manipulado)
- **Rate limiting**: Casdoor tem proteção própria contra força bruta (configurar via UI)
- **Senha admin seed**: definida via `CASDOOR_ADMIN_PASSWORD` ou gerada automaticamente
- Ao alternar `AUTH_PROVIDER`, reinicie o backend
- `POST /api/v1/cadastro` e `POST /api/v1/login` retornam **410** no modo Casdoor
