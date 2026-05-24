#!/usr/bin/env node

const http = require('node:http');

const {
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
} = require('./plataforma-core');

const API_BASE = process.env.API_URL || 'http://localhost:8080';
const [comando, ...args] = process.argv.slice(2);

function imprimirAjuda() {
  console.log(`Uso:
  pnpm run shell:modulos
  pnpm run shell:perfis
  node scripts/plataforma-shell.js identidades
  node scripts/plataforma-shell.js sessao <identidade>
  node scripts/plataforma-shell.js acesso <perfil> <modulo>
  node scripts/plataforma-shell.js acesso-sessao <sessao> <modulo>
  node scripts/plataforma-shell.js auditoria <perfil> <modulo>
  node scripts/plataforma-shell.js auditoria-sessao <sessao> <modulo>
   node scripts/plataforma-shell.js cadastrar <nome> <cpf> <email> <senha> [estrangeiro]
   node scripts/plataforma-shell.js login <cpf> <senha>
   node scripts/plataforma-shell.js solicitar-redefinicao <email>
   node scripts/plataforma-shell.js redefinir-senha <token> <senha>`);
}

function imprimirModulos() {
  for (const modulo of listarModulos()) {
    console.log(`${modulo.id} | ${modulo.nome} | ${modulo.status}`);
  }
}

function imprimirPerfis() {
  for (const perfil of listarPerfis()) {
    console.log(`${perfil.id} | ${perfil.nome}`);
  }
}

function imprimirIdentidadesDemonstrativas() {
  for (const identidade of listarIdentidadesDemonstrativas()) {
    const perfil = buscarPerfil(identidade.perfil);
    console.log(`${identidade.id} | ${identidade.nome} | ${perfil ? perfil.nome : identidade.perfil}`);
  }
}

function imprimirSessao(identidadeId) {
  const identidade = buscarIdentidadeDemonstrativa(identidadeId);

  if (!identidade) {
    console.error(`Identidade demonstrativa não encontrada: ${identidadeId}`);
    process.exit(1);
  }

  console.log(JSON.stringify(criarSessaoDemonstrativa(identidade.id), null, 2));
}

function verificarAcesso(perfilId, moduloId) {
  const perfil = buscarPerfil(perfilId);
  const modulo = buscarModulo(moduloId);

  if (!perfil) {
    console.error(`Perfil não encontrado: ${perfilId}`);
    process.exit(1);
  }

  if (!modulo) {
    console.error(`Módulo não encontrado: ${moduloId}`);
    process.exit(1);
  }

  const permitido = podeAcessarModulo(perfil.id, modulo.id);
  console.log(`${perfil.nome} -> ${modulo.nome}: ${permitido ? 'permitido' : 'negado'}`);
  process.exit(permitido ? 0 : 2);
}

function verificarAcessoPorSessao(sessaoId, moduloId) {
  const sessao = obterSessaoDemonstrativa(sessaoId);
  const modulo = buscarModulo(moduloId);

  if (!sessao) {
    console.error(`Sessão demonstrativa não encontrada: ${sessaoId}`);
    process.exit(1);
  }

  if (!modulo) {
    console.error(`Módulo não encontrado: ${moduloId}`);
    process.exit(1);
  }

  const permitido = podeAcessarModulo(sessao.perfil, modulo.id);
  console.log(`${sessao.nome} (${sessao.perfilNome}) -> ${modulo.nome}: ${permitido ? 'permitido' : 'negado'}`);
  process.exit(permitido ? 0 : 2);
}

function simularAuditoria(perfilId, moduloId) {
  const permitido = podeAcessarModulo(perfilId, moduloId);
  const evento = criarEventoAuditoria({
    ator: 'usuario.validacao@fapitec.se.gov.br',
    perfil: perfilId,
    acao: 'acessar_modulo_dashboard',
    moduloId,
    resultado: permitido ? 'sucesso' : 'negado'
  });

  console.log(JSON.stringify(evento, null, 2));
  process.exit(permitido ? 0 : 2);
}

function chamarAPI(method, path, body) {
  return new Promise((resolve, reject) => {
    const url = new URL(path, API_BASE);
    const opts = {
      method,
      hostname: url.hostname,
      port: url.port,
      path: url.pathname,
      headers: body ? { 'Content-Type': 'application/json' } : {},
      timeout: 5000,
    };
    const req = http.request(opts, (res) => {
      let data = '';
      res.on('data', (chunk) => (data += chunk));
      res.on('end', () => {
        resolve({ status: res.statusCode, body: data });
      });
    });
    req.on('error', () => reject(new Error('API nao disponivel. Execute a API Go primeiro.')));
    req.on('timeout', () => { req.destroy(); reject(new Error('Timeout ao conectar na API.')); });
    if (body) req.write(JSON.stringify(body));
    req.end();
  });
}

