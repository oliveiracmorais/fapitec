#!/usr/bin/env node

const fs = require('node:fs');
const path = require('node:path');

const RAIZ = process.cwd();
const ORIGEM = '.aiox-core/development/agents';
const DESTINO = '.codex/agents';
const AGENTES = [
  'aiox-master', 'analyst', 'architect', 'data-engineer', 'dev', 'devops',
  'pm', 'po', 'qa', 'sm', 'squad-creator', 'ux-design-expert'
];
const MODO_CHECK = process.argv.includes('--check');

function syncAgente(nome) {
  const origemPath = path.join(RAIZ, ORIGEM, `${nome}.md`);
  const destinoPath = path.join(RAIZ, DESTINO, `${nome}.md`);

  if (!fs.existsSync(origemPath)) {
    return { nome, status: 'origem_ausente' };
  }

  const conteudo = fs.readFileSync(origemPath, 'utf-8');
  const conteudoDestino = conteudo.endsWith('\n')
    ? conteudo + `\n*AIOX Agent - Synced from .aiox-core/development/agents/${nome}.md*\n`
    : conteudo + `\n\n*AIOX Agent - Synced from .aiox-core/development/agents/${nome}.md*\n`;

  if (MODO_CHECK) {
    if (!fs.existsSync(destinoPath)) {
      return { nome, status: 'faltando' };
    }
    const atual = fs.readFileSync(destinoPath, 'utf-8');
    if (atual !== conteudoDestino) {
      return { nome, status: 'desatualizado' };
    }
    return { nome, status: 'ok' };
  }

  fs.mkdirSync(path.dirname(destinoPath), { recursive: true });
  fs.writeFileSync(destinoPath, conteudoDestino);
  return { nome, status: 'sincronizado' };
}

const resultados = AGENTES.map(syncAgente);
const erros = resultados.filter(r => r.status !== 'ok');

if (MODO_CHECK) {
  if (erros.length > 0) {
    console.error('Agentes desatualizados ou faltando:');
    for (const e of erros) {
      console.error(`  - ${e.nome}: ${e.status}`);
    }
    process.exit(1);
  }
  console.log(`${resultados.length} agentes sincronizados e atualizados.`);
} else {
  for (const r of resultados) {
    console.log(`  ${r.nome}: ${r.status}`);
  }
  if (erros.length > 0) process.exit(1);
}
