# Documento de Arquitetura Inicial

## 1. Objetivo

Este documento consolida as primeiras decisões arquiteturais do projeto, com base nos artefatos disponíveis em [Anexo I - Termo de Referência.pdf](/home/carlos/projetos/fapitec/docs/stories/Anexo%20I%20-%20Termo%20de%20Refer%C3%AAncia.pdf) e [Pedido Esclarecimento 002.pdf](/home/carlos/projetos/fapitec/docs/stories/Pedido%20Esclarecimento%20002.pdf).

O objetivo é orientar a evolução do sistema da `FAPITEC-SE` com foco em:

- aderência aos requisitos da `FAPITEC-SE`;
- redução de riscos e custos de manutenção ao longo do tempo;
- uso de DDD, Clean Architecture, Clean Code e TDD;
- preservação da linguagem ubíqua do domínio em português;
- conformidade com padrões operacionais e culturais do Brasil.

## 2. Escopo Inicial do Sistema

Com base no Termo de Referência, o sistema deverá contemplar ao menos os seguintes módulos de negócio:

1. Gestão de Editais
2. Inscrição e Seleção de Projetos
3. Solicitação e Concessão de Auxílios
4. Gestão de Bolsistas e Pesquisadores Vinculados
5. Prestação de Contas
6. Tomada de Contas Especial
7. Comunicação Institucional
8. Relatórios e Painéis de Indicadores

Além desses módulos, o sistema possui capacidades transversais relevantes:

- identidade e controle de acesso;
- gestão documental;
- trilha de auditoria;
- notificações e comunicação;
- segurança da informação e conformidade.

## 3. Diretrizes Arquiteturais

### 3.1 Princípios

- `DDD`: o modelo deve refletir o domínio real da fundação e seus processos.
- `Clean Architecture`: regras de negócio não devem depender de frameworks.
- `Clean Code`: nomes claros, coesão alta e baixo acoplamento.
- `TDD`: priorizar testes de domínio e casos de uso antes da implementação.
- `CLI First`: automações, validações e operações devem existir sem depender de interface gráfica.

### 3.2 Linguagem Ubíqua

Todos os elementos do domínio devem ser escritos em português:

- pastas e arquivos de domínio;
- tipos, structs, interfaces e casos de uso;
- atributos, métodos e eventos;
- nomes de módulos e fluxos de negócio.

Exemplos desejados:

- `edital`, `proponente`, `prestacao_de_contas`, `bolsista`, `tomada_de_contas_especial`

Exemplos a evitar no domínio:

- `grant`, `proposal`, `applicant`, `accountability`, `special_audit_case`

Exceções aceitáveis:

- nomes técnicos de bibliotecas, frameworks, protocolos e padrões externos;
- nomes exigidos por ferramentas de mercado, quando não houver alternativa prática.

### 3.3 Convenções Brasil-First

Todas as regras de apresentação e processamento devem considerar o contexto brasileiro:

- moeda em `BRL`;
- datas no padrão brasileiro na camada de apresentação;
- horário e calendário compatíveis com operação no Brasil;
- percentuais com regra explícita de precisão e arredondamento;
- CPF, passaporte, e-mail e demais documentos conforme regras do negócio;
- linguagem de interface e mensagens em português.

As convenções consolidadas de implementação estão registradas em [docs/architecture/convencoes-arquiteturais-e-linguagem-ubiqua.md](/home/carlos/projetos/fapitec/docs/architecture/convencoes-arquiteturais-e-linguagem-ubiqua.md).

## 4. Stack Inicial Recomendada

### 4.1 Frontend

- `Next.js`
- `React`
- `Tailwind CSS`
- `shadcn/ui`

### 4.2 Backend

- `Go`
- API HTTP orientada a recursos e casos de uso
- `PostgreSQL` como banco de dados relacional principal
- `sqlc` como abordagem preferencial de acesso a dados

### 4.3 Estratégia de Repositório

Monorepo com separação clara entre aplicações, pacotes compartilhados e documentação.

Estrutura inicial recomendada:

```text
apps/
  web/
  api/
packages/
  ui/
  config/
docs/
  architecture/
  stories/
```

### 4.4 Estrutura Física Inicial do Repositório

Detalhamento recomendado para a primeira iteração do monorepo:

```text
apps/
  web/
    src/
      app/
      componentes/
      modulos/
      aplicacao/
      infraestrutura/
      lib/
    public/
  api/
    cmd/
      api/
    internal/
      identidade_e_acesso/
      gestao_de_editais/
      inscricao_e_selecao_de_projetos/
      solicitacao_e_concessao_de_auxilios/
      gestao_de_bolsistas_e_pesquisadores/
      prestacao_de_contas/
      tomada_de_contas_especial/
      comunicacao_institucional/
      relatorios_e_indicadores/
      gestao_documental/
      auditoria/
      notificacoes/
    db/
      consultas/
      gerado/
      migracoes/
packages/
  ui/
  config/
docs/
  architecture/
  stories/
```

