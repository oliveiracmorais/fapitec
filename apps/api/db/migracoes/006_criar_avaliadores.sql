-- Migration 006: Criar tabelas de avaliadores e atribuicoes
-- Depende: 005_criar_propostas.sql

CREATE TABLE IF NOT EXISTS avaliadores (
    id BIGSERIAL PRIMARY KEY,
    usuario_id BIGINT NOT NULL,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL DEFAULT '',
    titulacao_maxima VARCHAR(100) NOT NULL DEFAULT '',
    area_conhecimento VARCHAR(255) NOT NULL DEFAULT '',
    instituicao VARCHAR(255) NOT NULL DEFAULT '',
    curriculo_resumido TEXT NOT NULL DEFAULT '',
    estado VARCHAR(20) NOT NULL DEFAULT 'ativo',
    data_cadastro TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    data_atualizacao TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS atribuicoes_editais (
    id BIGSERIAL PRIMARY KEY,
    avaliador_id BIGINT NOT NULL REFERENCES avaliadores(id) ON DELETE RESTRICT,
    edital_id BIGINT NOT NULL REFERENCES editais(id) ON DELETE RESTRICT,
    data_inicio TIMESTAMPTZ NOT NULL,
    data_fim TIMESTAMPTZ NOT NULL,
    status_convite VARCHAR(20) NOT NULL DEFAULT 'pendente',
    hash_anonimizacao VARCHAR(64) NOT NULL UNIQUE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_avaliadores_usuario_id ON avaliadores(usuario_id);
CREATE INDEX IF NOT EXISTS idx_avaliadores_cpf ON avaliadores(cpf);
CREATE INDEX IF NOT EXISTS idx_avaliadores_area_conhecimento ON avaliadores(area_conhecimento);
CREATE INDEX IF NOT EXISTS idx_avaliadores_estado ON avaliadores(estado);
CREATE INDEX IF NOT EXISTS idx_atribuicoes_avaliador_id ON atribuicoes_editais(avaliador_id);
CREATE INDEX IF NOT EXISTS idx_atribuicoes_edital_id ON atribuicoes_editais(edital_id);
CREATE INDEX IF NOT EXISTS idx_atribuicoes_status_convite ON atribuicoes_editais(status_convite);
