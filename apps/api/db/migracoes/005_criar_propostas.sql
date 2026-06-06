-- Migration 005: Criar tabelas do dominio inscricao_e_selecao_de_projetos
-- Depende: 004_estender_editais.sql

CREATE TABLE IF NOT EXISTS propostas (
    id BIGSERIAL PRIMARY KEY,
    edital_id BIGINT NOT NULL REFERENCES editais(id) ON DELETE RESTRICT,
    proponente_id BIGINT NOT NULL,
    protocolo VARCHAR(30) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'rascunho',
    versao INTEGER NOT NULL DEFAULT 1,
    -- Dados do proponente (embedded)
    proponente_nome VARCHAR(255) NOT NULL,
    proponente_cpf VARCHAR(14) NOT NULL,
    proponente_rg VARCHAR(20) NOT NULL DEFAULT '',
    proponente_genero VARCHAR(30) NOT NULL DEFAULT '',
    proponente_etnia VARCHAR(30) NOT NULL DEFAULT '',
    proponente_data_nascimento VARCHAR(10) NOT NULL DEFAULT '',
    proponente_endereco TEXT NOT NULL DEFAULT '',
    proponente_telefone VARCHAR(20) NOT NULL DEFAULT '',
    proponente_email VARCHAR(255) NOT NULL DEFAULT '',
    -- Dados academicos (embedded)
    academico_maior_titulacao VARCHAR(50) NOT NULL DEFAULT '',
    academico_curso VARCHAR(255) NOT NULL DEFAULT '',
    academico_instituicao VARCHAR(255) NOT NULL DEFAULT '',
    academico_ano_conclusao INTEGER NOT NULL DEFAULT 0,
    academico_area_conhecimento VARCHAR(255) NOT NULL DEFAULT '',
    -- Empresa vinculada
    empresa_vinculada VARCHAR(255) NOT NULL DEFAULT '',
    -- Valor total calculado
    valor_total_solicitado BIGINT NOT NULL DEFAULT 0,
    -- Timestamps
    data_submissao TIMESTAMPTZ,
    data_atualizacao TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS versoes_proposta (
    id BIGSERIAL PRIMARY KEY,
    proposta_id BIGINT NOT NULL REFERENCES propostas(id) ON DELETE CASCADE,
    versao INTEGER NOT NULL,
    -- Snapshot dos dados da proposta no momento da versao
    status VARCHAR(20) NOT NULL,
    proponente_nome VARCHAR(255) NOT NULL,
    proponente_cpf VARCHAR(14) NOT NULL,
    proponente_rg VARCHAR(20) NOT NULL DEFAULT '',
    proponente_genero VARCHAR(30) NOT NULL DEFAULT '',
    proponente_etnia VARCHAR(30) NOT NULL DEFAULT '',
    proponente_data_nascimento VARCHAR(10) NOT NULL DEFAULT '',
    proponente_endereco TEXT NOT NULL DEFAULT '',
    proponente_telefone VARCHAR(20) NOT NULL DEFAULT '',
    proponente_email VARCHAR(255) NOT NULL DEFAULT '',
    academico_maior_titulacao VARCHAR(50) NOT NULL DEFAULT '',
    academico_curso VARCHAR(255) NOT NULL DEFAULT '',
    academico_instituicao VARCHAR(255) NOT NULL DEFAULT '',
    academico_ano_conclusao INTEGER NOT NULL DEFAULT 0,
    academico_area_conhecimento VARCHAR(255) NOT NULL DEFAULT '',
    empresa_vinculada VARCHAR(255) NOT NULL DEFAULT '',
    valor_total_solicitado BIGINT NOT NULL DEFAULT 0,
    protocolo VARCHAR(30) NOT NULL DEFAULT '',
    data_submissao TIMESTAMPTZ,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS itens_orcamentarios (
    id BIGSERIAL PRIMARY KEY,
    proposta_id BIGINT NOT NULL REFERENCES propostas(id) ON DELETE CASCADE,
    descricao TEXT NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    quantidade INTEGER NOT NULL DEFAULT 1,
    valor_unitario BIGINT NOT NULL DEFAULT 0,
    valor_total BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS documentos_proposta (
    id BIGSERIAL PRIMARY KEY,
    proposta_id BIGINT NOT NULL REFERENCES propostas(id) ON DELETE CASCADE,
    tipo VARCHAR(50) NOT NULL,
    nome_arquivo VARCHAR(255) NOT NULL,
    caminho TEXT NOT NULL,
    data_upload TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pareceres (
    id BIGSERIAL PRIMARY KEY,
    proposta_id BIGINT NOT NULL REFERENCES propostas(id) ON DELETE CASCADE,
    etapa VARCHAR(50) NOT NULL,
    avaliador_id BIGINT NOT NULL,
    nota INTEGER NOT NULL DEFAULT 0,
    parecer_texto TEXT NOT NULL DEFAULT '',
    data TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_propostas_edital_id ON propostas(edital_id);
CREATE INDEX IF NOT EXISTS idx_propostas_proponente_id ON propostas(proponente_id);
CREATE INDEX IF NOT EXISTS idx_propostas_status ON propostas(status);
CREATE INDEX IF NOT EXISTS idx_propostas_protocolo ON propostas(protocolo);
CREATE INDEX IF NOT EXISTS idx_itens_orcamentarios_proposta_id ON itens_orcamentarios(proposta_id);
CREATE INDEX IF NOT EXISTS idx_documentos_proposta_proposta_id ON documentos_proposta(proposta_id);
CREATE INDEX IF NOT EXISTS idx_pareceres_proposta_id ON pareceres(proposta_id);
CREATE INDEX IF NOT EXISTS idx_versoes_proposta_proposta_id ON versoes_proposta(proposta_id);