async function cadastrarUsuario(nome, cpf, email, senha, estrangeiro) {
  try {
    const res = await chamarAPI('POST', '/api/v1/cadastro', {
      nome, cpf, email, confirmacao_email: email, senha, confirmacao_senha: senha,
      estrangeiro: estrangeiro === 'true',
    });
    if (res.status === 201) {
      const dados = JSON.parse(res.body);
      console.log(`Usuario cadastrado: ${dados.nome} (ID: ${dados.id})`);
      process.exit(0);
    } else {
      const err = JSON.parse(res.body);
      console.error(`Erro: ${err.erro}`);
      process.exit(1);
    }
  } catch (e) {
    console.error(e.message);
    process.exit(1);
  }
}

async function loginUsuario(cpf, senha) {
  try {
    const res = await chamarAPI('POST', '/api/v1/login', { cpf, senha });
    if (res.status === 200) {
      const dados = JSON.parse(res.body);
      console.log(`Autenticado: ${dados.nome} (${dados.email})`);
      process.exit(0);
    } else {
      const err = JSON.parse(res.body);
      console.error(`Erro: ${err.erro}`);
      process.exit(1);
    }
  } catch (e) {
    console.error(e.message);
    process.exit(1);
  }
}

async function solicitarRedefinicaoSenha(email) {
  try {
    const res = await chamarAPI('POST', '/api/v1/solicitar-redefinicao-senha', { email });
    const dados = JSON.parse(res.body);
    console.log(dados.mensagem);
    process.exit(0);
  } catch (e) {
    console.error(e.message);
    process.exit(1);
  }
}

async function redefinirSenha(token, senha) {
  try {
    const res = await chamarAPI('POST', '/api/v1/redefinir-senha', { token, senha, confirmacao_senha: senha });
    if (res.status === 200) {
      const dados = JSON.parse(res.body);
      console.log(dados.mensagem);
      process.exit(0);
    } else {
      const err = JSON.parse(res.body);
      console.error(`Erro: ${err.erro}`);
      process.exit(1);
    }
  } catch (e) {
    console.error(e.message);
    process.exit(1);
  }
}

function simularAuditoriaPorSessao(sessaoId, moduloId) {
  const sessao = obterSessaoDemonstrativa(sessaoId);

  if (!sessao) {
    const evento = criarEventoAuditoria({
      ator: 'sistema',
      acao: 'acessar_modulo_dashboard',
      moduloId,
      resultado: 'negado',
      contexto: { motivo: 'sessao_ausente_ou_invalida' }
    });
    console.log(JSON.stringify(evento, null, 2));
    process.exit(2);
  }

  const permitido = podeAcessarModulo(sessao.perfil, moduloId);
  const evento = criarEventoAuditoria({
    ator: sessao.ator,
    perfil: sessao.perfil,
    acao: 'acessar_modulo_dashboard',
    moduloId,
    resultado: permitido ? 'sucesso' : 'negado',
    contexto: { sessao: sessao.id, identidadeDemonstrativa: sessao.identidadeId }
  });

  console.log(JSON.stringify(evento, null, 2));
  process.exit(permitido ? 0 : 2);
}

switch (comando) {
  case 'modulos':
    imprimirModulos();
    break;
  case 'perfis':
    imprimirPerfis();
    break;
  case 'identidades':
    imprimirIdentidadesDemonstrativas();
    break;
  case 'sessao':
    imprimirSessao(args[0]);
    break;
  case 'acesso':
    verificarAcesso(args[0], args[1]);
    break;
  case 'acesso-sessao':
    verificarAcessoPorSessao(args[0], args[1]);
    break;
  case 'auditoria':
    simularAuditoria(args[0], args[1]);
    break;
  case 'auditoria-sessao':
    simularAuditoriaPorSessao(args[0], args[1]);
    break;
  case 'cadastrar':
    cadastrarUsuario(args[0], args[1], args[2], args[3], args[4]);
    break;
  case 'login':
    loginUsuario(args[0], args[1]);
    break;
  case 'solicitar-redefinicao':
    solicitarRedefinicaoSenha(args[0]);
    break;
  case 'redefinir-senha':
    redefinirSenha(args[0], args[1]);
    break;
  default:
    imprimirAjuda();
    process.exit(comando ? 1 : 0);
}
