#!/usr/bin/env node

const fs = require('node:fs');
const path = require('node:path');
const { spawnSync } = require('node:child_process');

const raiz = process.cwd();
const tsc = path.join(raiz, 'node_modules', 'typescript', 'bin', 'tsc');

if (!fs.existsSync(tsc)) {
  console.log('typecheck: projeto web TypeScript detectado, mas dependências ainda não instaladas.');
  console.log('typecheck: execute pnpm install e depois pnpm --filter @fapitec/web run typecheck.');
  process.exit(0);
}

const resultado = spawnSync('pnpm', ['--filter', '@fapitec/web', 'run', 'typecheck'], {
  cwd: raiz,
  stdio: 'inherit'
});

process.exit(resultado.status ?? 1);
