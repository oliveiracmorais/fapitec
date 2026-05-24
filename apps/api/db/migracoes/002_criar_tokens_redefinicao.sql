CREATE TABLE IF NOT EXISTS tokens_redefinicao (
    id BIGSERIAL PRIMARY KEY,
    usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    token VARCHAR(64) NOT NULL UNIQUE,
    expirado_em TIMESTAMPTZ NOT NULL,
    consumido_em TIMESTAMPTZ,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tokens_redefinicao_token ON tokens_redefinicao(token);
CREATE INDEX idx_tokens_redefinicao_usuario_id ON tokens_redefinicao(usuario_id);
