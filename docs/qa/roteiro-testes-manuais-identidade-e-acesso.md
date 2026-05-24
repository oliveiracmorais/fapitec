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

## Cenário 5 — Validação de Conformidade (Story 1.10)

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