Objetivo de cada área:

- `apps/web`: aplicação web da plataforma;
- `apps/api`: backend principal em `Go`;
- `apps/api/db/consultas`: consultas SQL organizadas para geração com `sqlc`;
- `apps/api/db/gerado`: código gerado por `sqlc`;
- `apps/api/db/migracoes`: versionamento de mudanças estruturais do banco;
- `packages/ui`: base de componentes compartilháveis do frontend;
- `packages/config`: convenções compartilhadas do ecossistema frontend;
- `docs`: artefatos de arquitetura, stories e decisões.

## 5. Estrutura Lógica Recomendada

### 5.1 Frontend

```text
apps/web/src/
  app/
  componentes/
  modulos/
  aplicacao/
  infraestrutura/
  lib/
```

Responsabilidades:

- `app`: rotas, layouts e composição de páginas;
- `componentes`: componentes visuais e reutilizáveis;
- `modulos`: organização por contexto funcional;
- `aplicacao`: casos de uso do frontend e orquestração;
- `infraestrutura`: clientes HTTP, adapters e integrações;
- `lib`: utilitários compartilhados.

### 5.2 Backend

```text
apps/api/
  cmd/api/
  internal/
    identidade_e_acesso/
    gestao_de_editais/
    inscricao_e_selecao_de_projetos/
    solicitacao_e_concessao_de_auxilios/
    gestao_de_bolsistas_e_pesquisadores/
    prestacao_de_contas/
    tomada_de_contas_especial/
    comunicacao_institucional/
    relatorios_e_indicadores/
    gestao_documental/
    auditoria/
    notificacoes/
```

Cada contexto deve seguir, sempre que possível, a estrutura:

```text
dominio/
aplicacao/
interfaces/
infraestrutura/
```

Exemplo de detalhamento por contexto:

```text
internal/gestao_de_editais/
  dominio/
    entidades/
    objetos_de_valor/
    servicos/
    repositorios/
  aplicacao/
    casos_de_uso/
    dto/
  interfaces/
    http/
  infraestrutura/
    persistencia/
    http/
```

Regras:

- `dominio`: entidades, objetos de valor, serviços de domínio e contratos;
- `aplicacao`: casos de uso e orquestração;
- `interfaces`: HTTP, filas, serialização e entrada/saída;
- `infraestrutura`: persistência, storage e integrações externas.
- consultas SQL devem ser explicitamente versionadas e utilizadas via `sqlc`, evitando ORM como padrão arquitetural.

### 5.3 Banco de Dados e Acesso a Dados

Estratégia inicial recomendada:

- `PostgreSQL` como banco principal;
- SQL escrito explicitamente e versionado no repositório;
- geração de código via `sqlc` para consultas e comandos;
- separação entre modelo relacional e modelo de domínio;
- evitar vazamento de structs de persistência para o núcleo do domínio.

Estrutura inicial recomendada para persistência:

```text
apps/api/db/
  consultas/
    identidade_e_acesso/
    gestao_de_editais/
    inscricao_e_selecao_de_projetos/
    solicitacao_e_concessao_de_auxilios/
    gestao_de_bolsistas_e_pesquisadores/
    prestacao_de_contas/
    tomada_de_contas_especial/
    comunicacao_institucional/
    relatorios_e_indicadores/
  migracoes/
  gerado/
```

Diretriz:

- cada contexto deve possuir consultas SQL próprias ou claramente agrupadas;
- o código gerado por `sqlc` deve ser tratado como detalhe de infraestrutura;
- o domínio deve depender de interfaces de repositório, não de structs geradas diretamente.

## 6. Decisões de Modelagem

### 6.1 Valores monetários

Em `Go`, valores monetários nunca devem usar `float32` ou `float64`.

Diretriz inicial:

- armazenar valores monetários em centavos, usando tipo inteiro;
- encapsular dinheiro em objeto de valor próprio, por exemplo `ValorMonetario`;
- centralizar regras de arredondamento e formatação;
- converter para exibição em `pt-BR` apenas nas bordas do sistema.

### 6.2 Percentuais

Percentuais não devem depender de arredondamento implícito.

Diretriz inicial:

- usar representação explícita de precisão;
- concentrar cálculo em objetos de valor ou serviços específicos;
- documentar a política de arredondamento por caso de uso.

### 6.3 Data e hora

Diretriz inicial:

- persistir valores temporais de forma padronizada;
- tratar timezone explicitamente;
- exibir datas e horários no formato brasileiro na camada de apresentação;
- evitar regras sensíveis a calendário espalhadas em componentes de interface.

## 7. Identidade e Acesso

O sistema exigirá, no mínimo:

- autenticação;
- autorização por perfis;
- trilha de auditoria;
- controle de acesso a módulos e operações;
- tratamento compatível com LGPD e segurança da informação.

