#!/usr/bin/env node

const fs = require('node:fs');
const path = require('node:path');
const { diretoriosObrigatorios, arquivosObrigatorios } = require('./estrutura-esperada');

const raiz = process.cwd();

function existe(relativo) {
  return fs.existsSync(path.join(raiz, relativo));
}

const faltantes = [
  ...diretoriosObrigatorios.filter((diretorio) => !existe(diretorio)),
  ...arquivosObrigatorios.filter((arquivo) => !existe(arquivo))
];

if (faltantes.length > 0) {
  console.error('Estrutura do projeto incompleta:');
  for (const item of faltantes) {
    console.error(`- ${item}`);
  }
  process.exit(1);
}

console.log('Estrutura do projeto validada com sucesso.');
