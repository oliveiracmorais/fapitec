# Backlog Inicial de Stories

Este arquivo organiza as primeiras stories fundacionais do projeto da `FAPITEC-SE` a partir do Termo de Referência e das decisões arquiteturais iniciais.

## Sequência sugerida

1. Story: estrutura base do monorepo
2. Story: convenções arquiteturais e linguagem ubíqua
3. Story: autenticação e cadastro inicial (especificação)
4. Story: autorização por perfis (especificação)
5. Story: trilha de auditoria (especificação)
6. Story: shell inicial e dashboard modular
7. Story: integração runtime mínima de autenticação, autorização e auditoria
8. Story: implementação de autenticação e cadastro (código)
9. Story: gestão documental base
10. Story: módulo de gestão de editais ✅

## Sequenciamento em ondas

### Onda 1 - Fundação mínima

- Story `1.1`: estrutura base do monorepo
- Story `1.2`: convenções arquiteturais e linguagem ubíqua
- Story `1.3`: autenticação e cadastro inicial (especificação) ✅
- Story `1.4`: autorização por perfis (especificação)
- Story `1.5`: trilha de auditoria mínima (especificação)
- Story `1.6`: shell inicial e dashboard modular
- Story `1.7`: integração runtime mínima de autenticação, autorização e auditoria
- Story `1.9`: implementação de autenticação e cadastro (código) ✅
- Story `1.10`: ajustes de conformidade (passaporte, mensagens, confirmação de e-mail) ✅
- Story `1.11`: recuperação de senha ✅
- Story `1.12`: gestão de editais 🔜

Objetivo:

- criar a base mínima necessária para que os módulos de negócio nasçam com segurança, rastreabilidade, coerência arquitetural e macroestrutura validável pela `FAPITEC-SE`.

### Onda 2 - Núcleo operacional

- Story futura: gestão documental base
- Story futura: notificações base
- Story futura: gestão de editais

Objetivo:

- viabilizar o primeiro núcleo operacional relevante da plataforma da `FAPITEC-SE`.

### Onda 3 - Entrada e seleção

- Story futura: inscrição de projetos
- Story futura: seleção e avaliação de projetos

Objetivo:

- permitir entrada, análise e decisão sobre propostas submetidas.

### Onda 4 - Execução

- Story futura: solicitação e concessão de auxílios
- Story futura: gestão de bolsistas e pesquisadores vinculados

Objetivo:

- suportar a operacionalização de auxílios, vínculos e execução associada.

### Onda 5 - Controle e conformidade

- Story futura: prestação de contas
- Story futura: tomada de contas especial
- Story futura: relatórios e painéis de indicadores

Objetivo:

- garantir governança, conformidade, rastreabilidade e apoio à tomada de decisão.

## Dependências principais

- autenticação antecede a maior parte dos fluxos internos e externos.
- autorização antecede o controle de acesso fino por perfil e módulo.
- auditoria mínima deve existir antes dos primeiros fluxos críticos de negócio.
- o shell inicial deve validar macroestrutura, nomes de módulos e acesso antes do detalhamento dos módulos de negócio.
- a integração runtime mínima deve remover o perfil fixo do dashboard e demonstrar sessão, autorização e auditoria antes da apresentação institucional.
- gestão documental base antecede módulos com anexos, comprovações e documentos formais.
- notificações base devem existir antes dos fluxos institucionais que dependem de comunicação operacional.
- gestão de editais funciona como eixo para diversos fluxos posteriores.

## Observações

- Os PDFs em [docs/stories](/home/carlos/projetos/fapitec/docs/stories) continuam sendo a fonte original dos requisitos disponíveis até aqui.
- As stories formais agora passam a ser redigidas em arquivos próprios, com critérios de aceite, checklist e file list.
- A avaliação de `Casdoor` e `Casbin` deve acontecer antes da story de autorização por perfis.
- O backlog deve respeitar o cronograma mensal de entregas previsto no Termo de Referência.
- As capacidades transversais devem ser implementadas de forma incremental, como fundação mínima para as ondas seguintes.
- A validação da macroestrutura pela `FAPITEC-SE` pode gerar ajustes de nomenclatura, perfis, módulos ou prioridade, que devem ser registrados em story, ADR ou backlog conforme o impacto.

## Stories fundacionais

- [1.1.estrutura-base-do-monorepo.md](/home/carlos/projetos/fapitec/docs/stories/1.1.estrutura-base-do-monorepo.md)
- [1.2.convencoes-arquiteturais-e-linguagem-ubiqua.md](/home/carlos/projetos/fapitec/docs/stories/1.2.convencoes-arquiteturais-e-linguagem-ubiqua.md)
- [1.3.autenticacao-e-cadastro-inicial.md](/home/carlos/projetos/fapitec/docs/stories/1.3.autenticacao-e-cadastro-inicial.md)
- [1.4.autorizacao-por-perfis.md](/home/carlos/projetos/fapitec/docs/stories/1.4.autorizacao-por-perfis.md)
- [1.5.trilha-de-auditoria-minima.md](/home/carlos/projetos/fapitec/docs/stories/1.5.trilha-de-auditoria-minima.md)
- [1.6.shell-inicial-e-dashboard-modular.md](/home/carlos/projetos/fapitec/docs/stories/1.6.shell-inicial-e-dashboard-modular.md)
- [1.7.integracao-runtime-minima-de-autenticacao-autorizacao-e-auditoria.md](/home/carlos/projetos/fapitec/docs/stories/1.7.integracao-runtime-minima-de-autenticacao-autorizacao-e-auditoria.md)
- [1.8.tela-de-login-simulado-com-auditoria.md](/home/carlos/projetos/fapitec/docs/stories/1.8.tela-de-login-simulado-com-auditoria.md)
- [1.9.implementacao-autenticacao-e-cadastro.md](/home/carlos/projetos/fapitec/docs/stories/1.9.implementacao-autenticacao-e-cadastro.md)
- [1.10.ajustes-de-conformidade-autenticacao-e-cadastro.md](/home/carlos/projetos/fapitec/docs/stories/1.10.ajustes-de-conformidade-autenticacao-e-cadastro.md)
- [1.11.recuperacao-de-senha.md](/home/carlos/projetos/fapitec/docs/stories/1.11.recuperacao-de-senha.md)
- [1.12.gestao-de-editais.md](/home/carlos/projetos/fapitec/docs/stories/1.12.gestao-de-editais.md)
- [1.13.identidade-e-acesso-ponta-a-ponta.md](/home/carlos/projetos/fapitec/docs/stories/1.13.identidade-e-acesso-ponta-a-ponta.md)
