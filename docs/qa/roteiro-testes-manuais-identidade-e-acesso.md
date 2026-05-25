# Roteiro de Testes Manuais — Identidade e Acesso

## Pré-requisitos

```bash
# Terminal 1 — Backend
cd apps/api && go run ./cmd/api

# Terminal 2 — Frontend
cd apps/web && pnpm run dev
```

## Cenário 1 — Cadastro de Usuário

**Via CLI (recomendado para validar API):**
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

**Via navegador:** `http://localhost:3000/recuperar-senha` (apenas recovery, cadastro não tem tela ainda)

### Subcenários

- [ ] CPF já cadastrado → `"CPF já cadastrado no sistema."`
- [ ] E-mail já cadastrado → `"E-mail já cadastrado. Utilize outro endereço ou recupere sua senha."`
- [ ] CPF inválido → `"CPF inválido. Verifique os dígitos."`
- [ ] Senha fraca → mensagem de validação
- [ ] E-mail e confirmação diferentes → `"O e-mail deve ser IGUAL ao e-mail principal."`
- [ ] Senha e confirmação diferentes → `"A senha deve ser IGUAL à primeira senha fornecida."`

---

## Cenário 2 — Autenticação (Login)

**Via CLI:**
```bash
pnpm run shell:login "52998224725" "Senha@123"
# Esperado: "Autenticado: João Silva (joao@example.com)"
```

**Via curl:**
```bash
curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"cpf":"52998224725","senha":"Senha@123"}'
```

### Subcenários

- [ ] Login com CPF/Passaporte ou senha inválidos → `"CPF/Passaporte ou Senha inválidos. Tente novamente."`
- [ ] 5 tentativas falhas consecutivas → bloqueio de 15 minutos
- [ ] Login com conta bloqueada → `"conta temporariamente bloqueada"`

---

## Cenário 3 — Recuperação de Senha

**Passo 1 — Solicitar redefinição:**
```bash
pnpm run shell:solicitar-redefinicao "joao@example.com"
# Esperado: "Se o e-mail estiver cadastrado, voce recebera um link de redefinicao de senha."

# Verifique o log do backend — deve mostrar o token gerado:
# [EMAIL-PLACEHOLDER] Para: joao@example.com | Token de redefinicao: <token>
```

**Passo 2 — Redefinir senha (copie o token do log):**
```bash
pnpm run shell:redefinir-senha "<token>" "NovaSenha@456"
# Esperado: "Senha redefinida com sucesso."
```

**Passo 3 — Login com nova senha:**
```bash
pnpm run shell:login "52998224725" "NovaSenha@456"
# Esperado: "Autenticado: João Silva (joao@example.com)"
```

### Subcenários

- [ ] Solicitar redefinição para e-mail inexistente → mensagem genérica, sem vazar info
- [ ] Token expirado → `"token invalido ou expirado"` (aguardar 1h ou usar token fake)
- [ ] Token inexistente → `"token invalido ou expirado"`
- [ ] Redefinir com senha fraca → erro de validação
- [ ] Redefinir com confirmação diferente → `"A senha deve ser IGUAL à primeira senha fornecida."`

---

## Cenário 4 — Auditoria

**Via CLI:**
```bash
# A CLI atual não tem comando específico; use curl:
curl -s http://localhost:8080/api/v1/auditoria | python3 -m json.tool
```

**Eventos esperados após executar todos os cenários acima:**
- [ ] `cadastro_de_usuario` — sucesso (João Silva)
- [ ] `login` — sucesso (João Silva)
- [ ] `solicitacao_redefinicao_senha` — sucesso
- [ ] `redefinicao_de_senha` — sucesso

---

## Cenário 5 — Endpoint Aliases (Conformidade com ProjetoFapitec.txt item 3.1)

### 5.1 POST /api/v1/register (alias do /cadastro)

Deve ter o mesmo comportamento de `POST /api/v1/cadastro`:

```bash
curl -s -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"nome":"Joao Aliasing","cpf":"111.222.333-44","email":"joao.alias@example.com","confirmacao_email":"joao.alias@example.com","senha":"Teste@123","confirmacao_senha":"Teste@123"}'
```

**Esperado:** HTTP 201, mesmo response de `/cadastro`.

### Subcenários

- [ ] **POST /api/v1/register com dados inválidos** → HTTP 400, mesma validação de `/cadastro`
- [ ] **POST /api/v1/register com CPF duplicado** → `"CPF já cadastrado no sistema."`

---

### 5.2 POST /api/v1/reset-password (alias do /solicitar-redefinicao-senha)

Deve ter o mesmo comportamento de `POST /api/v1/solicitar-redefinicao-senha`:

```bash
curl -s -X POST http://localhost:8080/api/v1/reset-password \
  -H "Content-Type: application/json" \
  -d '{"email":"joao.alias@example.com"}'
```

