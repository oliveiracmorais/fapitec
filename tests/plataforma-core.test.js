const test = require('node:test');
const assert = require('node:assert/strict');
const {
  buscarModulo,
  criarSessaoDemonstrativa,
  criarEventoAuditoria,
  listarIdentidadesDemonstrativas,
  listarModulos,
  listarPerfis,
  obterSessaoDemonstrativa,
  podeAcessarModulo
} = require('../scripts/plataforma-core');

test('catálogo contempla módulos oficiais do dashboard inicial', () => {
  const ids = listarModulos({ incluirAdministrativos: false }).map((modulo) => modulo.id);

  assert.deepEqual(ids, [
    'gestao_de_editais',
    'inscricao_e_selecao_de_projetos',
    'solicitacao_e_concessao_de_auxilios',
    'gestao_de_bolsistas_e_pesquisadores',
    'prestacao_de_contas',
    'tomada_de_contas_especial',
    'comunicacao_institucional',
    'relatorios_e_indicadores',
    'gestao_documental'
  ]);
});

test('perfil administrador acessa módulos administrativos', () => {
  assert.equal(podeAcessarModulo('administrador', 'identidade_e_acesso'), true);
  assert.equal(podeAcessarModulo('administrador', 'auditoria'), true);
});

test('perfil proponente não acessa área administrativa de auditoria', () => {
  assert.equal(podeAcessarModulo('proponente', 'auditoria'), false);
});

test('catálogo possui perfis mínimos definidos na story 1.4', () => {
  const perfis = listarPerfis().map((perfil) => perfil.id);

  assert.deepEqual(perfis, [
    'administrador',
    'gestor_fapitec',
    'avaliador',
    'proponente',
    'auditoria_controle'
  ]);
});

test('catálogo possui identidades demonstrativas sem dados pessoais reais', () => {
  const identidades = listarIdentidadesDemonstrativas();

  assert.equal(identidades.length, 5);
  assert.equal(identidades.some((identidade) => identidade.id === 'gestor_validacao'), true);
  assert.equal(identidades.every((identidade) => !String(identidade.id).includes('@')), true);
});

test('sessão demonstrativa deriva perfil da identidade controlada', () => {
  const sessao = criarSessaoDemonstrativa('gestor_validacao');

  assert.equal(sessao.id, 'sessao-gestor_validacao');
  assert.equal(sessao.perfil, 'gestor_fapitec');
  assert.equal(sessao.ator, 'identidade-demonstrativa:gestor_validacao');
  assert.equal(obterSessaoDemonstrativa(sessao.id).perfil, 'gestor_fapitec');
});

test('evento de auditoria registra ação, resultado e contexto do módulo', () => {
  const evento = criarEventoAuditoria({
    ator: 'usuario.validacao@fapitec.se.gov.br',
    perfil: 'gestor_fapitec',
    acao: 'acessar_modulo_dashboard',
    moduloId: 'gestao_de_editais',
    resultado: 'sucesso'
  });

  assert.equal(evento.acao, 'acessar_modulo_dashboard');
  assert.equal(evento.resultado, 'sucesso');
  assert.equal(evento.modulo, 'gestao_de_editais');
  assert.equal(evento.contexto.nomeModulo, buscarModulo('gestao_de_editais').nome);
  assert.match(evento.dataHora, /^\d{4}-\d{2}-\d{2}T/);
});

test('evento de auditoria aceita contexto de sessão sem registrar segredo', () => {
  const evento = criarEventoAuditoria({
    ator: 'identidade-demonstrativa:proponente_validacao',
    perfil: 'proponente',
    acao: 'acessar_modulo_dashboard',
    moduloId: 'auditoria',
    resultado: 'negado',
    contexto: { sessao: 'sessao-proponente_validacao', motivo: 'perfil_sem_permissao' }
  });

  assert.equal(evento.resultado, 'negado');
  assert.equal(evento.contexto.sessao, 'sessao-proponente_validacao');
  assert.equal(JSON.stringify(evento).includes('senha'), false);
  assert.equal(JSON.stringify(evento).includes('token'), false);
});
