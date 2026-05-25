CREATE TABLE IF NOT EXISTS editais (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT NOT NULL,
    data_inicio DATE NOT NULL,
    data_fim DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ativo',
    tipo_chamada VARCHAR(20) NOT NULL,
    nota_de_corte INTEGER NOT NULL DEFAULT 0,
    valor_global BIGINT NOT NULL DEFAULT 0,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_editais_status ON editais(status);
CREATE INDEX idx_editais_tipo_chamada ON editais(tipo_chamada);