**Esperado:** HTTP 200, `{"mensagem":"Se o e-mail estiver cadastrado, voce recebera um link de redefinicao de senha."}`

### Subcenários

- [ ] **E-mail inexistente** → mesma resposta genérica (não vaza informação)
- [ ] **Request sem body** → `{"erro":"requisicao invalida"}` (HTTP 400)

---

### 5.3 GET /api/v1/user-profile

**Por CPF:**
```bash
# Primeiro crie um usuário e guarde o CPF:
curl -s -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"nome":"Perfil Teste","cpf":"999.888.777-66","email":"perfil@example.com","confirmacao_email":"perfil@example.com","senha":"Teste@123","confirmacao_senha":"Teste@123"}'

# Consulte por CPF:
curl -s "http://localhost:8080/api/v1/user-profile?cpf=999.888.777-66"
```

**Esperado:** HTTP 200, JSON com `id`, `nome`, `documento`, `email`, `estrangeiro`, `criado_em`.

**Por e-mail:**
```bash
curl -s "http://localhost:8080/api/v1/user-profile?email=perfil@example.com"
```

**Esperado:** Mesmo response da consulta por CPF.

### Subcenários

- [ ] **Sem parâmetros** → `{"erro":"informe cpf ou email"}` (HTTP 400)
  ```bash
  curl -s "http://localhost:8080/api/v1/user-profile"
  ```
- [ ] **CPF inexistente** → `{"erro":"usuario nao encontrado"}` (HTTP 404)
  ```bash
  curl -s "http://localhost:8080/api/v1/user-profile?cpf=000.000.000-00"
  ```
- [ ] **E-mail inexistente** → `{"erro":"usuario nao encontrado"}` (HTTP 404)
  ```bash
  curl -s "http://localhost:8080/api/v1/user-profile? email=inexistente@test.com"
  ```

---

## Cenário 6 — Captcha (Turnstile) — Tentativas Excessivas

**Pré-requisito:** Backend rodando com `TURNSTILE_SECRET_KEY` configurada (padrão: chave de teste).

### 6.1 Captcha não aparece em tentativas normais

```bash
# Fazer login com credenciais inválidas 2 vezes — captcha NÃO deve aparecer
curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"cpf":"000.000.000-00","senha":"errada"}'
```

**Esperado:** HTTP 401, sem campo `captcha_token` no body, captcha não visível no frontend.

### 6.2 Captcha aparece após 3 tentativas falhas

1. No navegador, tente fazer login 3 vezes com credenciais inválidas
2. Na 4ª tentativa, o widget Turnstile deve aparecer acima do botão "Entrar"

**Via curl (simulando com `captcha_token` inválido):**
```bash
curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"cpf":"000.000.000-00","senha":"errada","captcha_token":"token_invalido"}'
```

**Esperado:** HTTP 401, `{"erro":"Validacao de captcha falhou. Tente novamente."}`

### 6.3 Captcha válido (chave de teste)

Com a chave de teste configurada, qualquer token é aceito:

```bash
curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"cpf":"123.456.789-09","senha":"senha_correta","captcha_token":"token_teste"}'
```

**Esperado:** HTTP 200 (login bem-sucedido) — a chave de teste sempre aprova.

### 6.4 Fluxo completo no frontend

1. Abrir `http://localhost:3000`
2. Tentar login 3 vezes com senha errada (ex: CPF `123.456.789-09`, senha `errada`)
3. Observar contagem de tentativas no backend (`GET /api/v1/auditoria` deve mostrar 3 `falha_de_login`)
4. Na 4ª tentativa, o widget Turnstile aparece
5. Resolver o captcha (com chave de teste, qualquer clique funciona)
6. Inserir senha correta → login autorizado

---

## Cenário 7 — Validação de Conformidade (Story 1.10)

**Passaporte (via curl):**
```bash
curl -s -X POST http://localhost:8080/api/v1/cadastro \
  -H "Content-Type: application/json" \
  -d '{"nome":"John Doe","cpf":"AB123456","email":"john@example.com","confirmacao_email":"john@example.com","senha":"Teste@123","confirmacao_senha":"Teste@123","estrangeiro":true}'
```

**Check-email e check-cpf (duplicidade):**
```bash
curl -s "http://localhost:8080/api/v1/check-email?email=joao@example.com"
# {"existe":true}

curl -s "http://localhost:8080/api/v1/check-cpf?cpf=52998224725"
# {"existe":true}
```

---

## Pontos de Atenção

- O repositório é **em memória** — dados são perdidos ao reiniciar o backend
- O placeholder de e-mail apenas **loga no console** — não envia e-mail real
- Para resetar o estado entre testes, reinicie o processo Go (Ctrl+C e `make run`)
- O proxy Next.js redireciona `/api/v1/*` para `localhost:8080` — ambos precisam estar rodando