Ferramentas em avaliação:

- `Casdoor` para autenticação e gestão de identidades;
- `Casbin` para autorização, com suporte a RBAC e ABAC.

Neste momento, a decisão ainda está em aberto e deve ser tratada como avaliação arquitetural formal, registrada em ADR.

## 8.1 Plataforma e Restrições Operacionais

O Termo de Referência indica restrições relevantes sobre infraestrutura e cronograma de execução.

Diretrizes iniciais:

- a hospedagem deve respeitar o ambiente e o local exigidos pela `FAPITEC-SE` no Termo de Referência;
- o planejamento das entregas deve considerar cronograma mensal;
- decisões de infraestrutura, segurança, observabilidade e implantação devem ser compatíveis com essas restrições contratuais.

## 9. Mapeamento Inicial de Contextos

### 9.1 Contextos de negócio

- `gestao_de_editais`
- `inscricao_e_selecao_de_projetos`
- `solicitacao_e_concessao_de_auxilios`
- `gestao_de_bolsistas_e_pesquisadores`
- `prestacao_de_contas`
- `tomada_de_contas_especial`
- `comunicacao_institucional`
- `relatorios_e_indicadores`

### 9.2 Contextos transversais

- `identidade_e_acesso`
- `gestao_documental`
- `auditoria`
- `notificacoes`

## 10. Estratégia de Qualidade

### 10.1 Testes

Pirâmide inicial recomendada:

- testes de domínio;
- testes de casos de uso;
- testes de integração de adapters;
- testes E2E para fluxos críticos.

### 10.2 TDD

Fluxo recomendado:

1. escrever o teste do comportamento esperado;
2. implementar o mínimo para passar;
3. refatorar mantendo clareza, isolamento e intenção do domínio.

### 10.3 Qualidade contínua

Os gates atuais do repositório cobrem o ecossistema Node. Para o projeto alvo, a tendência é ampliar para incluir também verificações do backend Go.

Direção futura:

- `npm run lint`
- `npm run typecheck`
- `npm test`
- `go test ./...`
- lint do ecossistema Go

## 11. Primeiras Stories Fundacionais

Sugestão de sequência inicial:

1. Estrutura base do monorepo
2. Convenções arquiteturais e linguagem ubíqua
3. Autenticação e cadastro inicial
4. Autorização por perfis
5. Trilha de auditoria
6. Shell inicial e dashboard modular
7. Gestão documental base
8. Gestão de editais

### 11.1 Estratégia de Sequenciamento

As capacidades transversais devem ser tratadas como alicerce evolutivo do sistema.

Diretriz:

- construir primeiro a base mínima reutilizável;
- evitar tentar esgotar todas as capacidades transversais antes de entregar valor de negócio;
- evoluir autenticação, autorização, auditoria, documentação e notificações em fatias compatíveis com os primeiros módulos.

Sequência recomendada em ondas:

1. `Onda 1 - Fundação mínima`
   estrutura base do monorepo, convenções arquiteturais, autenticação inicial, autorização inicial, auditoria mínima e shell inicial com dashboard modular.
2. `Onda 2 - Núcleo operacional`
   gestão documental base, notificações base e módulo de gestão de editais.
3. `Onda 3 - Entrada e seleção`
   inscrição e seleção de projetos.
4. `Onda 4 - Execução`
   solicitação e concessão de auxílios, gestão de bolsistas e pesquisadores.
5. `Onda 5 - Controle e conformidade`
   prestação de contas, tomada de contas especial, relatórios e indicadores.

Critério de aplicação:

- nenhuma capacidade transversal deve ser construída em nível "completo" antes dos primeiros módulos;
- cada capacidade deve ser implementada no menor nível suficiente para suportar a próxima entrega de negócio;
- a expansão posterior deve acontecer guiada pelas stories dos módulos.
- o dashboard inicial deve apresentar a macroestrutura dos módulos como placeholders validáveis, sem detalhar regras de negócio ainda não confirmadas pela `FAPITEC-SE`.

## 12. Riscos e Pontos em Aberto

- o material do pedido de esclarecimento indica que parte dos fluxos ainda é esboço inicial;
- a decisão final sobre `Casdoor` e `Casbin` depende de avaliação formal;
- ainda será necessário detalhar modelagem de banco, storage, observabilidade e estratégia de implantação;
- será necessário confirmar timezone operacional, integrações externas e perfis institucionais;
- a documentação de histórias ainda precisa ser consolidada em artefatos executáveis do AIOX.

## 13. Próximos Passos

1. Avaliar formalmente a estratégia de identidade e acesso.
2. Implementar autenticação, autorização e auditoria mínimas.
3. Implementar shell inicial com dashboard modular para validação da macroestrutura pela `FAPITEC-SE`.
3. Definir glossário oficial do domínio em português.
4. Estruturar o monorepo de acordo com os contextos definidos.
