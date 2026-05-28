#!/usr/bin/env node

const fs = require('node:fs');
const path = require('node:path');
const { spawnSync } = require('node:child_process');

const raiz = process.cwd();
const candidate = [
  path.join(raiz, 'node_modules', 'typescript', 'bin', 'tsc'),
  path.join(raiz, 'apps', 'web', 'node_modules', 'typescript', 'bin', 'tsc'),
];

const tsc = candidate.find((p) => fs.existsSync(p));

if (!tsc) {
  console.log('typecheck: typescript ainda não instalado.');
  console.log('typecheck: execute pnpm install e depois pnpm --filter @fapitec/web run typecheck.');
  process.exit(0);
}

const resultado = spawnSync('pnpm', ['--filter', '@fapitec/web', 'run', 'typecheck'], {
  cwd: raiz,
  stdio: 'inherit'
});

process.exit(resultado.status ?? 1);
