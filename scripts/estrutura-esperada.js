const diretoriosObrigatorios = [
  'apps/web/src/app',
  'apps/web/src/componentes',
  'apps/web/src/modulos',
  'apps/web/src/aplicacao',
  'apps/web/src/infraestrutura',
  'apps/web/src/lib',
  'apps/web/public',
  'apps/api/cmd/api',
  'apps/api/internal/identidade_e_acesso',
  'apps/api/internal/gestao_de_editais',
  'apps/api/internal/inscricao_e_selecao_de_projetos',
  'apps/api/internal/solicitacao_e_concessao_de_auxilios',
  'apps/api/internal/gestao_de_bolsistas_e_pesquisadores',
  'apps/api/internal/prestacao_de_contas',
  'apps/api/internal/tomada_de_contas_especial',
  'apps/api/internal/comunicacao_institucional',
  'apps/api/internal/relatorios_e_indicadores',
  'apps/api/internal/gestao_documental',
  'apps/api/internal/auditoria',
  'apps/api/internal/notificacoes',
  'apps/api/db/consultas',
  'apps/api/db/gerado',
  'apps/api/db/migracoes',
  'packages/ui',
  'packages/config',
  'docs/architecture',
  'docs/stories',
  'tests',
  'scripts'
];

const arquivosObrigatorios = [
  'AGENTS.md',
  '.aiox-core/constitution.md',
  'docs/architecture.md',
  'docs/architecture/adr-001-identidade-e-acesso.md',
  'docs/architecture/convencoes-arquiteturais-e-linguagem-ubiqua.md',
  'docs/stories/README.md',
  'docs/stories/1.1.estrutura-base-do-monorepo.md',
  'docs/stories/1.2.convencoes-arquiteturais-e-linguagem-ubiqua.md',
  'docs/stories/1.3.autenticacao-e-cadastro-inicial.md',
  'docs/stories/1.4.autorizacao-por-perfis.md',
  'docs/stories/1.5.trilha-de-auditoria-minima.md',
  'docs/stories/1.6.shell-inicial-e-dashboard-modular.md',
  'package.json'
];

module.exports = {
  diretoriosObrigatorios,
  arquivosObrigatorios
};
