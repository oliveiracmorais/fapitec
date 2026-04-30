const fs = require('node:fs');
const path = require('node:path');
const test = require('node:test');
const assert = require('node:assert/strict');
const { diretoriosObrigatorios, arquivosObrigatorios } = require('../scripts/estrutura-esperada');

const raiz = path.resolve(__dirname, '..');

test('estrutura fundacional do monorepo existe', () => {
  for (const diretorio of diretoriosObrigatorios) {
    assert.equal(fs.existsSync(path.join(raiz, diretorio)), true, diretorio);
  }

  for (const arquivo of arquivosObrigatorios) {
    assert.equal(fs.existsSync(path.join(raiz, arquivo)), true, arquivo);
  }
});
