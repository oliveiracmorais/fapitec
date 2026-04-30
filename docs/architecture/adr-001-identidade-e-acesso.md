# ADR-001 - Estratégia de Identidade e Acesso

## Status

Proposto

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

## Decisão em avaliação

Avaliar a adoção de uma abordagem composta:

- `Casdoor` como provedor de identidade;
- `Casbin` como mecanismo de autorização;
- integração do backend em `Go` com essas capacidades;
- frontend em `Next.js` consumindo autenticação e autorização de forma desacoplada.

## Critérios de avaliação

### 1. Aderência funcional

- autenticação por usuário e senha;
- recuperação de senha;
- suporte a MFA, se necessário no futuro;
- gestão de perfis e papéis;
- suporte a regras além de RBAC simples.

### 2. Aderência arquitetural

- compatibilidade com DDD e Clean Architecture;
- possibilidade de isolar fornecedores em adapters;
- baixo acoplamento do domínio a bibliotecas externas.

### 3. Aderência operacional

- facilidade de implantação;
- suporte a ambientes governamentais ou corporativos compatíveis com a `FAPITEC-SE`;
- observabilidade e trilha de auditoria;
- maturidade da integração com Go e Next.js.

### 4. Segurança e conformidade

- compatibilidade com boas práticas de segurança;
- aderência a requisitos de proteção de dados;
- capacidade de governança de acessos e permissões.

## Alternativas iniciais

### Alternativa A

Autenticação e autorização implementadas integralmente no backend próprio.

Vantagens:

- controle máximo;
- menor dependência de terceiros.

Desvantagens:

- maior custo de implementação e manutenção;
- maior risco de erros em uma área sensível.

### Alternativa B

`Casdoor` para identidade e autenticação, com autorização simples no backend.

Vantagens:

- acelera a parte de identidade;
- reduz esforço em fluxos comuns de autenticação.

Desvantagens:

- autorização mais limitada ou dispersa;
- risco de regras de acesso ficarem espalhadas.

### Alternativa C

`Casdoor` para identidade e `Casbin` para autorização.

Vantagens:

- separação clara entre autenticação e autorização;
- boa flexibilidade para RBAC e ABAC;
- alinhamento com sistema de múltiplos perfis e regras institucionais.

Desvantagens:

- maior custo de integração inicial;
- exige desenho cuidadoso de papéis, políticas e auditoria.

## Direção preliminar

A alternativa C parece promissora e merece avaliação mais detalhada antes da implementação.

Ainda não há decisão final.

## Consequências esperadas se confirmada

- o contexto transversal `identidade_e_acesso` poderá ser implementado com adapters para `Casdoor` e `Casbin`;
- o domínio deve continuar independente dessas bibliotecas;
- as regras de autorização devem ficar centralizadas e auditáveis;
- a integração deve ser coberta por testes automatizados.

## Próximos passos

1. Levantar requisitos de perfis e permissões do sistema.
2. Verificar aderência de `Casdoor` aos fluxos exigidos.
3. Verificar aderência de `Casbin` ao modelo de autorização desejado.
4. Decidir se a autenticação e autorização serão plenamente externalizadas ou parcialmente internalizadas.
