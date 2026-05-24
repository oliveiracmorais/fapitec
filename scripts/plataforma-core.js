const catalogo = require('../packages/config/src/catalogo-plataforma.json');

function listarPerfis() {
  return catalogo.perfis;
}

function listarIdentidadesDemonstrativas() {
  return catalogo.identidadesDemonstrativas;
}

function listarModulos({ incluirAdministrativos = true } = {}) {
  return catalogo.modulos.filter((modulo) => incluirAdministrativos || modulo.tipo !== 'administrativo');
}

function buscarModulo(id) {
  return catalogo.modulos.find((modulo) => modulo.id === id);
}

function buscarPerfil(id) {
  return catalogo.perfis.find((perfil) => perfil.id === id);
}

function buscarIdentidadeDemonstrativa(id) {
  return catalogo.identidadesDemonstrativas.find((identidade) => identidade.id === id);
}

function criarSessaoDemonstrativa(identidadeId) {
  const identidade = buscarIdentidadeDemonstrativa(identidadeId);

  if (!identidade) {
    return undefined;
  }

  const perfil = buscarPerfil(identidade.perfil);

  return {
    id: `sessao-${identidade.id}`,
    ator: `identidade-demonstrativa:${identidade.id}`,
    identidadeId: identidade.id,
    nome: identidade.nome,
    perfil: perfil ? perfil.id : identidade.perfil,
    perfilNome: perfil ? perfil.nome : identidade.perfil,
    demonstrativa: true
  };
}

function obterSessaoDemonstrativa(sessaoId) {
  if (!sessaoId || typeof sessaoId !== 'string') {
    return undefined;
  }

  const identidadeId = sessaoId.startsWith('sessao-') ? sessaoId.slice('sessao-'.length) : sessaoId;
  return criarSessaoDemonstrativa(identidadeId);
}

function podeAcessarModulo(perfilId, moduloId) {
  const perfil = buscarPerfil(perfilId);
  const modulo = buscarModulo(moduloId);

  if (!perfil || !modulo) {
    return false;
  }

  return modulo.perfisPermitidos.includes(perfil.id);
}

function criarEventoAuditoria({ ator = 'sistema', perfil, acao, moduloId, resultado, contexto = {} }) {
  const modulo = moduloId ? buscarModulo(moduloId) : undefined;

  return {
    id: `evt-${Date.now()}`,
    ator,
    perfil,
    acao,
    resultado,
    modulo: modulo ? modulo.id : undefined,
    contexto: {
      ...(modulo ? { nomeModulo: modulo.nome, tipoModulo: modulo.tipo } : {}),
      ...contexto
    },
    dataHora: new Date().toISOString()
  };
}

module.exports = {
  buscarIdentidadeDemonstrativa,
  buscarModulo,
  buscarPerfil,
  criarSessaoDemonstrativa,
  criarEventoAuditoria,
  listarIdentidadesDemonstrativas,
  listarModulos,
  listarPerfis,
  obterSessaoDemonstrativa,
  podeAcessarModulo
};
