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
const bases = ['.codex/agents', '.claude/rules', '.github/agents'];
const faltantes = [];

for (const agente of agentes) {
  const caminhosPossiveis = [
    path.join(raiz, '.codex/agents', `${agente}.md`),
    path.join(raiz, '.github/agents', `${agente}.agent.md`)
  ];

  if (!caminhosPossiveis.some((caminho) => fs.existsSync(caminho))) {
    faltantes.push(agente);
  }
}

for (const base of bases) {
  if (!fs.existsSync(path.join(raiz, base))) {
    faltantes.push(base);
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
