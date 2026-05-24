-- name: InserirTokenRedefinicao :exec
INSERT INTO tokens_redefinicao (usuario_id, token, expirado_em)
VALUES ($1, $2, $3);

-- name: BuscarTokenRedefinicao :one
SELECT id, usuario_id, token, expirado_em, consumido_em, criado_em
FROM tokens_redefinicao
WHERE token = $1 AND consumido_em IS NULL;

-- name: ConsumirTokenRedefinicao :exec
UPDATE tokens_redefinicao
SET consumido_em = NOW()
WHERE token = $1;
