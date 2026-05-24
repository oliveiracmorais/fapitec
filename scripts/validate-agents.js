#!/usr/bin/env node

const fs = require('node:fs');
const path = require('node:path');

const agentes = [
  'aiox-master',
  'analyst',
  'architect',
  'data-engineer',
  'dev',
  'devops',
  'pm',
  'po',
  'qa',
  'sm',
  'squad-creator',
  'ux-design-expert'
];

const raiz = process.cwd();
const faltantes = [];

for (const agente of agentes) {
  const caminho = path.join(raiz, '.codex/agents', `${agente}.md`);
  if (!fs.existsSync(caminho)) {
    faltantes.push(agente);
  }
}

if (faltantes.length > 0) {
  console.error('Configuracao de agentes incompleta:');
  for (const item of faltantes) {
    console.error(`- ${item}`);
  }
  process.exit(1);
}

console.log('Configuracao de agentes validada com sucesso.');
