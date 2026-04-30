# Convenções Arquiteturais e Linguagem Ubíqua

Este documento consolida as convenções mínimas para implementação da plataforma modular da `FAPITEC-SE`.

As regras abaixo derivam dos artefatos oficiais em `docs/stories/`, do documento de arquitetura inicial e das instruções do projeto em `AGENTS.md`.

## 1. Princípios De Implementação

- `DDD`: o modelo deve expressar o domínio real da `FAPITEC-SE`.
- `Clean Architecture`: regras de negócio não devem depender de framework, banco, HTTP, UI ou ferramentas externas.
- `Clean Code`: nomes claros, coesão alta, baixo acoplamento e responsabilidades explícitas.
- `TDD`: priorizar testes de domínio e casos de uso antes de adapters, integração e UI.
- `CLI First`: automações e verificações devem existir sem depender de interface gráfica.

## 2. Linguagem Ubíqua Em Português

Elementos de domínio devem usar português institucional:

- pastas e arquivos de domínio;
- tipos, structs, interfaces, atributos, métodos e eventos;
- casos de uso, módulos e fluxos de negócio;
- nomes de permissões e ações auditáveis quando representarem conceitos internos do sistema.

Termos preferenciais incluem:

- `edital`
- `proponente`
- `bolsista`
- `prestacao_de_contas`
- `tomada_de_contas_especial`
- `identidade_e_acesso`
- `gestao_documental`
- `auditoria`
- `notificacoes`

Termos genéricos em inglês devem ser evitados quando houver termo institucional claro em português. Exceções são permitidas para nomes técnicos de bibliotecas, protocolos, padrões externos e contratos exigidos por ferramentas.

## 3. Convenções Brasil-First

O sistema deve considerar o contexto operacional brasileiro:

- moeda em `BRL`;
- datas exibidas em padrão brasileiro na camada de apresentação;
- horários tratados com timezone explícito e exibidos conforme operação no Brasil;
- mensagens de interface e validação em português;
- CPF, passaporte, e-mail e demais documentos conforme regras do negócio;
- percentuais com precisão e arredondamento definidos pelo caso de uso.

Valores persistidos devem ser padronizados nas bordas apropriadas. A apresentação deve formatar datas, horários, moeda e percentuais em `pt-BR`.

## 4. Valores Monetários

No backend `Go`, valores monetários nunca devem usar `float32` ou `float64`.

Diretriz inicial:

- representar dinheiro em centavos com tipo inteiro, quando suficiente;
- encapsular regras monetárias em objeto de valor do domínio, por exemplo `ValorMonetario`;
- centralizar arredondamento, soma, comparação e formatação;
- converter para exibição em `BRL` apenas nas bordas do sistema.

## 5. Percentuais

Percentuais devem ter precisão explícita.

Diretriz inicial:

- evitar arredondamento implícito;
- documentar a política de arredondamento por caso de uso;
- encapsular cálculos sensíveis em objetos de valor ou serviços de domínio;
- não espalhar regra de percentual em componentes de UI ou adapters.

## 6. Persistência E Banco De Dados

Diretriz inicial:

- `PostgreSQL` é o banco relacional principal;
- `sqlc` é a estratégia preferencial de acesso a dados no backend `Go`;
- SQL deve ser explícito, versionado e organizado por contexto quando fizer sentido;
- ORM não é padrão do projeto;
- structs geradas por `sqlc` não devem vazar para o núcleo do domínio;
- repositórios do domínio devem ser contratos, implementados por adapters de infraestrutura.

## 7. Organização Por Contexto

Cada contexto do backend deve evoluir, quando fizer sentido, com:

```text
dominio/
aplicacao/
interfaces/
infraestrutura/
```

Responsabilidades:

- `dominio`: entidades, objetos de valor, serviços de domínio e contratos;
- `aplicacao`: casos de uso, orquestração e DTOs de aplicação;
- `interfaces`: HTTP, filas, serialização e entrada/saída;
- `infraestrutura`: persistência, storage, provedores externos e detalhes técnicos.

## 8. Restrições Operacionais

Decisões de infraestrutura, segurança, hospedagem e implantação devem respeitar o Termo de Referência.

O planejamento deve considerar entregas mensais e dependências entre capacidades fundacionais e módulos de negócio.

## 9. Relação Com Identidade E Acesso

`Casdoor` e `Casbin` seguem como avaliação arquitetural, conforme `docs/architecture/adr-001-identidade-e-acesso.md`.

Até decisão final:

- o domínio não deve depender diretamente dessas ferramentas;
- autenticação, autorização e auditoria devem ser tratadas como capacidades transversais;
- integrações externas devem ficar isoladas em adapters;
- regras de autorização e eventos auditáveis devem ser rastreáveis e testáveis.
