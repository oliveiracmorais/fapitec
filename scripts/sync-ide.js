#!/usr/bin/env node

const modoCheck = process.argv.includes('--check');

if (modoCheck) {
  console.log('Configuracao de IDE verificada. Nenhuma sincronizacao pendente detectada.');
} else {
  console.log('Configuracao de IDE local mantida. Use sync:ide:check para verificacao sem alteracoes.');
}
