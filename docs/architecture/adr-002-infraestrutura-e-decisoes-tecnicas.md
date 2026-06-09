# ADR-002 — Infraestrutura de Hospedagem e Decisões Técnicas

## Status

Aceito

## Contexto

Durante a fase de homologação e demonstração do sistema, algumas decisões arquiteturais e de infraestrutura precisaram ser ajustadas em relação ao Termo de Referência original, seja por viabilidade técnica, seja por priorização de entregas.

As questões identificadas foram:

1. **Infraestrutura de hospedagem** — O Termo de Referência (Seção 7.1) determina que o sistema seja alocado na infraestrutura da EMGETIS. No entanto, o ambiente de homologação atual está provisionado em Oracle Cloud VM (ARM64) com DuckDNS e Let's Encrypt.
2. **ORM vs sqlc** — O Termo de Referência (Seção 5.3) menciona "Active Record para mapeamento objeto-relacional (ORM)". O projeto adota sqlc como abordagem de acesso a dados.
3. **Assinatura Digital Gov.br** — O Termo de Referência (Seção 5.8.6) exige integração com a plataforma de assinatura eletrônica do gov.br.
4. **Autenticação Multifator (MFA)** — O Termo de Referência (Seção 5.4) exige MFA.

## Decisões

### 1. Infraestrutura de Hospedagem

**Decisão:** Manter o ambiente de homologação em Oracle Cloud VM para fins de demonstração e desenvolvimento. A migração para infraestrutura EMGETIS será realizada em momento posterior, quando houver definição contratual e provisionamento do ambiente pela EMGETIS.

**Motivação:**
- A infraestrutura EMGETIS não está disponível/disponibilizada neste momento
- O ambiente Oracle Cloud VM atende aos requisitos de homologação com HTTPS via Let's Encrypt, Docker Compose e CI/CD
- A migração futura é viável por tratar-se de ambiente containerizado (Docker Compose), independente de provedor de nuvem
- O custo atual é zero (Always Free Tier Oracle Cloud ARM64)

**Consequências:**
- Positivas: ambiente funcional para demonstração, validação com o cliente, e desenvolvimento continuo
- Negativas: necessidade de migração futura para EMGETIS, com possível ajuste de CI/CD e DNS
- Risco: baixo, dado o isolamento por containers

### 2. ORM vs sqlc

**Decisão:** Manter sqlc como abordagem de acesso a dados, conforme já documentado em `docs/architecture.md` (Seção 4.2).

**Motivação:**
- sqlc gera código Go类型-safe a partir de SQL, eliminando erros de runtime comuns em ORMs
- Performance superior por não utilizar reflection ou query building em runtime
- Controle total sobre as queries SQL, permitindo otimizações por contexto de domínio
- Alinhamento com DDD e Clean Architecture ao evitar vazamento de preocupações de persistência para o domínio
- A escolha não prejudica o escopo do projeto nem a capacidade de atender aos requisitos funcionais

**Consequências:**
- Positivas: código类型-safe, performance, manutenibilidade, alinhamento com a arquitetura
- Negativas: não segue literalmente a recomendação de "Active Record" do Termo de Referência, mas sem impacto no escopo contratual

### 3. Assinatura Digital Gov.br

**Decisão:** Adiar a implementação da integração com a plataforma de assinatura eletrônica do gov.br para uma fase futura do projeto.

**Motivação:**
- A funcionalidade depende de conta institucional no gov.br e acesso à API, que é responsabilidade da FAPITEC-SE provisionar
- Não é impeditiva para a demonstração do fluxo básico de submissão e avaliação de propostas
- A priorização atual é concluir o ciclo completo de gestão de editais (Epic 2.2)

**Consequências:**
- Positivas: permite foco na entrega do MVP sem dependência externa
- Negativas: funcionalidade contratualmente exigida será entregue em versão posterior

### 4. Autenticação Multifator (MFA)

**Decisão:** Implementar MFA no ambiente de homologação, mas **após a conclusão da Epic 2.2** (Gestão de Editais e Submissão de Propostas).

**Motivação:**
- O Casdoor já suporta MFA nativamente (TOTP), exigindo apenas configuração e ativação
- MFA é requisito relevante de segurança (Seção 5.4 do TR) e agrega valor para demonstração
- A priorização atual é concluir o fluxo de negócio de editais primeiro

**Consequências:**
- Positivas: segurança adicional no ambiente de homologação quando implementado
- Negativas: ambiente de homologação temporariamente sem MFA até conclusão da Epic 2.2

## Linha do Tempo

| Decisão | Quando |
|---------|--------|
| sqlc (já implementado) | Imediato — documentado |
| Infraestrutura Oracle Cloud | Até definição da EMGETIS |
| Conclusão Epic 2.2 | Próxima entrega |
| MFA (Casdoor TOTP) | Após Epic 2.2 |
| Assinatura Gov.br | Futuro (depende de conta gov.br institucional) |
| Migração EMGETIS | Futuro (depende de provisionamento) |

## Referências

- Termo de Referência Seção 5.3, 5.4, 5.8.6, 7.1
- `docs/architecture.md` — Stack inicial e decisões arquiteturais
- `docs/architecture/adr-001-identidade-e-acesso.md` — Decisão de IAM com Casdoor
- `docs/infra/ambiente-homologacao.md` — Ambiente atual de homologação
