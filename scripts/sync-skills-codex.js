#!/usr/bin/env node

const modoGlobal = process.argv.includes('--global');

if (modoGlobal) {
  console.log('Sincronizacao global de skills do Codex reservada para configuracao explicita futura.');
} else {
  console.log('Sincronizacao local de skills do Codex reservada para configuracao explicita futura.');
}
